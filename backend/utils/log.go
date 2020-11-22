package utils

import (
	"github.com/labstack/echo/v4"
)

func Log(logger echo.Logger, message string, do func() error) error {
	logger.Debug(message + "...")
	err := do()
	if err != nil {
		logger.Error(err, message+"... error!")
	} else {
		logger.Debug(message + "... done!")
	}
	return err
}
