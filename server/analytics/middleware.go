package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"

	"github.com/desmos-labs/caerus/server/types"
)

// PostHog represents a Gin middleware that allows to send analytics data to the PostHog instance
func PostHog() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Go to the next step of the request processing
		context.Next()

		// Add the user identity
		userAddress, exists := context.Get(types.SessionDesmosAddressKey)
		if !exists {
			context.Next()
			return
		}

		distinctID := userAddress.(string)

		// Capture the page view
		Enqueue(posthog.Capture{
			DistinctId: distinctID,
			Event:      "$pageview",
			Properties: posthog.NewProperties().
				Set(KeyContentURL, context.FullPath()),
		})
	}
}
