/*
Copyright © 2023 Curtis Jewell <swordsman@curtisjewell.name>

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
	"internal/config"
	"internal/database"

	"context"
	"encoding/json"

	//"fmt"
	"net/http"
	//"strings"
	"strconv"
	//"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	htpasswd "github.com/tg123/go-htpasswd"
)

type contextKey int

const (
	ConfigKey contextKey = iota
	userKey
)

type BasicAuth struct {
	authRequired bool
	authorizer   *htpasswd.File
}

func NewBasicAuthMiddleware(authRequired bool, authorizer *htpasswd.File) BasicAuth {
	return BasicAuth{
		authRequired: authRequired,
		authorizer:   authorizer,
	}
}

func (ba *BasicAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !ba.authRequired {
			next.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Bulbistry OCI Artifact Registry", charset="UTF-8"`)
			NoLogin(w)
			return
		}

		ok = ba.authorizer.Match(username, password)

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Bulbistry OCI Artifact Registry", charset="UTF-8"`)
			InvalidLogin(w)
			return
		}

		ctxWithUser := context.WithValue(ctx, userKey, username)
		rWithUser := r.WithContext(ctxWithUser)
		next.ServeHTTP(w, rWithUser)
	})
}

func getDB(r *http.Request) (*database.Database, error) {
	db, err := database.NewDatabase(viper.GetString("database_file"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetV2Check(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", `text/plain`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func writeError(err error, w http.ResponseWriter) {

}

func writeNotFound(err error, w http.ResponseWriter) {

}

func HeadRedirectManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db, err := getDB(r)
	if err != nil {
		writeError(err, w)
		return
	}

	mt, err := db.GetManifestTag(vars["manifestName"], vars["reference"])
	if err != nil {
		writeNotFound(err, w)
		return
	}

	commonRedirectManifest(false, *mt, w)
}

func HeadRedirectNamespacedManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := getDB(r)
	if err != nil {
		writeError(err, w)
		return
	}

	mt, err := db.GetNamespacedManifestTag(vars["namespace"], vars["manifestName"], vars["reference"])
	if err != nil {
		writeNotFound(err, w)
		return
	}

	commonRedirectManifest(false, *mt, w)
}

func GetRedirectManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db, err := getDB(r)
	if err != nil {
		writeError(err, w)
		return
	}

	mt, err := db.GetManifestTag(vars["manifestName"], vars["reference"])
	if err != nil {
		writeNotFound(err, w)
		return
	}

	commonRedirectManifest(true, *mt, w)
}

func GetRedirectNamespacedManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := getDB(r)
	if err != nil {
		writeError(err, w)
		return
	}

	mt, err := db.GetNamespacedManifestTag(vars["namespace"], vars["manifestName"], vars["reference"])
	if err != nil {
		writeNotFound(err, w)
		return
	}

	commonRedirectManifest(true, *mt, w)
}

func commonRedirectManifest(hasBody bool, mt database.ManifestTag, w http.ResponseWriter) {
	if mt.ID != 0 {
		url := config.GetManifestURL(mt)
		w.Header().Set("ETag", `"`+mt.Sha256+`"`)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Docker-Content-Digest", mt.Sha256)
		if hasBody {
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusPermanentRedirect)
			w.Write([]byte("Redirected"))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetTags(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// kick out with a 500 error.
	}

	nStr := r.FormValue("n")
	var n int
	if nStr != "" {
		n, err = strconv.Atoi(nStr)
		if err != nil {
			// Kick out with a 500 error
		}
	}

	db, err := getDB(r)

	vars := mux.Vars(r)
	tags, err := db.GetTags(vars["manifestName"], n, r.FormValue("last"))
	if err != nil {
		// Kick out with a 404 or 500 as appropriate.
	}

	jsonTags, err := json.Marshal(tags)
	if err != nil {
		// 500 error
	}

	w.Header().Set("Content-Type", `application/json`)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonTags)
}

func GetNamespacedTags(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// kick out with a 500 error.
	}

	nStr := r.FormValue("n")
	var n int
	if nStr != "" {
		n, err = strconv.Atoi(nStr)
		if err != nil {
			// Kick out with a 500 error
		}
	}

	db, err := getDB(r)

	vars := mux.Vars(r)
	tags, err := db.GetNamespacedTags(vars["namespace"], vars["manifestName"], n, r.FormValue("last"))
	if err != nil {
		// Kick out with a 404 or 500 as appropriate.
	}

	jsonTags, err := json.Marshal(tags)
	if err != nil {
		// 500 error
	}

	w.Header().Set("Content-Type", `application/json`)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonTags)
}

// Routines up to this point are at least psuedocoded

func GetManifest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", `text/plain`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func GetNamespacedManifest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", `text/plain`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func UploadLocation(namespace *string, name string, uuid uuid.UUID) string {
	return ""
}

func StoreInProgressUpload(namespace *string, name string, nano string) {

}

//func PostBlobUpload(w http.ResponseWriter, r *http.Request) {
//	err := r.ParseForm()
//	if err != nil {
//		// kick out with a 500 error.
//	}
//
//	ctx := r.Context()
//	f := ctx.Value(ConfigKey)
//	if f == nil {
//		// Kick out with a 500 error
//	}
//
//	bc, ok := f.(BulbistryConfig)
//	if !ok {
//		w.Write(ConfigError(errors.New("Configuration not loadable")))
//		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
//		return
//	}
//
//	vars := mux.Vars(r)
//	name := vars["manifestName"]
//
//	key := fmt.Sprint(time.Now().UnixNano())
//
//	StoreInProgressUpload(nil, name, key)
//
//	blobUUID, _ := GenerateBlobUUID(bc, name, key)
//
//	w.Header().Set("Location", UploadLocation(nil, name, blobUUID))
//	w.WriteHeader(http.StatusAccepted)
//}

func PostNamespacedBlobUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", `text/plain`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
