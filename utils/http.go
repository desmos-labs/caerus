package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"unicode"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// --------------------------------------------------------------------------------------------------------------------

type HttpError struct {
	StatusCode codes.Code
	Response   string
	Headers    map[string]string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Response)
}

// WrapErr wraps the given error into a new one that contains the given status code and response
func WrapErr(statusCode codes.Code, res string) *HttpError {
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

func UnwrapError(ctx context.Context, err error) error {
	if httpErr, ok := err.(*HttpError); ok {
		if len(httpErr.Headers) != 0 {
			err := grpc.SendHeader(ctx, metadata.New(httpErr.Headers))
			if err != nil {
				return err
			}
		}

		return status.Error(status.Code(err), httpErr.Response)
	}
	return status.Error(codes.Internal, ucFirst(err.Error()))
}
