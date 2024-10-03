package helpers

import (
	"errors"
	"strconv"
)

func ParseAndValidateID(idStr string) (uint64, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id must be non-zero value")
	}

	return id, nil
}
