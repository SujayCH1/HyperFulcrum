package handlers

import (
	"net/http"

	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/pkg/utils"
)

type ShardKeyHandler struct {
	service *metadata.ShardKeysService
}

func NewShardKeyHandler(service *metadata.ShardKeysService) *ShardKeyHandler {
	return &ShardKeyHandler{service: service}
}

func (h *ShardKeyHandler) AddShardKey(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	payload, ok := middleware.GetShardKeyPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve payload",
			nil,
		)
		return
	}

	key, err := h.service.AddShardKey(
		r.Context(),
		projectID,
		payload.TableName,
		payload.KeyColumn,
	)
	if err != nil {
		writeHandlerError(
			w,
			"Project, table, or column not found",
			"Failed to add shard key",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusCreated,
		"Shard key added successfully",
		dto.NewShardKeyResponse(key),
	)
}

func (h *ShardKeyHandler) DeleteShardKey(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	keyID := r.PathValue("id")

	err := h.service.DeleteShardKey(r.Context(), projectID, keyID)
	if err != nil {
		writeHandlerError(w, "Shard key not found", "Failed to delete shard key", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShardKeyHandler) GetShardKey(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	tableName := r.PathValue("tableName")

	key, err := h.service.GetShardKey(r.Context(), projectID, tableName)
	if err != nil {
		writeHandlerError(w, "Shard key not found", "Failed to get shard key", err)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Shard key retrieved successfully",
		dto.NewShardKeyResponse(key),
	)
}

func (h *ShardKeyHandler) ListShardKeys(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	keys, err := h.service.ListShardKeys(r.Context(), projectID)
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to list shard keys", err)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Shard keys retrieved successfully",
		dto.NewShardKeyListResponse(keys),
	)
}
