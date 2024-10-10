package main

import "fmt"

type NotFoundError struct {
	Resource string
	ID       string
}

type DatabaseError struct {
	Operation string
	Err       error
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource %s with ID %s not found", e.Resource, e.ID)
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}
