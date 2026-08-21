package handlers

import (
	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/pkg/utils"
	"net/http"
	"strconv"
)

type NodeHandler struct {
	service *metadata.NodeService
}

func NewNodeHandler(
	service *metadata.NodeService,
) *NodeHandler {

	return &NodeHandler{
		service: service,
	}
}

func (h *NodeHandler) AddNode(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	payload, ok := middleware.GetNodePayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload from context", nil)
		return
	}

	node, err := h.service.AddNode(r.Context(), projectID, payload.Type, payload.Name)
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to add node", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node added successfully", node)
}

func (h *NodeHandler) ListNodes(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	nodes, err := h.service.ListNodes(r.Context(), projectID)
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to list nodes", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Nodes retrieved successfully", nodes)
}

func (h *NodeHandler) RemoveNode(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")

	err := h.service.RemoveNode(r.Context(), nodeID)
	if err != nil {
		writeHandlerError(w, "Node not found", "Failed to remove node", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NodeHandler) UpdateNodeName(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	name := r.URL.Query().Get("name")

	if err := dto.ValidateNodeName(name); err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid node name", err)
		return
	}

	err := h.service.UpdateNodeName(r.Context(), nodeID, name)
	if err != nil {
		writeHandlerError(w, "Node not found", "Failed to update node name", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NodeHandler) UpdateNodeStatus(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	statusStr := r.URL.Query().Get("status")

	if statusStr == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Status is required", nil)
		return
	}

	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid status value", err)
		return
	}

	err = h.service.UpdateNodeStatus(r.Context(), nodeID, status)
	if err != nil {
		writeHandlerError(w, "Node not found", "Failed to update node status", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NodeHandler) UpdateNodeType(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	nodeType := r.URL.Query().Get("type")

	if err := dto.ValidateNodeType(nodeType); err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid node type", err)
		return
	}

	err := h.service.UpdateNodeType(r.Context(), nodeID, nodeType)
	if err != nil {
		writeHandlerError(w, "Node not found", "Failed to update node type", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
