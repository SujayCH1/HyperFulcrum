package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/lib/pq"

	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/internal/repository"
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

	if errors.Is(err, metadata.ErrInvalidNodeType) ||
		errors.Is(err, metadata.ErrInvalidConnection) ||
		errors.Is(err, metadata.ErrTopologySelfRelation) ||
		errors.Is(err, metadata.ErrTopologyRoleMismatch) {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, message, nil)
		return
	}

	if errors.Is(err, metadata.ErrProjectRunning) ||
		errors.Is(err, metadata.ErrProjectHasNodes) ||
		errors.Is(err, metadata.ErrProjectHasTopology) ||
		errors.Is(err, metadata.ErrDuplicateNodeName) ||
		errors.Is(err, metadata.ErrNodeActive) ||
		errors.Is(err, metadata.ErrNodeInTopology) ||
		errors.Is(err, metadata.ErrConnectionExists) ||
		errors.Is(err, metadata.ErrConnectionNotFound) ||
		errors.Is(err, metadata.ErrReplicaAlreadyUsed) ||
		errors.Is(err, metadata.ErrShardIsReplica) ||
		errors.Is(err, metadata.ErrDuplicateTopology) ||
		errors.Is(err, metadata.ErrSchemaNotLocked) ||
		errors.Is(err, repository.ErrSchemaNotLocked) ||
		errors.Is(err, repository.ErrSchemaLocked) ||
		errors.Is(err, repository.ErrSchemaActivated) ||
		errors.Is(err, repository.ErrSchemaHasKeys) ||
		errors.Is(err, repository.ErrSchemaRevision) {
		utils.WriteJSONErrorResponse(w, http.StatusConflict, message, nil)
		return
	}

	if errors.Is(err, repository.ErrSchemaEmpty) {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, message, nil)
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
