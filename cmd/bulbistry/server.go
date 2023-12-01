package main

import (
	"internal/config"

	"context"
	"net"
	"net/http"
	"time"

	tbv "github.com/csjewell/bulbistry"
	"github.com/gorilla/mux"
	htpasswd "github.com/tg123/go-htpasswd"
	"github.com/urfave/cli/v2"
)

func RunServer(ctx *cli.Context) error {
	cfg, err := config.ReadConfig(ctx.String("config"))
	if err != nil {
		executionLog.Fatal(err.Error())
	}

	needAuth := false
	if cfg.HTPasswdFile != "" {
		authorizer, err = htpasswd.New(cfg.HTPasswdFile, htpasswd.DefaultSystems, nil)
		if err != nil {
			executionLog.Fatal(err.Error())
		}
		needAuth = true
	}

	r := mux.NewRouter().StrictSlash(false).SkipClean(true).UseEncodedPath()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://example.com")
		w.Header().Set("Access-Control-Max-Age", "86400")
	}).Methods(http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	s := r.PathPrefix("/v2").Subrouter()
	ba := tbv.NewBasicAuthMiddleware(needAuth, authorizer)
	s.Use(ba.Middleware)

	s.HandleFunc("/", tbv.GetV2Check).Methods(http.MethodGet)

	//s.HandleFunc("/{manifestName}/manifests/{reference:sha256:[0-9a-fA-F]+}", tbv.HeadManifest).Methods(http.MethodHead)
	//s.HandleFunc("/{namespace}/{manifestName}/manifests/{reference:sha256:[0-9a-fA-F]+}", tbv.HeadNamespacedManifest).Methods(http.MethodHead)

	s.HandleFunc("/{manifestName}/manifests/{reference:sha256:[0-9a-fA-F]+}", tbv.GetManifest).Methods(http.MethodGet).Name("Manifest")
	s.HandleFunc("/{namespace}/{manifestName}/manifests/{reference:sha256:[0-9a-fA-F]+}", tbv.GetNamespacedManifest).Methods(http.MethodGet).Name("NamespacedManifest")

	s.HandleFunc("/{manifestName}/manifests/{reference}", tbv.HeadRedirectManifest).Methods(http.MethodHead)
	s.HandleFunc("/{namespace}/{manifestName}/manifests/{reference}", tbv.HeadRedirectNamespacedManifest).Methods(http.MethodHead)

	s.HandleFunc("/{manifestName}/manifests/{reference}", tbv.GetRedirectManifest).Methods(http.MethodGet)
	s.HandleFunc("/{namespace}/{manifestName}/manifests/{reference}", tbv.GetRedirectNamespacedManifest).Methods(http.MethodGet)

	// Paths to implement for pulling:
	// GET /v2/<name>/manifests/<reference> *
	// GET /v2/<name>/blobs/<digest>
	// HEAD /v2/<name>/manifests/<reference>
	// HEAD /v2/<name>/blobs/<digest>

	// s.HandleFunc("/{manifestName}/blobs/uploads", tbv.PostBlobUpload).Methods(http.MethodPost)
	// s.HandleFunc("/{namespace}/{manifestName}/blobs/uploads", tbv.PostNamespacedBlobUpload).Methods(http.MethodPost)

	// Paths to implement for pushing blobs
	// POST /v2/<name>/blobs/uploads/ returns 202 Accepted, and header Location: <location>
	// Then a PUT to <location>?digest=<digest> that returns 201 Created with Location: <blob URL>
	// OR POST /v2/<name>/blobs/uploads/?digest=<digest> (or can return a 202 Accepted as above to require the PUT)
	// OR POST-PATCH-PUT format.
	// Or mount from another <name> if we know it.

	// Path to implement for pushing manifests
	// PUT /v2/<name>/manifests/<reference>

	// Path to implement for Content Discovery
	// GET /v2/<name>/tags/list (returns JSON)
	s.HandleFunc("/{manifestName}/tags/list", tbv.GetTags).Methods(http.MethodGet)
	s.HandleFunc("/{namespace}/{manifestName}/tags/list", tbv.GetNamespacedTags).Methods(http.MethodGet)

	// Path to implement content deletion
	// DELETE /v2/<name>/manifests/<tag>    // tags
	// DELETE /v2/<name>/manifests/<digest> // manifests
	// DELETE /v2/<name>/blobs/<digest>     // blobs

	r.HandleFunc("/", handlerNotFound).Methods(http.MethodGet, http.MethodPost, http.MethodHead)

	svr := &http.Server{
		Handler:        r,
		Addr:           cfg.GetListenOn(),
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
		BaseContext:    htcontext,
	}

	executionLog.Fatal(svr.ListenAndServe())
	return nil
}

func htcontext(_ net.Listener) context.Context {
	ctx := context.WithValue(context.Background(), tbv.ConfigKey, cfg)
	return ctx
}

func handlerNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	return
}
