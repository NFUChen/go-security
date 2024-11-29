package web

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

func RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		url := ctx.Request().URL.Path

		if strings.Contains(url, "upload") {
			return next(ctx)
		}

		// Save the original body reader
		req := ctx.Request()
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

		return next(ctx)
	}
}
