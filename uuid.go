package bulbistry

// This generates version 5 UUID's based on the URL.

import (
	"github.com/google/uuid"
)

func GenerateManifestUUID(cfg Config, manifestName string, reference string) (string, error) {

	url := cfg.GetExternalURL()
	url = url + manifestName + "/manifests/" + reference

	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(url)).String(), nil
}

func GenerateBlobUUID(cfg Config, manifestName string, reference string) (string, error) {

	url := cfg.GetExternalURL()
	url = url + manifestName + "/blobs/" + reference

	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(url)).String(), nil
}
