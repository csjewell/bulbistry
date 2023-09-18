package bulbistry

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type BulbistryConfig struct {
	ExternalUrl struct {
	    Scheme string    `yaml:"scheme"`
	    HostName string  `yaml:"hostname"`
	    Port int         `yaml:"port"`
	}
	ListenOn struct {
	    IP string     `yaml:"ip"`
	    Port int      `yaml:"port"`
	}
	DatabaseFile string `yaml:"database_file"`
	HTPasswdFile string `yaml:"htpasswd_file"`
}

func (bc BulbistryConfig) GetExternalUrl() string {
	return bc.ExternalUrl.Scheme + "://" + bc.ExternalUrl.HostName + ":" + string(bc.ExternalUrl.Port) + "/v2/"
}

func (bc BulbistryConfig) GetListenOn() string {
	return bc.ListenOn.IP + ":" + string(bc.ListenOn.Port)
}

func readConfig (filename string) (*BulbistryConfig, error) {
    yml, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var bc BulbistryConfig
    err = yaml.Unmarshal([]byte(yml), &bc)
    if err != nil {
	return nil, err
    }

    return &bc, nil
}
