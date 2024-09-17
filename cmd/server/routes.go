package main

import "net/http"

type routeMetadata struct {
	Name string
}

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	reqIDMW := RequestIDMiddleware()
	logMW := LoggingMiddleware(app.logger)

	handle := func(pattern string, hnd http.Handler, meta routeMetadata) {
		mux.Handle(pattern, reqIDMW(logMW(hnd, meta), meta))
	}

	handle(
		"GET /healthcheck",
		app.handleHealthcheck(),
		routeMetadata{
			Name: "healthcheck",
		},
	)

	handle(
		"GET /api/todos",
		app.handleGetAllTodos(),
		routeMetadata{
			Name: "getAllTodos",
		},
	)

	handle(
		"POST /api/todos",
		app.handleCreateTodo(),
		routeMetadata{
			Name: "createTodo",
		},
	)

	handle(
		"GET /api/todos/{id}",
		app.handleGetTodo(),
		routeMetadata{
			Name: "getTodo",
		},
	)

	handle(
		"PUT /api/todos/{id}",
		app.handleUpdateTodo(),
		routeMetadata{
			Name: "updateTodo",
		},
	)

	handle(
		"DELETE /api/todos/{id}",
		app.handleDeleteTodo(),
		routeMetadata{
			Name: "deleteTodo",
		},
	)

	return mux
}
