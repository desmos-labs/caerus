package files

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/caerus/server/authentication"
	"github.com/desmos-labs/caerus/server/utils"
)

// Register registers all the request handlers related to the media functionalities
func Register(router *gin.Engine, handler *Handler) {
	authMiddleware := authentication.NewMiddleware(handler)

	router.Group("/media").

		// ----------------------------------------
		// --- Authenticated routes
		// ----------------------------------------

		POST("", authMiddleware, func(c *gin.Context) {
			// Parse the request
			filePath, err := handler.GetFileFromRequest(c.Request)
			if err != nil {
				utils.HandleError(c, err)
				return
			}

			// Handle the request
			tempFilePath, res, err := handler.UploadFile(filePath)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			os.Remove(tempFilePath)

			c.JSON(http.StatusOK, &res)
		}).

		// ----------------------------------------
		// --- Unauthenticated routes
		// ----------------------------------------

		GET("/:fileName", func(c *gin.Context) {
			// Parse the request
			fileName := c.Param("fileName")

			// Handle the request
			tempFilePath, err := handler.GetFile(fileName)
			if err != nil {
				utils.HandleError(c, err)
				return
			}
			defer os.Remove(tempFilePath)

			// Serve the file
			c.File(tempFilePath)
		})

}
