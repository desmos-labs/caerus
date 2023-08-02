package notifications

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/caerus/server/authentication"
	"github.com/desmos-labs/caerus/server/types"
	"github.com/desmos-labs/caerus/server/utils"
)

func Register(router *gin.Engine, handler *Handler) {
	authMiddleware := authentication.NewMiddleware(handler)

	notificationsRoutesGroup := router.Group("/notifications")

	// ----------------------------------------
	// --- Notifications routes
	// ----------------------------------------

	notificationsRoutesGroup.
		POST("", func(c *gin.Context) {
			// Read the body
			body, err := c.GetRawData()
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Get the token
			token, err := utils.GetTokenValue(c)
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
			req.Token = token

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

	notificationsRoutesGroup.Group("", authMiddleware).
		POST("/tokens", func(c *gin.Context) {
			// Read the body
			body, err := c.GetRawData()
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Parse the request
			req, err := handler.ParseRegisterDeviceTokenRequest(body)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			req.UserAddress = c.MustGet(types.SessionDesmosAddressKey).(string)

			// Handle the request
			err = handler.HandleRegisterDeviceTokenRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			c.String(http.StatusOK, "Device token registered successfully")
		})
}
