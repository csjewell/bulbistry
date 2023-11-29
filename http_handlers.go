package bulbistry

import (
	"context"
	"errors"
	//"fmt"
	"net/http"
	//"strings"
	"strconv"
	//"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
			w.Write(NoLogin())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ok = ba.authorizer.Match(username, password)

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Bulbistry OCI Artifact Registry", charset="UTF-8"`)
			w.Write(InvalidLogin())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctxWithUser := context.WithValue(ctx, userKey, username)
		rWithUser := r.WithContext(ctxWithUser)
		next.ServeHTTP(w, rWithUser)
	})
}

func GetV2Check(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", `text/plain`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HeadRedirectManifest(_ http.ResponseWriter, _ *http.Request) {
//func HeadRedirectManifest(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(e)
	// Get manifest SHA and content type based on request parameters
	// commonRedirectManifest(0, manifest)
}

func HeadRedirectNamespacedManifest(_ http.ResponseWriter, _ *http.Request) {
// func HeadRedirectNamespacedManifest(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(e)
	// Get manifest SHA and content type based on request parameters
	// commonRedirectManifest(0, manifest)
}

func GetRedirectManifest(_ http.ResponseWriter, _ *http.Request) {
// func GetRedirectManifest(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(e)
	// Get manifest SHA and content type based on request parameters
	// commonRedirectManifest(1, manifest)
}

func GetRedirectNamespacedManifest(_ http.ResponseWriter, _ *http.Request) {
// func GetRedirectNamespacedManifest(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(e)
	// Get manifest SHA and content type based on request parameters
	// commonRedirectManifest(1, manifest)
}

func commonRedirectManifest(hasBody bool, mt ManifestTag, url string, w http.ResponseWriter) {
	if mt.ID != 0 {
		w.Header().Set("ETag", `"`+mt.Sha256+`"`)
		w.Header().Set("Content-Type", mt.ContentType)
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

	ctx := r.Context()
	f := ctx.Value(ConfigKey)
	if f == nil {
		// Kick out with a 500 error
	}

	bc, ok := f.(Config)
	if !ok {
		w.Write(ConfigError(errors.New("Configuration not loadable")))
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	nStr := r.FormValue("n")
	var n int
	if nStr != "" {
		n, err = strconv.Atoi(nStr)
		if err != nil {
			// Kick out with a 500 error
		}
	}

	vars := mux.Vars(r)
	db := NewDatabase(bc)
	db.GetTags(vars["manifestName"], n, r.FormValue("last"))
	// Print out JSON tag list
}

func GetNamespacedTags(_ http.ResponseWriter, _ *http.Request) {
// func GetNamespacedTags(w http.ResponseWriter, r *http.Request) {
	//	vars := mux.Vars(r)

	// Get manifest SHA and content type based on request parameters
	// Print out JSON tag list
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
