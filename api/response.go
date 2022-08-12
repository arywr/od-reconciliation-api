package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func APIResponse(code int, status string, data interface{}) Response {
	response := Response{
		Code:   code,
		Status: status,
		Data:   data,
	}

	return response
}

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type ErrorValidationResponse struct {
	Code   int               `json:"code"`
	Status string            `json:"status"`
	Error  []ValidationError `json:"data"`
}

type ErrorResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Error  error  `json:"error"`
}

func APIErrorResponse(code int, status string, err error) ErrorResponse {
	response := ErrorResponse{
		Code:   code,
		Status: status,
		Error:  err,
	}

	return response
}

func APIValidationResponse(code int, status string, verr validator.ValidationErrors) ErrorValidationResponse {
	errs := []ValidationError{}

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, ValidationError{Field: f.Field(), Reason: err})
	}

	response := ErrorValidationResponse{
		Code:   code,
		Status: status,
		Error:  errs,
	}

	return response
}
