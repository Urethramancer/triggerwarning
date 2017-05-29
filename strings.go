package main

const (
	// StatusOK means no problems.
	StatusOK = iota
	// StatusNoCredentials means wrong inputs to login.
	StatusNoCredentials = iota
	// StatusNoAuth means no or wrong auth token was passed.
	StatusNoAuth = iota
	// StatusAuthFailed means username or password was wrong, or account non-existent.
	StatusAuthFailed = iota
	// StatusNoAccess means the token given does not have access to the requested data.
	StatusNoAccess = iota
	// StatusError is for generic error messages.
	StatusError = iota
)

var messages = map[int]string{
	StatusOK:            "OK",
	StatusNoCredentials: "Missing username or password.",
	StatusNoAuth:        "No valid authorisation token provided.",
	StatusAuthFailed:    "Login failed.",
	StatusNoAccess:      "Not accessible with provided credentials.",
}

// API endpoints
const (
	// PathWatch requires 'name', 'user' and 'token' parameters.
	PathWatch = "/watch"
	// PathUnwatch requires 'name', 'user' and 'token' parameters.
	PathUnwatch = "/unwatch"
	// PathTrigger requires 'name', 'user' and 'token' parameters.
	PathTrigger = "/trigger"
)
