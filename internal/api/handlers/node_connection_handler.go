package handlers

import (
	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/internal/repository"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

type NodeConnectionHandler struct {
	service *metadata.NodeConnectionService
}

func NewNodeConnectionHandler(
	service *metadata.NodeConnectionService,
) *NodeConnectionHandler {

	return &NodeConnectionHandler{
		service: service,
	}
}

func (s *NodeConnectionHandler) AddNodeConnection(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("nodeId")

	payload, ok := middleware.GetNodeConnectionPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}

	var nodeConnection repository.NodeConnection
	nodeConnection.NodeId = nodeID
	nodeConnection.Host = payload.Host
	nodeConnection.Port = payload.Port
	nodeConnection.DatabaseName = payload.DatabaseName
	nodeConnection.Username = payload.Username
	nodeConnection.Password = payload.Password

	err := s.service.AddConnection(
		r.Context(),
		&nodeConnection,
	)

	if err != nil {
		writeHandlerError(w, "Node not found", "Failed to add node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusCreated,
		"Node connection successfully added",
		dto.NewNodeConnectionResponse(nodeConnection),
	)
}

func (s *NodeConnectionHandler) RemoveNodeConnection(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("nodeId")

	err := s.service.RemoveConnection(r.Context(), nodeID)

	if err != nil {
		writeHandlerError(w, "Node connection not found", "Failed to remove node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Node connection successfully removed", nil)
}

func (s *NodeConnectionHandler) UpdateNodeConnection(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("nodeId")

	payload, ok := middleware.GetNodeConnectionPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}

	var nodeConnection repository.NodeConnection
	nodeConnection.NodeId = nodeID
	nodeConnection.Host = payload.Host
	nodeConnection.Port = payload.Port
	nodeConnection.DatabaseName = payload.DatabaseName
	nodeConnection.Username = payload.Username
	nodeConnection.Password = payload.Password

	err := s.service.UpdateConnection(r.Context(), &nodeConnection)
	if err != nil {
		writeHandlerError(w, "Node connection not found", "Failed to update node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Node connection successfully updated",
		dto.NewNodeConnectionResponse(nodeConnection),
	)
}

func (s *NodeConnectionHandler) GetNodeConnectionByID(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("nodeId")

	nodeConn, err := s.service.GetConnectionByNodeID(r.Context(), nodeID)
	if err != nil {
		writeHandlerError(w, "Node connection not found", "Failed to get node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Node connection retrieved successfully",
		dto.NewNodeConnectionResponse(nodeConn),
	)
}
