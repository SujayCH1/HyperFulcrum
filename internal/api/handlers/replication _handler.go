package handlers

import (
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/replication"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

type ReplicationHandler struct {
	service *replication.ReplicationService
}

func NewReplicationHandler(
	service *replication.ReplicationService,
) *ReplicationHandler {
	return &ReplicationHandler{
		service: service,
	}
}

func (h *ReplicationHandler) CreateReplication(
	w http.ResponseWriter,
	r *http.Request,
) {
	payload, ok := middleware.GetReplicationCreatePayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve payload",
			nil,
		)
		return
	}

	topology, err := h.service.CreateReplication(
		r.Context(),
		payload.ProjectID,
		payload.ShardNodeID,
		payload.ReplicaNodeID,
	)
	if err != nil {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to create replication",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusCreated,
		"Replication created successfully",
		topology,
	)
}

func (h *ReplicationHandler) DeleteReplication(
	w http.ResponseWriter,
	r *http.Request,
) {
	payload, ok := middleware.GetReplicationDeletePayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve payload",
			nil,
		)
		return
	}

	err := h.service.DeleteReplication(
		r.Context(),
		payload.RelationID,
		payload.ProjectID,
	)
	if err != nil {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to delete replication",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Replication deleted successfully",
		nil,
	)
}

func (h *ReplicationHandler) PromoteReplica(
	w http.ResponseWriter,
	r *http.Request,
) {
	payload, ok := middleware.GetReplicationPromotePayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve payload",
			nil,
		)
		return
	}

	err := h.service.PromoteReplica(
		r.Context(),
		payload.RelationID,
		payload.ShardNodeID,
		payload.ReplicaNodeID,
	)
	if err != nil {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to promote replica",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Replica promoted successfully",
		nil,
	)
}
