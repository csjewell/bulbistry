package bulbistry

// This generates version 5 UUID's based on the URL.

import (
	"internal/config"
	"internal/database"
	"internal/urls"

	"github.com/google/uuid"
)

func GenerateManifestUUID(cfg config.Config, mt database.ManifestTag) (string, error) {
	url := urls.GetManifestURL(cfg, mt)

	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(url)).String(), nil
}

func GenerateBlobUUID(cfg config.Config, manifestName string, reference string) (string, error) {

//	url := cfg.GetExternalURL()
//	url = url + manifestName + "/blobs/" + reference

//	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(url)).String(), nil
return "", nil
}
