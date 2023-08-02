package authentication

import (
	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/caerus/server/types"
	"github.com/desmos-labs/caerus/server/utils"
)

// NewMiddleware represents a Gin middleware that allows to authenticate a request.
// After the request has been authenticated properly, the types.SessionTokenKey and
// types.SessionDesmosAddressKey context's keys will be set to contain the user details.
func NewMiddleware(source Source) gin.HandlerFunc {
	return func(context *gin.Context) {
		// Get the token
		token, err := utils.GetTokenValue(context)
		if err != nil {
			utils.HandleError(context, err)
			return
		}

		// Verify the session
		session, err := source.GetSession(token)
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
