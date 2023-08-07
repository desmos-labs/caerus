package notifications

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/caerus/server/authentication"
	"github.com/desmos-labs/caerus/server/types"
	"github.com/desmos-labs/caerus/server/utils"
)

// Register allows to register all the routes related to notifications
func Register(router *gin.Engine, handler *Handler) {
	appAuthMiddleware := authentication.NewAppAuthMiddleware(handler)
	userAuthMiddleware := authentication.NewUserAuthMiddleware(handler)

	notificationsRoutesGroup := router.Group("/notifications")

	// ----------------------------------------
	// --- Notifications routes
	// ----------------------------------------

	notificationsRoutesGroup.
		POST("", appAuthMiddleware, func(c *gin.Context) {
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

	// ----------------------------------------
	// --- Device token routes
	// ----------------------------------------

	notificationsRoutesGroup.
		POST("/app-tokens", appAuthMiddleware, func(c *gin.Context) {
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
		POST("/user-tokens", userAuthMiddleware, func(c *gin.Context) {
			// Read the body
			body, err := c.GetRawData()
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Parse the request
			req, err := handler.ParseRegisterUserDeviceTokenRequest(body)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			req.UserAddress = c.MustGet(types.SessionDesmosAddressKey).(string)

			// Handle the request
			err = handler.HandleRegisterUserDeviceTokenRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			c.String(http.StatusOK, "Device token registered successfully")
		})
}
