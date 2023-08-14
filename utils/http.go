package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

// GetTokenValue returns the token value associated with the given context, reading it from the Authorization header
func GetTokenValue(c *gin.Context) (string, error) {
	headerValue := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(headerValue, "Bearer"))
	if token == "" {
		return "", WrapErr(http.StatusUnauthorized, "wrong Authorization header value")
	}

	return token, nil
}

// --------------------------------------------------------------------------------------------------------------------

type HttpError struct {
	StatusCode int
	Response   string
	Headers    map[string]string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Response)
}

// WrapErr wraps the given error into a new one that contains the given status code and response
func WrapErr(statusCode int, res string) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Response:   res,
	}
}

func (e *HttpError) WithHeaders(headers map[string]string) *HttpError {
	e.Headers = headers
	return e
}

// --------------------------------------------------------------------------------------------------------------------

func NewTooManyRequestsError(res string) *HttpError {
	// Compute when the user will be able to retry
	currentTime := time.Now().UTC()
	nextDayMidnight := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, 0, 0, 0, 0, time.UTC)
	retryAfter := nextDayMidnight.Sub(currentTime).Seconds()

	return WrapErr(http.StatusTooManyRequests, res).WithHeaders(map[string]string{
		http.CanonicalHeaderKey("Retry-After"): fmt.Sprintf("%d", int(retryAfter)),
	})
}

// --------------------------------------------------------------------------------------------------------------------

func ucFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// UnwrapErr unwraps the given error returning the status code and response
func UnwrapErr(err error) (statusCode int, res string, headers map[string]string) {
	if httpErr, ok := err.(*HttpError); ok {
		return httpErr.StatusCode, httpErr.Response, httpErr.Headers
	}
	return http.StatusInternalServerError, ucFirst(err.Error()), nil
}

// --------------------------------------------------------------------------------------------------------------------

type errorJsonResponse struct {
	Error string `json:"error"`
}

// HandleError handles the given error by returning the proper response
func HandleError(c *gin.Context, err error) {
	statusCode, res, headers := UnwrapErr(err)
	c.Abort()
	c.Error(err)
	for key, value := range headers {
		c.Header(key, value)
	}
	c.JSON(statusCode, errorJsonResponse{Error: res})
}

// --------------------------------------------------------------------------------------------------------------------
