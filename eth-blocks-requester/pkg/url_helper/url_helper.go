package url_helper

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

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

func GetRandId() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := math.MaxInt - 100
	randomInt := rand.Intn(max-min+1) + min
	return randomInt
}
