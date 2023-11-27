package bulbistry

// This generates version 5 UUID's based on the URL.

import (
	"github.com/google/uuid"
)

func GenerateManifestUUID(bc BulbistryConfig, manifestName string, reference string) (string, error) {

	url := bc.GetExternalUrl()
	url = url + manifestName + "/manifests/" + reference

	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(url)).String(), nil
}

func GenerateBlobUUID(bc BulbistryConfig, manifestName string, reference string) (string, error) {

	url := bc.GetExternalUrl()
	url = url + manifestName + "/blobs/" + reference

	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(url)).String(), nil
}
