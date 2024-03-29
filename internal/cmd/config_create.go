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
	"log/slog"
	"net/url"
	"os/user"
	"path"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createConfigCmd represents the createConfig command
var configCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"config", "create_config"},
	Short:   "Creates a Bulbistry configuration",
	Long:    `Creates a Bulbistry configuration interactively`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return createConfig(cmd)
	},
}

func init() {
	configCmd.AddCommand(configCreateCmd)
}

// createConfig creates the bulbistry configuration in an interactive fashion.
func createConfig(cmd *cobra.Command) error {

	validateFile := func(input string) error {
		if input == "" {
			return errors.New("Filename cannot be empty")
		}
		// TODO: Other checks (exists as a directory, etc.)
		return nil
	}

	prompt := promptui.Prompt{
		Label:     "What is the registry hostname",
		Default:   "registry.local",
		AllowEdit: true,
		// Validate:  validateHostname,
	}

	hostName, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get hostname %v\n", err)
	}

	prompt = promptui.Prompt{
		Label:     "What is the URL for the registry",
		Default:   "https://" + hostName + "/",
		AllowEdit: true,
		// Validate:  validateURL,
	}

	urlString, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get registry URL %v\n", err)
	}

	urlRegistry, _ := url.Parse(urlString)
	if err != nil {
		return fmt.Errorf("Registry URL was not a URL: %v", err)
	}

	prompt = promptui.Prompt{
		Label:   "Where should the blobs be stored",
		Default: "/www-data/blob",
		// Validate:  validateDirectory,
		AllowEdit: true,
	}

	blobDirectory, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get blob storage directory %v\n", err)
	}
	viper.Set("blob.directory", blobDirectory)

	menu := promptui.Select{
		Label:     "Is this storage directory served by a proxy",
		CursorPos: 0,
		Items:     []string{"Yes", "No"},
	}

	_, isProxyStr, err := menu.Run()
	if err != nil {
		return fmt.Errorf("Did not get whether blobs are proxied %v\n", err)
	}

	isProxied := false
	if isProxyStr == "Yes" {
		isProxied = true
	}
	viper.Set("blob.proxied", isProxied)

	urlBlobDefault := urlRegistry.JoinPath("/blob")

	prompt = promptui.Prompt{
		Label:     "Where is the URL to where blobs are stored",
		Default:   urlBlobDefault.String(),
		AllowEdit: true,
		// Validate:  validateURL,
	}

	urlString, err = prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get blob URL %v\n", err)
	}

	urlBlob, err := url.Parse(urlString)
	if err != nil {
		return fmt.Errorf("Blob URL was not a URL: %v\n", err)
	}

	userInfo, _ := user.Current()
	prompt = promptui.Prompt{
		Label:     "Where is the SQLite database",
		Default:   path.Join(userInfo.HomeDir, "bulbistry.db"),
		Validate:  validateFile,
		AllowEdit: true,
	}

	databaseFile, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get SQLite database filename %v\n", err)
	}
	viper.Set("file.database", databaseFile)

	prompt = promptui.Prompt{
		Label:     "Where is the password file",
		Default:   path.Join(userInfo.HomeDir, ".htpasswd"),
		Validate:  validateFile,
		AllowEdit: true,
	}

	htpasswdFile, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get password filename %v\n", err)
	}
	viper.Set("file.htpasswd", htpasswdFile)

	menu = promptui.Select{
		Label:     "Where should the registry listen",
		CursorPos: 0,
		Items:     []string{"127.0.0.1", "0.0.0.0"},
	}

	_, ip, err := menu.Run()
	if err != nil {
		return fmt.Errorf("Did not get IP %v\n", err)
	}
	viper.Set("registry.ip", ip)

	menu = promptui.Select{
		Label:     "What port should the registry listen on",
		CursorPos: 2,
		Items:     []string{"80", "8080", "8088"},
	}

	_, portStr, err := menu.Run()
	if err != nil {
		return fmt.Errorf("Did not get port %v\n", err)
	}

	port, _ := strconv.Atoi(portStr)
	viper.Set("registry.port", port)

	euPort, _ := strconv.Atoi(urlRegistry.Port())
	bPort, _ := strconv.Atoi(urlBlob.Port())

	viper.Set("registry.url.port", euPort)
	viper.Set("registry.url.scheme", urlRegistry.Scheme)
	viper.Set("registry.url.hostname", urlRegistry.Hostname())
	viper.Set("registry.url.path", urlRegistry.Path)
	viper.Set("blob.url.port", bPort)
	viper.Set("blob.url.scheme", urlBlob.Scheme)
	viper.Set("blob.url.hostname", urlBlob.Hostname())
	viper.Set("blob.url.path", urlBlob.Path)

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	slog.Info("Configuration file written, exiting...")
	return nil
}
