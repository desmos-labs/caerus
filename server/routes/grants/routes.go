package grants

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"

	"github.com/desmos-labs/caerus/server/analytics"
	"github.com/desmos-labs/caerus/server/authentication"
	"github.com/desmos-labs/caerus/server/runner"
	"github.com/desmos-labs/caerus/server/types"
	"github.com/desmos-labs/caerus/server/utils"
)

// RoutesRegistrar implements runner.RoutesRegister
func RoutesRegistrar(router *gin.Engine, ctx runner.Context) {
	Register(router, NewHandler(ctx.ChainClient, ctx.Database))
}

// Register registers the routes related to the grants module inside the given router.
func Register(router *gin.Engine, handler *Handler) {
	appAuthMiddleware := authentication.NewAppAuthMiddleware(handler)

	// ----------------------------------------
	// --- Funds endpoints
	// ----------------------------------------

	router.
		GET("/grant", appAuthMiddleware, func(c *gin.Context) {
			// Get the request body
			body, err := c.GetRawData()
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Parse the request
			req, err := handler.ParseRequestFeeGrantRequest(body)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			req.AppID = c.MustGet(types.SessionAppID).(string)

			// Handle the request
			err = handler.HandleFeeGrantRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Log the event
			analytics.Enqueue(posthog.Capture{
				DistinctId: req.AppID,
				Event:      "Requested fee grant",
				Properties: posthog.NewProperties().
					Set(analytics.KeyUserAddress, req.DesmosAddress),
			})

			c.String(http.StatusOK, "Fee grant requested successfully. You and the user will receive a notification once it has been granted")
		})
}
