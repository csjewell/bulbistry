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
	v "internal/version"

	"errors"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type bulbistryConfigError struct {
	configKeys []string
	error
}

func newConfigError(keys []string) bulbistryConfigError {
	return bulbistryConfigError{
		keys,
		errors.New("configuration entry required: " + strings.Join(keys, ", ")),
	}
}

var cfgFile string
var debugLogging bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "bulbistry",
	Version: v.Version(),
	Short:   "A pared-down container registry server, perfect for self-hosting",
	Long: `bulbistry version ` + v.Version() + `
	A pared-down (soon to be OCI-compliant) container registry server, perfect for self-hosting.
	By default, it starts the server on the port specified in the configuration or the environment.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig(false)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Server not hooked in yet.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bulbistry/env)")
	rootCmd.PersistentFlags().BoolVar(&debugLogging, "debug", false, "Turn on debug logging (default is false)")
}

// initConfig reads in config file and ENV variables if set.
// The boolean flag specifies whether to set defaults (true) or return errors (false)
func initConfig(_ bool) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigType("dotenv")
		viper.SetConfigFile(cfgFile)
	} else {
		// Find executable directory.
		executable, err := os.Executable()
		if err != nil {
			return err
		}
		executableDir := path.Dir(executable)

		// Search config in same directory as executable with name ".env".
		viper.SetConfigType("dotenv")
		viper.AddConfigPath(executableDir)
		viper.SetConfigName(".env")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Using config file:", viper.ConfigFileUsed())

		settingKeys := viper.AllKeys()
		settings := make(map[string]any, 50)

		for _, key := range settingKeys {
			if strings.Contains(key, "_") {
				newKey := strings.ReplaceAll(key, "_", ".")
				settings[newKey] = viper.Get(key)
			}
		}

		viper.Reset()

		for k, v := range settings {
			viper.Set(k, v)
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	configKeyErr := make([]string, 5)

	if viper.IsSet("file.database") {
		viper.SetDefault("file.database", viper.Get("file.database"))
	} else {
		configKeyErr = append(configKeyErr, "FILE_DATABASE")
	}

	viper.SetDefault("file.htpasswd", "")

	if viper.IsSet("registry.url.hostname") {
		viper.SetDefault("registry.url.hostname", viper.Get("registry.url.hostname"))
	} else {
		configKeyErr = append(configKeyErr, "REGISTRY_URL_HOSTNAME")
	}

	if viper.IsSet("blob.directory") {
		viper.SetDefault("blob.directory", viper.Get("blob.directory"))
	} else {
		configKeyErr = append(configKeyErr, "BLOB_DIRECTORY")
	}

	viper.SetDefault("registry.url.port", 80)
	viper.SetDefault("registry.url.path", "/")
	viper.SetDefault("registry.url.scheme", "http")

	viper.SetDefault("blob.proxied", false)

	if viper.GetBool("blob.proxied") {
		if !viper.IsSet("blob.url.hostname") {
			configKeyErr = append(configKeyErr, "BLOB_URL_HOSTNAME")
		}
	} else {
		viper.SetDefault("blob.url.hostname", viper.GetString("registry.url.hostname"))
	}

	viper.SetDefault("blob.url.port", 80)
	viper.SetDefault("blob.url.path", "/blob")
	viper.SetDefault("blob.url.scheme", "http")

	if len(configKeyErr) > 0 {
		return newConfigError(configKeyErr)
	}

	return nil
}
