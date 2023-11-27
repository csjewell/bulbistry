package bulbistry

import (
	"encoding/json"
	"fmt"
)

type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message`
	Details interface{} `json:"details"`
}

type ErrorArray []ErrorInfo

type ResponseText struct {
	Errors interface{}
}

type UnauthorizedDetails struct {
	Attempted string `json:"username_attempted"`
	Allowed   string `json:"username_allowed"`
}

func NoLogin() []byte {
	ret, _ := json.Marshal(ResponseText{Errors: ErrorArray{ErrorInfo{Code: "DENIED", Message: "Please provide username and password", Details: nil}}})
	return ret
}

func ConfigError(err error) []byte {
	str := fmt.Sprintf("Configuration error: %v", err)
	ret, _ := json.Marshal(ResponseText{Errors: ErrorArray{ErrorInfo{Code: "DENIED", Message: str, Details: nil}}})
	return ret
}

func InvalidLogin() []byte {
	ret, _ := json.Marshal(ResponseText{Errors: ErrorArray{ErrorInfo{Code: "DENIED", Message: "Username or password is incorrect", Details: nil}}})
	return ret
}

func Unauthorized(usernameAttempted string, usernameAllowed string) []byte {
	ret, _ := json.Marshal(ResponseText{
		Errors: ErrorArray{ErrorInfo{
			Code:    "UNAUTHORIZED",
			Message: "Username is not allowed to create, update, or delete this resource",
			Details: UnauthorizedDetails{Attempted: usernameAttempted, Allowed: usernameAllowed},
		}},
	})
	return ret
}
