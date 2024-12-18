package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

func MultiPartFileToOsFile(src *multipart.FileHeader) (*os.File, error) {
	targetFile, err := os.Create(src.Filename)
	if err != nil {
		return nil, err
	}
	sourceFile, err := src.Open()
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return nil, err
	}
	// after copying the file, we need to reset the file pointer to the beginning of the file
	if _, err := targetFile.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	return targetFile, nil

}

func getQueryParamAsInt(ctx echo.Context, paramName string) (uint, error) {
	param := ctx.QueryParam(paramName)
	if param == "" {
		return 0, fmt.Errorf("%s is required", paramName)
	}

	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", paramName)
	}

	return uint(paramInt), nil
}

func GetUserIdFromQueryParam(ctx echo.Context) (uint, error) {
	return getQueryParamAsInt(ctx, "user_id")
}

func GetCategoryIdFromQueryParam(ctx echo.Context) (uint, error) {
	return getQueryParamAsInt(ctx, "category_id")
}

func GetProductIdFromQueryParam(ctx echo.Context) (uint, error) {
	return getQueryParamAsInt(ctx, "product_id")
}
