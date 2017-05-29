package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Result is returned from many JSON requests.
type Result struct {
	Status     string `json:"status,omitempty"`
	StatusCode int    `json:"statuscode"`
	Token      string `json:"token,omitempty"`
}

// SetStatus fills in the status text from a code.
func (r *Result) SetStatus(code int) {
	s, ok := messages[code]
	if !ok {
		r.StatusCode = StatusOK
		r.Status = messages[StatusOK]
		return
	}
	r.StatusCode = code
	r.Status = s
}

func respond(w io.Writer, r interface{}) error {
	res, err := json.Marshal(r)
	if err != nil {
		return err
	}

	_, err = w.Write(res)
	if err != nil {
		return err
	}

	return nil
}

// LoadJSON loads a JSON file into a structure.
func LoadJSON(fn string, out interface{}) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}

// SaveJSON writes any structure to file as JSON.
func SaveJSON(fn string, data interface{}) error {
	var out []byte
	out, err := json.MarshalIndent(data, "\t", "")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fn, out, 0600)
}
