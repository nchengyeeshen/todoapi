package main

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	codeBadRequest = "bad_request"
	codeNotFound   = "not_found"
)

func (app *Application) serverError(ctx context.Context, w http.ResponseWriter, err error) {
	app.logger.WarnContext(ctx, "Server error", "err", err)
	app.jsonResponse(ctx, w, http.StatusInternalServerError, map[string]any{
		"code":        "internal",
		"description": "Internal server error. Try again later.",
	})
}

func (app *Application) clientError(
	ctx context.Context,
	w http.ResponseWriter,
	statusCode int,
	errCode, description string,
) {
	app.jsonResponse(ctx, w, statusCode, map[string]any{
		"code":        errCode,
		"description": description,
	})
}

func (app *Application) jsonResponse(ctx context.Context, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		app.logger.WarnContext(ctx, "Write response", "err", err)
	}
}
