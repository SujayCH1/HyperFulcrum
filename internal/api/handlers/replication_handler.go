package handlers

import (
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
	utils.WriteJSONErrorResponse(
		w,
		http.StatusNotImplemented,
		"Replication API is not implemented",
		nil,
	)
}

func (h *ReplicationHandler) DeleteReplication(
	w http.ResponseWriter,
	r *http.Request,
) {
	utils.WriteJSONErrorResponse(
		w,
		http.StatusNotImplemented,
		"Replication API is not implemented",
		nil,
	)
}

func (h *ReplicationHandler) PromoteReplica(
	w http.ResponseWriter,
	r *http.Request,
) {
	utils.WriteJSONErrorResponse(
		w,
		http.StatusNotImplemented,
		"Replication API is not implemented",
		nil,
	)
}
