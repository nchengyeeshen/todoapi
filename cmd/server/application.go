package main

import "log/slog"

type Application struct {
	logger *slog.Logger
	repo   TodoRepository
}

// NewApplication returns a new [Application].
func NewApplication(logger *slog.Logger, repo TodoRepository) *Application {
	return &Application{
		logger: logger,
		repo:   repo,
	}
}
