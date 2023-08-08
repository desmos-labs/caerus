package applications

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
	Register(router, NewHandler(ctx.Database))
}

// Register allows to register all the routes related to applications
func Register(router *gin.Engine, handler *Handler) {
	appAuthMiddleware := authentication.NewAppAuthMiddleware(handler)

	applicationRouter := router.Group("/applications")

	// ----------------------------------------
	// --- Application endpoints
	// ----------------------------------------

	applicationRouter.Group("/me", appAuthMiddleware).
		POST("/notification-tokens", func(c *gin.Context) {
			// Read the body
			body, err := c.GetRawData()
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Parse the request
			req, err := handler.ParseRegisterAppDeviceTokenRequest(body)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			req.AppID = c.MustGet(types.SessionAppID).(string)

			// Handle the request
			err = handler.HandleRegisterAppDeviceTokenRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			c.String(http.StatusOK, "Device token registered successfully")
		}).
		DELETE("", func(c *gin.Context) {
			// Parse the request
			req := &DeleteApplicationRequest{
				AppID: c.MustGet(types.SessionAppID).(string),
			}

			// Handle the request
			err := handler.HandleDeleteApplicationRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			c.String(http.StatusOK, "Application deleted successfully")
		})
}
