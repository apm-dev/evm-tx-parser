package common

import (
	"fmt"
	"net/http"

	"github.com/apm-dev/evm-tx-parser/src/domain"
	"github.com/pkg/errors"
)

func ErrToHttpCodeAndMessage(err error, methodName string) (int, string) {
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, fmt.Sprintf("failed to %s", methodName)
	}
}
