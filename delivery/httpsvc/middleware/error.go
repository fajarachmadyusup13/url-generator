package middleware

import "github.com/gin-gonic/gin"

// CustomHTTPError is a custom error struct that implements the error interface.
type CustomHTTPError struct {
	Code    int
	Message string
}

func (e CustomHTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new CustomHTTPError instance.
func NewHTTPError(code int, message string) error {
	return CustomHTTPError{
		Code:    code,
		Message: message,
	}
}

// Custom error handling middleware
func CustomErrorMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		// Retrieve the last error from the context
		lastError := c.Errors.Last()

		if err, ok := lastError.Err.(CustomHTTPError); ok {
			c.JSON(err.Code, gin.H{
				"error": err.Message,
			})
			c.Abort()
		}
	}
}
