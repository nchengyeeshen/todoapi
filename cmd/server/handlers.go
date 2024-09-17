package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *Application) handleHealthcheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK!"))
	})
}

func (app *Application) handleGetAllTodos() http.Handler {
	type todo struct {
		ID          string    `json:"id"`
		Status      string    `json:"status"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
	type response struct {
		Todos []todo `json:"todos"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		todos, err := app.repo.GetAll(ctx)
		if err != nil {
			app.serverError(ctx, w, fmt.Errorf("get todos: %v", err))
			return
		}

		var resp response
		resp.Todos = make([]todo, 0, len(todos))
		for _, t := range todos {
			resp.Todos = append(resp.Todos, todo{
				ID:          t.ID,
				Status:      t.Status,
				Description: t.Description,
				CreatedAt:   t.CreatedAt,
				UpdatedAt:   t.UpdatedAt,
			})
		}
		app.jsonResponse(ctx, w, http.StatusOK, &resp)
	})
}

func (app *Application) handleCreateTodo() http.Handler {
	type request struct {
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			app.clientError(
				ctx,

				w,
				http.StatusUnsupportedMediaType,
				codeBadRequest,
				"Content-Type must be application/json",
			)
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			app.logger.WarnContext(ctx, "Decode request body", "err", err)
			app.clientError(ctx, w, http.StatusUnprocessableEntity, codeBadRequest, "Cannot decode request body")
			return
		}

		_, err := app.repo.Create(ctx, Todo{
			Status:      req.Status,
			Description: req.Description,
		})
		if err != nil {
			app.serverError(ctx, w, fmt.Errorf("update todo: %v", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func (app *Application) handleGetTodo() http.Handler {
	type response struct {
		ID          string
		Status      string
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("id")
		if id == "" {
			app.clientError(ctx, w, http.StatusBadRequest, codeBadRequest, "ID must be provided")
			return
		}

		todo, err := app.repo.Get(ctx, id)
		if err != nil {
			if errors.Is(err, ErrTodoNotFound) {
				app.clientError(ctx, w, http.StatusNotFound, codeNotFound, "Todo not found")
				return
			}
			app.serverError(ctx, w, fmt.Errorf("get todo: %v", err))
			return
		}

		app.jsonResponse(ctx, w, http.StatusOK, &response{
			ID:          todo.ID,
			Status:      todo.Status,
			Description: todo.Description,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		})
	})
}

func (app *Application) handleUpdateTodo() http.Handler {
	type request struct {
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("id")
		if id == "" {
			app.clientError(ctx, w, http.StatusBadRequest, codeBadRequest, "ID must be provided")
			return
		}

		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			app.clientError(
				ctx,

				w,
				http.StatusUnsupportedMediaType,
				codeBadRequest,
				"Content-Type must be application/json",
			)
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			app.logger.WarnContext(ctx, "Decode request body", "err", err)
			app.clientError(ctx, w, http.StatusUnprocessableEntity, codeBadRequest, "Cannot decode request body")
			return
		}

		err := app.repo.Update(ctx, Todo{
			ID:          id,
			Status:      req.Status,
			Description: req.Description,
		})
		if err != nil {
			if errors.Is(err, ErrTodoNotFound) {
				app.clientError(ctx, w, http.StatusNotFound, codeNotFound, "Todo not found")
				return
			}
			app.serverError(ctx, w, fmt.Errorf("update todo: %v", err))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})
}

func (app *Application) handleDeleteTodo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("id")
		if id == "" {
			app.clientError(ctx, w, http.StatusBadRequest, codeBadRequest, "ID must be provided")
			return
		}

		err := app.repo.Delete(ctx, id)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})
}

type Todo struct {
	ID          string
	Status      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
