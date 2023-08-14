package authentication

import (
	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

// NewUserAuthMiddleware represents a Gin middleware that allows to authenticate a request made from a user.
// After the request has been authenticated properly, the types.SessionTokenKey and
// types.SessionDesmosAddressKey context's keys will be set to contain the user details.
func NewUserAuthMiddleware(source Source) gin.HandlerFunc {
	return func(context *gin.Context) {
		// Get the token
		token, err := utils.GetTokenValue(context)
		if err != nil {
			utils.HandleError(context, err)
			return
		}

		// Verify the session
		session, err := source.GetUserSession(token)
		if err != nil {
			utils.HandleError(context, err)
			return
		}

		// Set the context variables
		context.Set(types.SessionTokenKey, token)
		context.Set(types.SessionDesmosAddressKey, session.DesmosAddress)

		// Go to the next step of the request processing
		context.Next()
	}
}

// NewAppAuthMiddleware represents a Gin middleware that allows to authenticate a request made from an application.
// After the request has been authenticated properly, the types.SessionAppID context's key will be set to
// contain the application details.
func NewAppAuthMiddleware(source Source) gin.HandlerFunc {
	return func(context *gin.Context) {
		// Get the token
		token, err := utils.GetTokenValue(context)
		if err != nil {
			utils.HandleError(context, err)
			return
		}

		// Verify the session
		appToken, err := source.GetAppToken(token)
		if err != nil {
			utils.HandleError(context, err)
			return
		}

		// Set the context variables
		context.Set(types.SessionTokenKey, token)
		context.Set(types.SessionAppID, appToken.AppID)

		// Go to the next step of the request processing
		context.Next()
	}
}
