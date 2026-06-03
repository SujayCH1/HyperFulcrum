package services

import (
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/repository"
	"hyperfulcrum/pkg/utils"
	"net/http"
	"strconv"
)

type NodeService struct {
	repo *repository.NodeRepository
}

func NewNodeService(repo *repository.NodeRepository) *NodeService {
	return &NodeService{repo: repo}
}

func (h *NodeService) AddNode(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "projectID is required", nil)
		return
	}

	payload, ok := middleware.GetNodePayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload from context", nil)
		return
	}

	err := h.repo.NodeAdd(r.Context(), projectID, payload.Type, payload.Name)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to add node", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node added successfully", nil)
}

func (h *NodeService) ListNodes(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "projectID is required", nil)
		return
	}

	nodes, err := h.repo.NodeList(r.Context(), projectID)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to list nodes", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Nodes retrieved successfully", nodes)
}

func (h *NodeService) RemoveNode(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	if nodeID == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "nodeID is required", nil)
		return
	}

	err := h.repo.NodeRemove(r.Context(), nodeID)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to remove node", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Node removed successfully", nil)
}

func (h *NodeService) UpdateNodeName(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	name := r.URL.Query().Get("name")

	if nodeID == "" || name == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "nodeID and name are required", nil)
		return
	}

	err := h.repo.NodeUpdateName(r.Context(), nodeID, name)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to update node name", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Node name updated successfully", nil)
}

func (h *NodeService) UpdateNodeStatus(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	statusStr := r.URL.Query().Get("status")

	if nodeID == "" || statusStr == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "nodeID and status are required", nil)
		return
	}

	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid status value", err)
		return
	}

	err = h.repo.NodeUpdateStatus(r.Context(), nodeID, status)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to update node status", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Node status updated successfully", nil)
}

func (h *NodeService) UpdateNodeType(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	nodeType := r.URL.Query().Get("type")

	if nodeID == "" || nodeType == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "nodeID and type are required", nil)
		return
	}

	err := h.repo.NodeUpdateType(r.Context(), nodeID, nodeType)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to update node type", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Node type updated successfully", nil)
}
