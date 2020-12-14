package handler

import (
	"context"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/router"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func GetEchoContext(c echo.Context, handler string) (context.Context, *zap.Logger) {
	ctx := GetContext(c)
	l := logging.WithContext(ctx).With(zap.String("handler", handler)).Named("handler." + handler)
	return ctx, l
}

var HttpErrorHandler = func(err error, c echo.Context) {

	_, l := GetEchoContext(c, "")
	l.Error(err.Error(), zap.Error(err))

	if ws, ok := err.(*errors.WebServiceException); ok {
		c.JSON(ws.Status, &errors.ErrorResponse{
			Message:    ws.Message,
			Code:       ws.Code,
			StatusCode: ws.Status,
		})
		return
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		validationError := err.(validator.ValidationErrors)

		var validErrors []errors.ValidError
		for _, fieldError := range validationError {
			validErrors = append(validErrors, errors.ValidError{
				Tag:             fieldError.Tag(),
				ActualTag:       fieldError.ActualTag(),
				Namespace:       fieldError.Namespace(),
				StructNamespace: fieldError.StructNamespace(),
				Field:           fieldError.Field(),
				StructField:     fieldError.StructField(),
				Param:           fieldError.Param(),
				Kind:            fieldError.Kind().String(),
				Type:            fieldError.Type().String(),
			})
		}

		translated := validationError.Translate(router.DefaultTranslator)

		response := &errors.ErrorResponse{
			Message:    validationError.Error(),
			Code:       "ErrValidation",
			StatusCode: http.StatusBadRequest,
			Validation: translated,
			Errors:     validErrors,
		}

		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusInternalServerError, &errors.ErrorResponse{
		Message:    "Internal server error",
		Code:       "ErrInternalServerError",
		StatusCode: http.StatusInternalServerError,
	})

}
