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

// The code that handles the parsed configuration for Bulbistry
package config

import (
	"errors"
	"fmt"
	"internal/database"
	"net/url"

	"github.com/spf13/viper"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.

type bulbistryConfigError struct {
	configKey string
	error
}

func newConfigError(key, err string) bulbistryConfigError {
	return bulbistryConfigError{key, errors.New(err + ": " + key)}
}

func getHostname(urlKey string) string {
	urlSub := viper.Sub(urlKey)
	scheme := urlSub.GetString("scheme")
	port := urlSub.GetInt("port")
	host := urlSub.GetString("hostname")

	if scheme == "http" && port == 80 {
		return host
	}

	if scheme == "https" && port == 443 {
		return host
	}

	return host + ":" + fmt.Sprint(port)
}

func getURL(urlKey string) *url.URL {
	urlSub := viper.Sub(urlKey)
	return &url.URL{
		Scheme: urlSub.GetString("scheme"),
		Host:   getHostname(urlKey),
		Path:   urlSub.GetString("path"),
	}
}

func GetExternalOrigin() *url.URL {
	return getURL("external_url")
}

func GetExternalURL() *url.URL {
	return getURL("external_url").JoinPath("/v2/")
}

// GetListenOn gets the IP and port that the registry is configured to listen on
func GetListenOn() string {
	return viper.GetString("listen_on.ip") + ":" + fmt.Sprint(viper.GetInt("listen_on.port"))
}

// CheckConfig checks that the current configuration is valid
func CheckConfig() error {

	return nil
}

// GetManifestURL gets the URL to retrieve a particular manifest
func GetManifestURL(mt database.ManifestTag) string {
	if mt.Namespace == "" {
		return GetExternalURL().JoinPath(mt.Name, "/manifest/", mt.Sha512).String()
	}
	return GetExternalURL().JoinPath(mt.Namespace, mt.Name, "/manifest/", mt.Sha512).String()
}

// GetBlobURL gets the blob storage base URL.
func GetBlobURL() string {
	//return GetExternalURL().Scheme + "://" + GetExternalURL().HostName + ":" + GetExternalURL().Port + "/v2/"
	return ""
}
