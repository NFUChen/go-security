package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"strconv"
)

func GetUserIdFromQueryParam(ctx echo.Context) (uint, error) {
	stringUserID := ctx.QueryParam("user_id")
	if stringUserID == "" {
		return 0, fmt.Errorf("user_id is required")
	}

	userID, err := strconv.Atoi(stringUserID)
	if err != nil {
		return 0, fmt.Errorf("user_id must be a number")
	}

	return uint(userID), nil
}
