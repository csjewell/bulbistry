/*
Copyright © 2023 Curtis Jewell <bulbistry@curtisjewell.name>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package bulbistry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorInfo struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details *interface{} `json:"details,omitempty"`
}

type ResponseText struct {
	Errors []ErrorInfo `json:"errors"`
}

type UnauthorizedDetails struct {
	Attempted string `json:"username_attempted"`
	Allowed   string `json:"username_allowed"`
}

func NoLogin(w http.ResponseWriter) {
	ret, _ := json.Marshal(ResponseText{
		Errors: []ErrorInfo{
			{
				Code:    "DENIED",
				Message: "Please provide username and password",
				Details: nil,
			},
		},
	})

	w.Header().Set("Content-Type", `application/json`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(ret)
	return
}

func ConfigError(w http.ResponseWriter, err error) {
	str := fmt.Sprintf("Configuration error: %v", err)

	ret, _ := json.Marshal(ResponseText{
		Errors: []ErrorInfo{
			{
				Code:    "DENIED",
				Message: str,
				Details: nil,
			},
		},
	})

	w.Header().Set("Content-Type", `application/json`)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(ret)
	return
}

func InvalidLogin(w http.ResponseWriter) {
	ret, _ := json.Marshal(ResponseText{
		Errors: []ErrorInfo{
			{
				Code:    "DENIED",
				Message: "Username or password is incorrect",
				Details: nil,
			},
		},
	})

	w.Header().Set("Content-Type", `application/json`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(ret)
	return
}

func RepositoryNotFound(w http.ResponseWriter, name string) {
	ret, _ := json.Marshal(ResponseText{
		Errors: []ErrorInfo{
			{
				Code:    "NAME_UNKNOWN",
				Message: "Repository named " + name + " cannot be found",
				Details: nil,
			},
		},
	})

	w.Header().Set("Content-Type", `application/json`)
	w.WriteHeader(http.StatusNotFound)
	w.Write(ret)
	return
}

//func Unauthorized(usernameAttempted string, usernameAllowed string) ResponseText {
//	// ret, _ := json.Marshal(
//	ret := ResponseText{
//		Errors: []ErrorInfo{
//			ErrorInfo{
//				Code:    "UNAUTHORIZED",
//				Message: "Username is not allowed to create, update, or delete this resource",
//				Details: &UnauthorizedDetails{Attempted: usernameAttempted, Allowed: usernameAllowed},
//			},
//		},
//	}
//	return ret
//}
