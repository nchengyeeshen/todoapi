package main

import (
	"context"
	"errors"
	"maps"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
)

type TodoRepository interface {
	Create(ctx context.Context, todo Todo) (string, error)
	Get(ctx context.Context, id string) (Todo, error)
	GetAll(ctx context.Context) ([]Todo, error)
	Update(ctx context.Context, todo Todo) error
	Delete(ctx context.Context, id string) error
}

var _ TodoRepository = (*InMemoryTodoRepository)(nil)

type InMemoryTodoRepository struct {
	mu      sync.RWMutex
	m       map[string]Todo
	counter atomic.Int64
	now     func() time.Time
}

func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		m:   make(map[string]Todo),
		now: time.Now,
	}
}

func (r *InMemoryTodoRepository) Create(ctx context.Context, todo Todo) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	newID := strconv.FormatInt(r.counter.Add(1), 10)
	currentTime := r.now().UTC()

	todo.ID = newID
	todo.CreatedAt = currentTime
	todo.UpdatedAt = currentTime

	r.m[todo.ID] = todo

	return todo.ID, nil
}

func (r *InMemoryTodoRepository) Get(ctx context.Context, id string) (Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, ok := r.m[id]
	if !ok {
		return Todo{}, ErrTodoNotFound
	}

	return todo, nil
}

func (r *InMemoryTodoRepository) GetAll(ctx context.Context) ([]Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return slices.Collect(maps.Values(r.m)), nil
}

func (r *InMemoryTodoRepository) Update(ctx context.Context, todo Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.m[todo.ID]
	if !ok {
		return ErrTodoNotFound
	}

	currentTime := r.now().UTC()

	r.m[todo.ID] = Todo{
		ID:          existing.ID,
		Status:      todo.Status,
		Description: todo.Description,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   currentTime,
	}

	return nil
}

func (r *InMemoryTodoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.m, id)

	return nil
}
