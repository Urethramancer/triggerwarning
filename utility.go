package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/sha3"
)

func pr(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}

func check(e error) {
	if e != nil {
		pr("Error: %s", e.Error())
		os.Exit(2)
	}
}

func genString(size int, complex bool) string {
	// This should be password-friendly
	valid := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-")
	if complex {
		// This isn't
		valid = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!'#$%&/()=?@*^<>-.:,;|[]{}")
	}
	pw := make([]byte, size)
	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(valid))))
		if err != nil {
			return ""
		}
		c := valid[n.Int64()]
		pw[i] = c
	}
	return string(pw)
}

func hashString(s string) string {
	hash := sha3.New256()
	hash.Reset()
	_, err := io.WriteString(hash, s)
	if err != nil {
		crit("Error: %s", err.Error())
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func getTime(t string) time.Duration {
	s := stripNonNumeric(t)
	count, err := strconv.Atoi(s)
	if err != nil {
		// A minute should be acceptable if conversion fails for some reason.
		return time.Minute
	}
	unit := getTimeUnit(t)
	return time.Duration(count) * unit
}

func stripNonNumeric(s string) string {
	num := func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}
	return strings.Map(num, s)
}

func stripNonLetter(s string) string {
	num := func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return -1
	}
	return strings.Map(num, s)
}

func getTimeUnit(s string) time.Duration {
	u := stripNonLetter(s)
	switch u {
	case "s":
		return time.Second
	case "m":
		return time.Minute
	case "h":
		return time.Hour
	case "d":
		return time.Hour * 24
	case "w":
		return time.Hour * 24 * 7
	case "M":
		return time.Hour * 24 * 30
	}
	return time.Second
}

// getVisitorIP simplistically tries to get the X-Real-IP header first,
// then whatever address the request itself reports. Tested with nginx.
func getVisitorIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
