package handlers

import (
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
	payload, ok := middleware.GetNodeConnectionPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}

	var nodeConnection repository.NodeConnection
	nodeConnection.NodeId = payload.NodeID
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
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to add node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node connection successfully added", payload)
}

func (s *NodeConnectionHandler) RemoveNodeConnection(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetNodeConnectionPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}

	err := s.service.RemoveConnection(r.Context(), payload.NodeID)

	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to remove node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node connection successfully added", payload)
}

func (s *NodeConnectionHandler) UpdateNodeConnection(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetNodeConnectionPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}

	var nodeConnection repository.NodeConnection
	nodeConnection.NodeId = payload.NodeID
	nodeConnection.Host = payload.Host
	nodeConnection.Port = payload.Port
	nodeConnection.DatabaseName = payload.DatabaseName
	nodeConnection.Username = payload.Username
	nodeConnection.Password = payload.Password

	err := s.service.UpdateConnection(r.Context(), &nodeConnection)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to update node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node connection successfully added", payload)
}

func (s *NodeConnectionHandler) GetNodeConnectionByID(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetNodeConnectionPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}

	nodeConn, err := s.service.GetConnectionByNodeID(r.Context(), payload.NodeID)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to get node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node connection successfully added", &nodeConn)
}
