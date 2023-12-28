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
package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// db "internal/database"
	"github.com/distribution/distribution/v3"
	"github.com/distribution/distribution/v3/manifest/manifestlist"

	// "github.com/distribution/distribution/v3/manifest/ocischema"
	"github.com/distribution/distribution/v3/manifest/schema2"
	"github.com/imroc/req/v3"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
)

// registryPullCmd represents the registryPull command
var registryPullCmd = &cobra.Command{
	Use:   "pull [container]",
	Short: "Pulls a container from another registry into bulbistry",
	Long: `Pulls a container from another registry into bulbistry.
	The container being pulled must be publicly accessible.`,
	Run: func(cmd *cobra.Command, args []string) {
		download(cmd, args)
	},
}

func init() {
	registryCmd.AddCommand(registryPullCmd)
}

type connection struct {
	client      *req.Client
	site        string
	packageName string
	authHeader  string
}

type genericToken struct {
	Token        string `json:"token"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	IssuedAt     string `json:"issued_at,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func getPackageInfo(packageNameIn string) (connection, string) {
	var possiblePackageName string
	var tag string

	if strings.Contains(packageNameIn, "@") {
		possiblePackageName, tag, _ = strings.Cut(packageNameIn, "@")
	} else {
		tag = "latest"
		possiblePackageName = packageNameIn
	}

	var conn connection

	possibleSite, shortPackageName, isCut := strings.Cut(possiblePackageName, "/")
	if isCut && strings.Contains(possibleSite, ".") {
		conn.site = possibleSite
		conn.packageName = shortPackageName
	} else {
		conn.site = "registry.docker.io"
		conn.packageName = possiblePackageName
	}

	conn.client = req.NewClient().
		SetTimeout(5 * time.Second).
		SetUserAgent("bulbistry-dl/0.0.5")

	conn.client.SetRedirectPolicy(req.NoRedirectPolicy())

	return conn, tag
}

func (conn *connection) NewRequest(accept []string) *req.Request {
	rqst := conn.client.NewRequest()
	if conn.authHeader != "" {
		rqst.SetHeader(http.CanonicalHeaderKey("Authorization"), conn.authHeader)
	}
	for _, s := range accept {
		rqst.SetHeader(http.CanonicalHeaderKey("Accept"), s)
	}
	return rqst
}

func (conn *connection) getAuthHeader(resp *req.Response) {
	authTag := make(map[string]string, 4)

	auth := resp.Header[http.CanonicalHeaderKey("WWW-Authenticate")][0]
	if auth == "" {
		return
	}

	authType, authInfo, _ := strings.Cut(auth, " ")
	if authType != "Bearer" {
		log.Fatal("Do not know this authentication format")
	}

	entries := strings.Split(authInfo, ",")
	for _, s := range entries {
		parts := strings.Split(s, "=")
		value := parts[1]
		if len(value) > 0 && value[0] == '"' {
			value = value[1:]
		}
		if len(value) > 0 && value[len(value)-1] == '"' {
			value = value[:len(value)-1]
		}
		authTag[parts[0]] = value
	}

	if authTag["scope"] == "repository:user/image:pull" {
		authTag["scope"] = "repository:" + conn.packageName + ":pull"
	}

	var tokenAuth genericToken
	if strings.Contains(authTag["realm"], "ghcr.io") || strings.Contains(authTag["realm"], "docker.io") {
		rqst := conn.client.NewRequest().
			AddQueryParam("service", authTag["service"]).
			SetSuccessResult(&tokenAuth)

		if authTag["scope"] != "" {
			rqst = rqst.AddQueryParam("scope", authTag["scope"])
		}

		resp, err := rqst.Get(authTag["realm"])

		if err != nil {
			log.Fatal(err)
		}

		if resp.IsSuccessState() {
			conn.authHeader = "Bearer " + tokenAuth.Token
		}
	}
}

func download(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cobra.CheckErr(fmt.Errorf("bulbistry-dl needs a package to download"))
	}

	conn, tag := getPackageInfo(args[0])

	resp, err := conn.client.NewRequest().Get("https://" + conn.site + "/v2/")
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 401 {
		conn.getAuthHeader(resp)
	} else if resp.StatusCode != 200 {
		log.Fatal("URL not a registry")
	}

	req := conn.NewRequest([]string{
		// "application/json",
		v1.MediaTypeImageManifest,
		schema2.MediaTypeManifest,
		manifestlist.MediaTypeManifestList,
		v1.MediaTypeImageIndex,
	})

	resp, err = req.Head("https://" + conn.site + "/v2/" + conn.packageName + "/manifests/" + tag)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(errors.New("Could not find the tag"))
	}

	var tagCanonical string
	if resp.StatusCode == 200 {
		tagCanonical = resp.Header[http.CanonicalHeaderKey("Docker-Content-Digest")][0]
		if tagCanonical != "" {
			// Store the mapping between tag and tagCanonical in the database
		} else {
			tagCanonical = tag
		}
	} else {
		log.Fatal(errors.New("Could not find the canonical tag"))
	}

	req = conn.NewRequest([]string{
		// "application/json",
		v1.MediaTypeImageManifest,
		schema2.MediaTypeManifest,
		manifestlist.MediaTypeManifestList,
		v1.MediaTypeImageIndex,
	})

	resp, err = req.Get("https://" + conn.site + "/v2/" + conn.packageName + "/manifests/" + tagCanonical)
	if err != nil {
		log.Fatal(err)
	}

	ct := resp.Header[http.CanonicalHeaderKey("Content-Type")][0]
	manifestBytes, _ := resp.ToBytes()

	f, err := os.OpenFile(tagCanonical[7:]+".manifest.blob", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(manifestBytes)
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	manifest, _, err := distribution.UnmarshalManifest(ct, manifestBytes)
	if err != nil {
		log.Fatal(err)
	}

	if ct == manifestlist.MediaTypeManifestList {
		processManifestList(conn, manifest)
	} else if ct == v1.MediaTypeImageIndex {
		processManifestIndex(conn, manifest)
	} else {
		processManifest(conn, manifest)
	}

	// GET "https://" + site + "/" PackageName + "/manifests/" + tag as a manifest(list)
	//     store in manifest-tags, foreach tag mentioned:
	//     GET "https://" + site + "/" PackageName + "/manifests/" + tag as a manifest
	//         store in manifests, foreach config or layer mentioned,
	//             GET "https://" + site + "/" PackageName + "/blobs/" + tag as the type given
	//             store in blobs and blobindex.

}

func processManifestIndex(conn connection, manifest distribution.Manifest) {

	i := 1
	i = i + 2
	//     store in manifest-tags, foreach tag mentioned:
	//     GET "https://" + site + "/" PackageName + "/manifests/" + tag as a manifest
	// ProcessManifest(conn, text)
}

func processManifestList(conn connection, manifest distribution.Manifest) {

	for _, subManifest := range manifest.References() {
		retrieveManifest(conn, subManifest)
	}
}

func retrieveBlob(conn connection, subManifest distribution.Describable) {
	tag := string(subManifest.Descriptor().Digest)

	req := conn.NewRequest([]string{})

	resp, err := req.Get("https://" + conn.site + "/v2/" + conn.packageName + "/blobs/" + tag)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 307 {
		url := resp.Header.Get("Location")
		resp, err = req.Get(url)
		if err != nil {
			log.Fatal(err)
		}
	}

	manifestBytes, _ := resp.ToBytes()

	f, err := os.OpenFile(tag[7:]+".blob", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(manifestBytes)
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func retrieveManifest(conn connection, subManifest distribution.Describable) {
	tag := string(subManifest.Descriptor().Digest)

	req := conn.NewRequest([]string{
		// "application/json",
		v1.MediaTypeImageManifest,
		schema2.MediaTypeManifest,
		manifestlist.MediaTypeManifestList,
		v1.MediaTypeImageIndex,
	})

	resp, err := req.Get("https://" + conn.site + "/v2/" + conn.packageName + "/manifests/" + tag)
	if err != nil {
		log.Fatal(err)
	}

	ct := resp.Header.Get("Content-Type")
	manifestBytes, _ := resp.ToBytes()

	f, err := os.OpenFile(tag[7:]+".manifest.blob", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(manifestBytes)
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	manifest, _, err := distribution.UnmarshalManifest(ct, manifestBytes)
	processManifest(conn, manifest)
}

func processManifest(conn connection, manifest distribution.Manifest) {
	for _, blob := range manifest.References() {
		retrieveBlob(conn, blob)
	}
}
