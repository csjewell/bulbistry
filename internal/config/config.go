// The code that handles and stores the configuration for Bulbistry
package config

import (
	v "internal/version"

	"errors"
	"net/url"
	"os"

	yaml "github.com/goccy/go-yaml"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.


type ConfigURL struct {
	Scheme    string `yaml:"scheme"`
	HostName  string `yaml:"hostname"`
	Port      int    `yaml:"port,omitempty"`
	Path      string `yaml:"path"`
}

type ConfigListenOn struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

// A Config variable contains the bulbistry configuration
type Config struct {
	Version       string         `yaml:"version"`
	ExternalURL   ConfigURL      `yaml:"external_url"`
	BlobURL       ConfigURL      `yaml:"blob_url"`
	ListenOn      ConfigListenOn `yaml:"listen_on"`
	BlobIsProxied bool           `yaml:"is_proxied"`
	DatabaseFile  string         `yaml:"database_file"`
	HTPasswdFile  string         `yaml:"htpasswd_file"`
	BlobDirectory string         `yaml:"blob_directory"`
}

type bulbistryConfigError struct {
	configKey string
	error
}

func newConfigError(key, err string) bulbistryConfigError {
	return bulbistryConfigError{key, errors.New(err + ": " + key)}
}

func (u ConfigURL) getHostname() string {
	scheme := u.Scheme
	port   := u.Port
	host   := u.HostName

	if scheme == "http" && port == 80 {
		return host
	}

	if scheme == "https" && port == 443 {
		return host
	}

	return host + ":" + string(port)
}

func (u ConfigURL) getURL() *url.URL {
	return &url.URL{
		Scheme: u.Scheme,
		Host:   u.getHostname(),
		Path:   u.Path,
	}
}

func (cfg Config) GetExternalURL() *url.URL {
	return cfg.ExternalURL.getURL().JoinPath("/v2/")
}

// GetListenOn gets the IP and port that the registry is configured to listen on
func (cfg Config) GetListenOn() string {
	return cfg.ListenOn.IP + ":" + string(cfg.ListenOn.Port)
}

// SaveConfig saves the current configuration to a YAML file.
func (cfg Config) SaveConfig(filename string) error {
	cfg.Version = v.Version()

	yml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, []byte(yml), 0640)
	if err != nil {
		return err
	}

	return nil
}

// ReadConfig reads the current configuration from a YAML file
func ReadConfig(filename string) (*Config, error) {
	yml, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal([]byte(yml), &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.DatabaseFile == "" {
		return nil, newConfigError("database_file", "configuration entry required")
	}

	if cfg.ListenOn.Port == 0 {
		cfg.ListenOn.Port = 28080
	}

	if cfg.ListenOn.IP == "" {
		cfg.ListenOn.IP = "127.0.0.1"
	}

	if cfg.ExternalURL.HostName == "" {
		return nil, newConfigError("hostname", "configuration entry required")
	}

	if cfg.ExternalURL.Port == 0 {
		cfg.ExternalURL.Port = 80
	}

	if cfg.ExternalURL.Scheme == "" {
		cfg.ExternalURL.Scheme = "http"
	}

	if cfg.BlobIsProxied {
		if cfg.BlobURL.HostName == "" {
			return nil, newConfigError("hostname", "configuration entry required")
		}
	} else {
		if cfg.BlobURL.HostName == "" {
			cfg.BlobURL.HostName = cfg.ExternalURL.HostName
		}
	}

	if cfg.BlobURL.Port == 0 {
		cfg.ExternalURL.Port = 80
	}

	if cfg.ExternalURL.Scheme == "" {
		cfg.ExternalURL.Scheme = "http"
	}

	return &cfg, nil
}
