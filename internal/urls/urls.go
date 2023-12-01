// The code that creates URLS for bulbistry, sans-router.
package urls

import (
	"internal/config"
	"internal/database"
)

// GetManifestURL gets the URL to retrieve a particular manifest
func GetManifestURL(cfg config.Config, mt database.ManifestTag) string {
	if mt.Namespace == "" {
		return cfg.GetExternalURL().JoinPath(mt.Name, "/manifest/", mt.Sha512).String()
	}
	return cfg.GetExternalURL().JoinPath(mt.Namespace, mt.Name, "/manifest/", mt.Sha512).String()
}

// GetBlobURL gets the blob storage base URL.
func GetBlobURL(cfg config.Config) string {
	return cfg.ExternalURL.Scheme + "://" + cfg.ExternalURL.HostName + ":" + string(cfg.ExternalURL.Port) + "/v2/"
}
