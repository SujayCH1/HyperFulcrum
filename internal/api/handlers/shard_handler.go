package handlers

import (
	"net/http"

	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/pkg/utils"
)

type ShardHandler struct{ service *metadata.ShardService }

func NewShardHandler(service *metadata.ShardService) *ShardHandler {
	return &ShardHandler{service: service}
}

func (h *ShardHandler) AddShard(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetShardPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}
	shard, err := h.service.AddShard(r.Context(), r.PathValue("projectId"), payload.Name, payload.PrimaryNodeID)
	if err != nil {
		writeHandlerError(w, "Project or primary node not found", "Failed to add shard", err)
		return
	}
	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Shard added successfully", shard)
}

func (h *ShardHandler) ListShards(w http.ResponseWriter, r *http.Request) {
	shards, err := h.service.ListShards(r.Context(), r.PathValue("projectId"))
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to list shards", err)
		return
	}
	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Shards retrieved successfully", shards)
}

func (h *ShardHandler) GetShard(w http.ResponseWriter, r *http.Request) {
	shard, err := h.service.GetShard(r.Context(), r.PathValue("id"))
	if err != nil {
		writeHandlerError(w, "Shard not found", "Failed to get shard", err)
		return
	}
	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Shard retrieved successfully", shard)
}

func (h *ShardHandler) RenameShard(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if err := dto.ValidateShardName(name); err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid shard name", err)
		return
	}
	if err := h.service.RenameShard(r.Context(), r.PathValue("id"), name); err != nil {
		writeHandlerError(w, "Shard not found", "Failed to rename shard", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ShardHandler) RemoveShard(w http.ResponseWriter, r *http.Request) {
	if err := h.service.RemoveShard(r.Context(), r.PathValue("id")); err != nil {
		writeHandlerError(w, "Shard not found", "Failed to remove shard", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
