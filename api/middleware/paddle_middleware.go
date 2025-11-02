package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/gin-gonic/gin"
)

func PaddleWebhookVerifier(secretKey string) gin.HandlerFunc {
	verifier := paddle.NewWebhookVerifier(secretKey)

	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		verifyReq := c.Request.Clone(c.Request.Context())
		verifyReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		verified := false
		handler := verifier.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			verified = true
		}))

		responseRecorder := &responseWriter{ResponseWriter: c.Writer, statusCode: http.StatusOK}
		handler.ServeHTTP(responseRecorder, verifyReq)

		if !verified || responseRecorder.statusCode != http.StatusOK {
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		c.Next()
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
