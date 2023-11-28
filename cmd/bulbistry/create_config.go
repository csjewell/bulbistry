package main

import (
	"errors"
	"fmt"
	//"net/url"
	"os/user"
	"path"
	"strings"

	tbv "github.com/csjewell/bulbistry"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func CreateConfig(ctx *cli.Context) error {

	//type BulbistryConfig_URL struct {
	//	Scheme    string `yaml:"scheme"`
	//	HostName  string `yaml:"hostname"`
	//	Port      int    `yaml:"port"`
	//	Directory string `yaml:"dir"`
	//}

	//type BulbistryConfig_ListenOn struct {
	//	IP   string `yaml:"ip"`
	//	Port int    `yaml:"port"`
	//}

	//type BulbistryConfig struct {
	//	ExternalUrl   BulbistryConfig_URL      `yaml:"external_url,inline"`
	//	BlobUrl       BulbistryConfig_URL      `yaml:"blob_url,inline"`
	//	ListenOn      BulbistryConfig_ListenOn `yaml:"listen_on,inline"`
	//	BlobIsProxied bool                     `yaml:"is_proxied"`
	//	DatabaseFile  string                   `yaml:"database_file"`
	//	HTPasswdFile  string                   `yaml:"htpasswd_file"`
	//	BlobDirectory string                   `yaml:"blob_directory"`
	//}

	validateFile := func(input string) error {
		if input == "" {
			return errors.New("Filename cannot be empty")
		}
		// TODO: Other checks (exists as a directory, etc.)
		return nil
	}

	prompt := promptui.Prompt{
		Label:   "What is the registry hostname",
		Default: "registry.local",
		// Validate:  validateFile,
		AllowEdit: true,
	}

	hostName, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get hostname %v\n", err)
	}

	prompt = promptui.Prompt{
		Label:   "What is the URL for the registry",
		Default: "https://" + hostName + "/",
		// Validate:  validateDirectory,
		AllowEdit: true,
	}

	//url, err := prompt.Run()
	_, err = prompt.Run()

	if err != nil {
		return fmt.Errorf("Did not get registry URL %v\n", err)
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

	menu := promptui.Select{
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

	port, _ := strings.Atoi(portStr)

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

	menu = promptui.Select{
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

	config = tbv.BulbistryConfig{
		ListenOn:      tbv.BulbistryConfig_ListenOn{IP: ip, Port: port},
		HTPasswdFile:  htpasswdFile,
		DatabaseFile:  databaseFile,
		BlobDirectory: blobDirectory,
		BlobIsProxied: isProxied,
	}

	// config.Save()

	return fmt.Errorf("Not Implemented")
}
