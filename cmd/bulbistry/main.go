package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
//    "database/sql"

//    tbv "git.curtisjewell.dev/bulbistry/bulbistry"
	"github.com/urfave/cli/v2"
	"github.com/gorilla/mux"
//	"github.com/tg123/go-htpasswd"
)

// github.com/google/uuid

// var ht htpasswd

func CreateConfig (ctx *cli.Context) error {
    return fmt.Errorf("Not Implemented")
}

func InitializeDatabase (ctx *cli.Context) error {
    return fmt.Errorf("Not Implemented")
}

func MigrateDatabase (ctx *cli.Context) error {
    return fmt.Errorf("Not Implemented")
}

func main() {

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n", c.App.Version)
	}

	app := &cli.App{
		Name:     "Bulbistry",
		HelpName: "A pared-down container registry, perfect for self-hosting",
		Authors:  []*cli.Author{ { Name: "Curtis Jewell", Email: "swordsman@curtisjewell.name", } },
		Copyright: "Copyright (c) 2003 Curtis Jewell",
		//Version:  tbv.FormatVersion(),
		Action:   func (c *cli.Context) error { return RunServer(c) },
		Flags:    []cli.Flag{
			&cli.BoolFlag{
				Name:   "create-config",
				Value:  false,
				Usage:  "Generate a basic configuration",
                Action: func (c *cli.Context, b bool) error {
					if b { return nil }
					return CreateConfig(c)
				},
			},
			&cli.BoolFlag{
				Name:   "init-db",
				Value:  false,
				Usage:  "Initialize application database",
                Action: func (c *cli.Context, b bool) error {
					if b { return nil }
					return InitializeDatabase(c)
				},
			},
			&cli.BoolFlag{
				Name:   "migrate-db",
				Value:  false,
				Usage:  "Migrate the database",
                Action: func (c *cli.Context, b bool) error {
					if b { return nil }
					return MigrateDatabase(c)
				},
			},
			&cli.StringFlag{
				Name:      "config",
				Aliases:   []string{"c"}, 
				Value:     "bulbistry.yml",
				Usage:     "Load configuration from FILE",
				TakesFile: true,
				EnvVars:   []string{"BULBISTRY_CONFIG"}, 
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Output(2, err.Error())
		os.Exit(1)
	}
}

func RunServer (ctx *cli.Context) error {

    if ctx.Bool("create-config") {
		return nil
	}

	if ctx.Bool("init-db") {
		return nil
	}

	if ctx.Bool("migrate-db") {
		return nil
	}

	fmt.Println("Read configuration file")
//	bc, err := tbv.readConfig()

	fmt.Println("Start server")

//	ht = htpasswd.New(bc.HTPasswdFile, htpasswd.DefaultSystems, nil)

	r := mux.NewRouter()
	r.Use()
	// r.Use(mux.CORSMethodMiddleware(r))

//	s := r.PathPrefix("/v2").Subrouter()

//	s.HandleFunc("/", tbv.V2Handler)

	// Paths to implement for pulling:
    // GET /v2/<name>/manifests/<reference>
	// GET /v2/<name>/blobs/<digest>
    // HEAD /v2/<name>/manifests/<reference>
	// HEAD /v2/<name>/blobs/<digest>

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

	// Path to implement content deletion
	// DELETE /v2/<name>/manifests/<tag>    // tags
	// DELETE /v2/<name>/manifests/<digest> // manifests
	// DELETE /v2/<name>/blobs/<digest>     // blobs

	r.HandleFunc("/", handlerNotFound)
	http.Handle("/", r)

	svr := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		BaseContext:    htcontext,
	}

	log.Fatal(svr.ListenAndServe());
	return nil;
}

func htcontext (l net.Listener) context.Context { 
	ctx := context.Background()
//	ctx = ctx.WithValue(tbv.HTPasswordKey, ht)
    return ctx;
}

func handlerNotFound(_ http.ResponseWriter, _ *http.Request) {
    return
}
