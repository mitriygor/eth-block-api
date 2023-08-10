package url_helper

import "strings"

func GetUrl(endpoint string, version string, key string) string {
	var builder strings.Builder

	builder.WriteString(endpoint)

	if version != "" {
		builder.WriteString("/")
		builder.WriteString(version)
	}

	if key != "" {
		builder.WriteString("/")
		builder.WriteString(key)
	}

	url := builder.String()

	return url
}
