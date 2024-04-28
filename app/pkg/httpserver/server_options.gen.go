// Code generated by options-gen. DO NOT EDIT.
package httpserver

import (
	fmt461e464ebed9 "fmt"
	"log/slog"
	"time"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	logger *slog.Logger,
	isProduction bool,
	address string,
	readHeaderTimeout time.Duration,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.logger = logger
	o.isProduction = isProduction
	o.address = address
	o.readHeaderTimeout = readHeaderTimeout

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("logger", _validate_Options_logger(o)))
	return errs.AsError()
}

func _validate_Options_logger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.logger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `logger` did not pass the test: %w", err)
	}
	return nil
}
