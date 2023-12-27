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
package cmd

import (
	"log/slog"
	"os"
	"path"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configInitCmd represents the configInit command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Writes out a basic configuration file",
	Long:  `Writes out a basic configuration file to be edited and used`,
	Run: func(cmd *cobra.Command, args []string) {
		err := configInitialize()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func configInitialize() error {
	var (
		file *os.File
		err  error
	)

	fileName := cfgFile
	if fileName == "" {
		executable, err := os.Executable()
		if err != nil {
			return err
		}
		executableDir := path.Dir(executable)
		fileName = path.Join(executableDir, ".env")
	}

	_, err = os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		file, err = os.OpenFile(fileName, os.O_WRONLY, 0755)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	envText := `#!/bin/sh

// FILE_DATABASE is where the SQLite database lives.
FILE_DATABASE={{.file.database}}

// FILE_HTPASSWD is where the users file lives.
// It can be empty, in which case, no authentication is done at all.
FILE_HTPASSWD={{.file.htpasswd}}

// REGISTRY_IP is the IP the registry server listens on.
// It should be set to 127.0.0.1 (listen locally) or 0.0.0.0 (listen globally)
// in most situations.
REGISTRY_IP={{.registry.ip}}

// REGISTRY_PORT is the port the registry server listens on.
REGISTRY_PORT={{.registry.port}}

// REGISTRY_URL_* is the different portions of the external registry URL.
REGISTRY_URL_HOSTNAME={{.registry.url.hostname}}
REGISTRY_URL_PATH={{.registry.url.path}}
REGISTRY_URL_PORT={{.registry.url.port}}
REGISTRY_URL_SCHEME={{.registry.url.scheme}}

// BLOB_DIRECTORY defines where the blobs are stored.
BLOB_DIRECTORY={{.blob.directory}}

// BLOB_PROXIED specifies whether to expect blobs to be proxied by an http proxy.
BLOB_PROXIED={{.blob.proxied}}

// BLOB_URL_* specifies the different parts of the URL that blobs will be accessible at.
// If BLOB_PROXIED is false, bulbistry will servve at this URL.
// If BLOB_PROXIED is true, then an external HTTP proxy is expected to serve these URLs.
BLOB_URL_HOSTNAME={{.blob.url.hostname}}
BLOB_URL_PATH={{.blob.url.path}}
BLOB_URL_PORT={{.blob.url.port}}
BLOB_URL_SCHEME={{.blob.url.scheme}}
`

	tmpl, err := template.New("dotenv").Parse(envText)
	if err != nil {
		return err
	}

	// Give defaults to parts of the configuration that may not exist.
	if !viper.IsSet("registry.url.hostname") {
		viper.Set("registry.url.hostname", "registry.localhost")
	}

	if !viper.IsSet("registry.ip") {
		viper.Set("registry.ip", "127.0.0.1")
	}

	if !viper.IsSet("registry.port") {
		viper.Set("registry.port", "8088")
	}

	if !viper.IsSet("blob.url.hostname") || viper.GetString("blob.url.hostname") == "" {
		viper.Set("blob.url.hostname", "registry.localhost")
	}

	if !viper.IsSet("blob.directory") {
		viper.Set("blob.directory", "/blob")
	}

	if !viper.IsSet("file.database") {
		viper.Set("file.database", "$HOME/bulbistry.db")
	}

	err = tmpl.Execute(file, viper.AllSettings())
	if err != nil {
		return err
	}
	slog.Info("Configuration file written")
	return nil
}

func init() {
	configCmd.AddCommand(configInitCmd)
}
