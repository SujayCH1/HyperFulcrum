package handlers

import (
	"net/http"

	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/pkg/utils"
)

type NodeRuntimeStateHandler struct {
	service *metadata.NodeRuntimeStateService
}

func NewNodeRuntimeStateHandler(service *metadata.NodeRuntimeStateService) *NodeRuntimeStateHandler {
	return &NodeRuntimeStateHandler{service: service}
}

func (h *NodeRuntimeStateHandler) GetByNodeID(w http.ResponseWriter, r *http.Request) {
	state, err := h.service.GetByNodeID(r.Context(), r.PathValue("id"))
	if err != nil {
		writeHandlerError(w, "Node runtime state not found", "Failed to get node runtime state", err)
		return
	}
	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Node runtime state retrieved successfully", state)
}

func (h *NodeRuntimeStateHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	states, err := h.service.ListByProject(r.Context(), r.PathValue("projectId"))
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to list node runtime states", err)
		return
	}
	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Node runtime states retrieved successfully", states)
}
