package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"

	"github.com/desmos-labs/caerus/server/analytics"

	"github.com/desmos-labs/caerus/server/authentication"
	"github.com/desmos-labs/caerus/server/types"
	"github.com/desmos-labs/caerus/server/utils"
)

// Register registers all the routes that allow to perform user-related operations
func Register(router *gin.Engine, handler *Handler) {
	authMiddleware := authentication.NewUserAuthMiddleware(handler)

	// ----------------------------------------
	// --- Login endpoints
	// ----------------------------------------

	router.
		GET("/nonce/:address", func(c *gin.Context) {
			// Parse the request
			req := NewNonceRequest(c.Param("address"))

			// Handle the request
			res, err := handler.HandleNonceRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Log the event
			analytics.Enqueue(posthog.Capture{
				DistinctId: req.DesmosAddress,
				Event:      "Requested Nonce",
			})

			c.JSON(http.StatusOK, &res)
		}).
		POST("/login", func(c *gin.Context) {
			// Read the body
			body, err := c.GetRawData()
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Parse the request
			req, err := handler.ParseAuthenticateRequest(body)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Handle the request
			res, err := handler.HandleAuthenticationRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Log the request
			analytics.Enqueue(posthog.Capture{
				DistinctId: req.DesmosAddress,
				Event:      "Logged In",
			})

			c.JSON(http.StatusOK, &res)
		})

	// ----------------------------------------
	// --- Logout endpoints
	// ----------------------------------------

	router.POST("/logout", authMiddleware, func(c *gin.Context) {
		// Parse the request
		all, found := c.GetQuery("all")
		shouldLogoutAll := found && all == "true"

		req := &LogoutRequest{
			Token:            c.MustGet(types.SessionTokenKey).(string),
			LogoutAllDevices: shouldLogoutAll,
		}

		// Handle the request
		err := handler.HandleLogoutRequest(req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.String(http.StatusOK, "Logout successful")
	})

	// ----------------------------------------
	// --- User endpoints
	// ----------------------------------------

	router.Group("/me", authMiddleware).
		POST("/sessions", authMiddleware, func(c *gin.Context) {
			// Parse the request
			token := c.MustGet(types.SessionTokenKey).(string)

			// Handle the request
			err := handler.HandleRefreshSessionRequest(token)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			c.String(http.StatusOK, token)
		}).
		POST("/notification-tokens", func(c *gin.Context) {
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
		}).
		DELETE("", func(c *gin.Context) {
			// Parse the request
			req := &DeleteAccountRequest{
				UserAddress: c.MustGet(types.SessionDesmosAddressKey).(string),
			}

			// Handle the request
			err := handler.HandleDeleteAccountRequest(req)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Log the event
			analytics.Enqueue(posthog.Capture{
				DistinctId: req.UserAddress,
				Event:      "Started Account Deletion",
			})

			c.String(http.StatusOK, "Your account has been setup for deletion. It will complete in 14 days. You can cancel this anytime by logging in again during this time frame")
		})

	// ----------------------------------------
	// --- Hasura endpoints
	// ----------------------------------------

	router.
		GET("/hasura-session", func(c *gin.Context) {
			// This endpoint is different from others: if any error is raised
			// then status code 200 is returned, but with "X-Hasura-User-Role" set to "anonymous".
			// This is to adapt to Hasura Authentication specs:
			// https://hasura.io/docs/latest/auth/authentication/webhook/

			// Parse the request
			token, err := utils.GetTokenValue(c)
			if err != nil {
				c.JSON(http.StatusOK, handler.GetUnauthorizedHasuraSession())
				return
			}

			// Handle the request
			res, err := handler.HandleHasuraSessionRequest(token)
			if err != nil {
				c.JSON(http.StatusOK, handler.GetUnauthorizedHasuraSession())
				return
			}

			c.JSON(http.StatusOK, res)
		})
}
