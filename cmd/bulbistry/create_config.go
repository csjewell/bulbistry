package main

import (
	"internal/config"

	"errors"
	"fmt"
	"net/url"
	"os/user"
	"path"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func CreateConfig(ctx *cli.Context) error {

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

	urlBlobDefault := urlRegistry.JoinPath("/blob")

	prompt = promptui.Prompt{
		Label:     "Where is the URL to where blobs are stored",
		Default:   urlBlobDefault.String(),
		AllowEdit: true,
		// Validate:  validateURL,
	}

	urlString, err = prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get registry URL %v\n", err)
	}

	urlBlob, _ := url.Parse(urlString)

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

	menu = promptui.Select{
		Label:     "Where should the registry listen",
		CursorPos: 0,
		Items:     []string{"127.0.0.1", "0.0.0.0"},
	}

	_, ip, err := menu.Run()
	if err != nil {
		return fmt.Errorf("Did not get IP %v\n", err)
	}

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

	euPort, _ := strconv.Atoi(urlRegistry.Port())
	bPort, _ := strconv.Atoi(urlBlob.Port())

	cfg = config.Config{
		ExternalURL:   config.ConfigURL{
			Scheme:   urlRegistry.Scheme,
			HostName: urlRegistry.Hostname(),
			Port:     euPort,
			Path:     urlRegistry.Path,
		},
		BlobURL:       config.ConfigURL{
			Scheme:   urlBlob.Scheme,
			HostName: urlBlob.Hostname(),
			Port:     bPort,
			Path:     urlBlob.Path,
		},
		ListenOn:      config.ConfigListenOn{IP: ip, Port: port},
		HTPasswdFile:  htpasswdFile,
		DatabaseFile:  databaseFile,
		BlobDirectory: blobDirectory,
		BlobIsProxied: isProxied,
	}

	err = cfg.SaveConfig(ctx.String("config"))
	if err != nil {
		return err
	}

	return fmt.Errorf("Configuration file written, exiting...")
}
