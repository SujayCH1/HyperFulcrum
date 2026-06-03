package services

import (
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/repository"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

type NodeConnectionService struct {
	repo *repository.NodeConnectionRepository
}

func NewNodeConnectionService(repo *repository.NodeConnectionRepository) *NodeConnectionService {
	return &NodeConnectionService{repo: repo}
}

func (s *NodeConnectionService) AddNodeConnection(w http.ResponseWriter, r *http.Request) {
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

	err := s.repo.ConnectionAdd(
		r.Context(),
		&nodeConnection,
	)

	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to add node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node connection successfully added", payload)
}

func (s *NodeConnectionService) RemoveNodeConnection(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetNodeConnectionPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}

	err := s.repo.ConnectionRemove(r.Context(), payload.NodeID)

	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to remove node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node connection successfully added", payload)
}

func (s *NodeConnectionService) UpdateNodeConnection(w http.ResponseWriter, r *http.Request) {
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

	err := s.repo.ConnectionUpdate(r.Context(), &nodeConnection)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to update node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node connection successfully added", payload)
}

func (s *NodeConnectionService) GetNodeConnectionByID(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetNodeConnectionPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload", nil)
		return
	}

	nodeConn, err := s.repo.GetConnectionByNodeId(r.Context(), payload.NodeID)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to get node connection", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Node connection successfully added", &nodeConn)
}
