package bulbistry

import (
	"error"
	"os"

	"gopkg.in/yaml.v3"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type BulbistryConfig_URL struct {
	Scheme    string `yaml:"scheme"`
	HostName  string `yaml:"hostname"`
	Port      int    `yaml:"port"`
	Directory string `yaml:"dir"`
}

type BulbistryConfig_ListenOn struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type BulbistryConfig struct {
	ExternalUrl BulbistryConfig_URL `yaml:"external_url,inline"`
	BlobUrl     BulbistryConfig_URL `yaml:"blob_url,inline"`
	ListenOn    BulbistryConfig_ListenOn `yaml:"listen_on,inline"`
	BlobIsProxied bool `yaml:"is_proxied"`
	DatabaseFile string `yaml:"database_file"`
	HTPasswdFile string `yaml:"htpasswd_file"`
	BlobDirectory string `yaml:"blob_directory"`
}

type bulbistryConfigError struct {
	configKey string
	*error.Error
}

func NewConfigError(key, err string) bulbistryConfigError {
	return &bulbistryConfigError{ key, error.New(err + ": " + key) }
}

func (bc BulbistryConfig) GetExternalUrl() string {
	return bc.ExternalUrl.Scheme + "://" + bc.ExternalUrl.HostName + ":" + string(bc.ExternalUrl.Port) + "/v2/"
}

func (bc BulbistryConfig) GetBlobUrl() string {
	return bc.ExternalUrl.Scheme + "://" + bc.ExternalUrl.HostName + ":" + string(bc.ExternalUrl.Port) + "/v2/"
}

func (bc BulbistryConfig) GetListenOn() string {
	return bc.ListenOn.IP + ":" + string(bc.ListenOn.Port)
}

func ReadConfig(filename string) (*BulbistryConfig, error) {
	yml, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var bc BulbistryConfig
	err = yaml.Unmarshal([]byte(yml), &bc)
	if err != nil {
		return nil, err
	}

	if bc.DatabaseFile == "" {
		return nil, NewConfigError("database_file", "configuration entry required")
	}

	if bc.ListenOn.Port == 0 {
		bc.ListenOn.Port = 28080
	}

	if bc.ListenOn.IP == 0 {
		bc.ListenOn.IP = "127.0.0.1"
	}

	if bc.ExternalUrl.HostName == "" {
		return nil, NewConfigError("hostname", "configuration entry required")
	}

	if bc.ExternalUrl.Port == 0 {
		bc.ExternalUrl.Port = 80
	}

	if bc.ExternalUrl.Scheme == "" {
		bc.ExternalUrl.Scheme = "http"
	}

	if bc.BlobIsProxied {
		if bc.BlobUrl.HostName == "" {
			return nil, NewConfigError("hostname", "configuration entry required")
		}
	} else {
		if bc.BlobUrl.HostName == "" {
			bc.BlobUrl.HostName = bc.ExternalUrl.HostName
		}
	}

	if bc.BlobUrl.Port == 0 {
		bc.ExternalUrl.Port = 80
	}

	if bc.ExternalUrl.Scheme == "" {
		bc.ExternalUrl.Scheme = "http"
	}

	return &bc, nil
}
