package bulbistry

// This generates version 5 UUID's based on the URL.

import(
	"github.com/google/uuid"
)

func generateManifestUUID(bc BulbistryConfig, manifestName string, reference string) (uuid.UUID, error) {

	url := bc.GetExternalUrl();
	url = url + manifestName + "/manifests/" + reference 

    return uuid.NewSHA1(uuid.NameSpaceURL, []byte(url)), nil
}