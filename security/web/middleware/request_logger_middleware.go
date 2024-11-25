package web

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
)

func RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Save the original body reader
		req := c.Request()
		body := req.Body

		// Read the body
		var buf bytes.Buffer
		if body != nil {
			_, err := io.Copy(&buf, body)
			if err != nil {
				return err
			}
		}

		// Log the body
		log.Info().Msgf("Request Body: %s", buf.String())

		// Reset the body so it can be read again
		req.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

		return next(c)
	}
}
