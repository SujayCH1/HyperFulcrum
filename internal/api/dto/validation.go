package dto

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

var validUUID = validation.By(func(value any) error {
	id, ok := value.(string)
	if !ok {
		return errors.New("must be a valid UUID")
	}

	if _, err := uuid.Parse(id); err != nil {
		return errors.New("must be a valid UUID")
	}

	return nil
})
