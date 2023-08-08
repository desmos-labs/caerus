package notifications

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/caerus/server/authentication"
	"github.com/desmos-labs/caerus/server/runner"
	"github.com/desmos-labs/caerus/server/types"
	"github.com/desmos-labs/caerus/server/utils"
)

// RoutesRegistrar implements runner.RoutesRegister
func RoutesRegistrar(router *gin.Engine, ctx runner.Context) {
	Register(router, NewHandler(ctx.FirebaseClient, ctx.Database))
}

// Register allows to register all the routes related to notifications
func Register(router *gin.Engine, handler *Handler) {
	appAuthMiddleware := authentication.NewAppAuthMiddleware(handler)
	notificationsRouter := router.Group("/notifications", appAuthMiddleware)

	// ----------------------------------------
	// --- Notifications endpoints
	// ----------------------------------------

	notificationsRouter.
		POST("", func(c *gin.Context) {
			// Read the body
			body, err := c.GetRawData()
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Parse the request
			req, err := handler.ParseSendNotificationRequest(body)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			req.AppID = c.MustGet(types.SessionAppID).(string)

			// Handle the request
			err = handler.HandleSendNotificationRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			c.String(http.StatusOK, "Notification sent successfully")
		})
}
