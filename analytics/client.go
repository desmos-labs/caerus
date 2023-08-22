package analytics

import (
	"strconv"
	"strings"

	"github.com/posthog/posthog-go"

	"github.com/desmos-labs/caerus/utils"
)

var client posthog.Client

func init() {
	enabledStr := utils.GetEnvOr(EnvAnalyticsEnabled, "false")
	enabled, err := strconv.ParseBool(enabledStr)
	if err != nil || !enabled {
		return
	}

	apiKey := utils.GetEnvOr(EnvAnalyticsPostHogApiKey, "")
	if len(strings.TrimSpace(apiKey)) == 0 {
		return
	}

	postHogClient, err := posthog.NewWithConfig(apiKey, posthog.Config{
		Endpoint: utils.GetEnvOr(EnvAnalyticsPostHogHost, "https://eu.posthog.com"),
	})
	if err != nil {
		panic(err)
	}
	client = postHogClient
}

// Enqueue allows to enqueue the given event to the PostHog instance
func Enqueue(event posthog.Capture) {
	if client != nil {
		client.Enqueue(event)
	}
}

// Stop closes the PostHog client
func Stop() {
	if client != nil {
		client.Close()
	}
}
