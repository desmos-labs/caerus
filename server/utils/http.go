package utils

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type HttpError struct {
	StatusCode int
	Response   string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Response)
}

// WrapErr wraps the given error into a new one that contains the given status code and response
func WrapErr(statusCode int, res string) error {
	return &HttpError{
		StatusCode: statusCode,
		Response:   res,
	}
}

func ucFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// UnwrapErr unwraps the given error returning the status code and response
func UnwrapErr(err error) (statusCode int, res string) {
	if httpErr, ok := err.(*HttpError); ok {
		return httpErr.StatusCode, httpErr.Response
	}
	return http.StatusInternalServerError, ucFirst(err.Error())
}

type errorJsonResponse struct {
	Error string `json:"error"`
}

// HandleError handles the given error by returning the proper response
func HandleError(c *gin.Context, err error) {
	statusCode, res := UnwrapErr(err)
	c.Abort()
	c.Error(err)
	c.JSON(statusCode, errorJsonResponse{Error: res})
}

// GetTokenValue returns the token value associated with the given context, reading it from the Authorization header
func GetTokenValue(c *gin.Context) (string, error) {
	headerValue := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(headerValue, "Bearer"))
	if token == "" {
		return "", WrapErr(http.StatusUnauthorized, "wrong Authorization header value")
	}

	return token, nil
}
