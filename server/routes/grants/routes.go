package grants

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"

	"github.com/desmos-labs/caerus/server/analytics"
	"github.com/desmos-labs/caerus/server/authentication"
	"github.com/desmos-labs/caerus/server/types"
	"github.com/desmos-labs/caerus/server/utils"
)

func Register(router *gin.Engine, handler *Handler) {
	authMiddleware := authentication.NewMiddleware(handler)

	// ----------------------------------------
	// --- Funds endpoints
	// ----------------------------------------

	router.
		GET("/authorizations", authMiddleware, func(c *gin.Context) {
			// Parse the request
			token := c.MustGet(types.SessionTokenKey).(string)

			err := handler.HandleFeeGrantRequest(token)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Log the event
			analytics.Enqueue(posthog.Capture{
				DistinctId: c.MustGet(types.SessionDesmosAddressKey).(string),
				Event:      "Requested Authorizations",
			})

			c.String(http.StatusOK, "Authorizations requested successfully. You will receive a notification once they have been approved")
		})
}
