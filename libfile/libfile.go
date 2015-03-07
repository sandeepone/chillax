package libfile

import (
	"strings"
)

var MimeTypesMap = map[string]string{
	".txt":  "text/plain; charset=utf-8",
	".xml":  "text/xml; charset=utf-8",
	".css":  "text/css; charset=utf-8",
	".htm":  "text/html; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".gif":  "image/gif",
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".js":   "application/x-javascript",
	".pdf":  "application/pdf",
}

func GetExtensionByMime(mimeToCheck string) string {
	for ext, mimeType := range MimeTypesMap {
		if strings.HasPrefix(mimeType, mimeToCheck) {
			return ext
		}
	}

	return ""
}

func GetMimeByExtension(extensionToCheck string) string {
	for ext, mimeType := range MimeTypesMap {
		if strings.HasSuffix(ext, extensionToCheck) {
			return mimeType
		}
	}

	return ""
}
