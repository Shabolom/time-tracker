package handlers

import (
	"github.com/go-playground/validator/v10"
	"timeTracker/config"
	"timeTracker/internal/storage"
	"timeTracker/pkg/httputils"
)

// Handlers implements all the handler functions and has the dependencies that they use
type Handlers struct {
	Sender  *httputils.Sender
	Storage storage.StorageInterface
	EnvBox  config.ApiEnvConfig
}

// Validate is a singleton that provides validation services for in handlers.
var Validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
