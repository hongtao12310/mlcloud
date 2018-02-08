package dao

import "fmt"

const (
    USER_NOT_FOUND int = 10000
)

type JobNotFoundError struct {
    name, username string
}

type JobAlreadyExistError struct {
    name, username string
}

func (e JobNotFoundError) Error() string {
    return fmt.Sprintf("job %s of user %s not found", e.name, e.username)
}

func (e JobAlreadyExistError) Error() string {
    return fmt.Sprintf("job %s of user %s already exists", e.name, e.username)
}