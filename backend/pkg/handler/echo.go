package handler

import (
	"context"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

func GetEchoContext(c echo.Context, handler string) (context.Context, *zap.Logger) {
	ctx := GetContext(c)
	l := logging.
		WithContext(ctx).
		With(zap.String("handler", handler)).
		Named("handler." + handler)
	return ctx, l
}

var HttpErrorHandler = func(err error, c echo.Context) {

	_, l := GetEchoContext(c, "ErrorHandler")
	l.Error("", zap.Error(err))

	var wse *exceptions.WebServiceException
	if errors.As(err, &wse) {
		c.JSON(wse.Status, &exceptions.ErrorResponse{
			Message:    wse.Message,
			Code:       wse.Code,
			StatusCode: wse.Status,
		})
		return
	}

	if httpErr, ok := err.(*echo.HTTPError); ok {
		c.JSON(httpErr.Code, &exceptions.WebServiceException{
			Status: httpErr.Code,
		})
		return
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		validationError := err.(validator.ValidationErrors)

		var validErrors []exceptions.ValidError
		for _, fieldError := range validationError {
			validErrors = append(validErrors, exceptions.ValidError{
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

		translated := validationError.Translate(validation.DefaultTranslator)

		response := &exceptions.ErrorResponse{
			Message:    validationError.Error(),
			Code:       "ErrValidation",
			StatusCode: http.StatusBadRequest,
			Validation: translated,
			Errors:     validErrors,
		}

		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusInternalServerError, &exceptions.ErrorResponse{
		Message:    "Internal server error",
		Code:       "ErrInternalServerError",
		StatusCode: http.StatusInternalServerError,
	})

}
