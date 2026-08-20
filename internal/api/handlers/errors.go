package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/lib/pq"

	"hyperfulcrum/pkg/logger"
	"hyperfulcrum/pkg/utils"
)

func writeHandlerError(
	w http.ResponseWriter,
	notFoundMessage string,
	message string,
	err error,
) {
	if errors.Is(err, context.DeadlineExceeded) {
		utils.WriteJSONErrorResponse(w, http.StatusGatewayTimeout, "Request timed out", nil)
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		utils.WriteJSONErrorResponse(w, http.StatusNotFound, notFoundMessage, nil)
		return
	}

	var databaseError *pq.Error
	if errors.As(err, &databaseError) {
		switch databaseError.Code {
		case "23505", "23503":
			utils.WriteJSONErrorResponse(w, http.StatusConflict, message, nil)
			return
		case "23514":
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, message, nil)
			return
		}
	}

	logger.Logger.Error(
		message,
		"request_id", w.Header().Get("X-Request-ID"),
		"error", err,
	)
	utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, message, err)
}
