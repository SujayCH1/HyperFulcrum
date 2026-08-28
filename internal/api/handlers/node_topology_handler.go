package handlers

import (
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

type TopologyHandler struct {
	service *metadata.TopologyService
}

func NewTopologyHandler(
	service *metadata.TopologyService,
) *TopologyHandler {

	return &TopologyHandler{
		service: service,
	}
}

// NewTopoogyHandler is retained for compatibility with older callers.
func NewTopoogyHandler(service *metadata.TopologyService) *TopologyHandler {
	return NewTopologyHandler(service)
}

func (h *TopologyHandler) CreateTopology(
	w http.ResponseWriter,
	r *http.Request,
) {
	payload, ok := middleware.GetTopologyCreatePayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve payload",
			nil,
		)
		return
	}

	topology, err := h.service.CreateTopology(
		r.Context(),
		r.PathValue("projectId"),
		payload.ShardID,
		payload.StandbyNodeID,
	)

	if err != nil {
		writeHandlerError(
			w,
			"Topology not found",
			"Failed to create topology",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusCreated,
		"Topology created successfully",
		topology,
	)

}

func (h *TopologyHandler) DeleteTopology(
	w http.ResponseWriter,
	r *http.Request,
) {
	payload, ok := middleware.GetTopologyDeletePayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve payload",
			nil,
		)
		return
	}

	err := h.service.DeleteTopology(
		r.Context(),
		payload.RelationID,
		payload.ProjectID,
	)
	if err != nil {
		writeHandlerError(
			w,
			"Topology not found",
			"Failed to delete topology",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TopologyHandler) DeleteTopologyByPath(w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteTopology(r.Context(), r.PathValue("relationId"), r.PathValue("projectId"))
	if err != nil {
		writeHandlerError(w, "Topology not found", "Failed to delete topology", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TopologyHandler) ListTopologies(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectID := r.PathValue("projectId")

	topologies, err := h.service.ListTopologies(
		r.Context(),
		projectID,
	)
	if err != nil {
		writeHandlerError(
			w,
			"Project not found",
			"Failed to fetch topology",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Topology fetched successfully",
		topologies,
	)
}

func (h *TopologyHandler) GetTopologyByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectID := r.PathValue("projectId")
	relationID := r.PathValue("relationId")

	topology, err := h.service.GetTopologyByID(
		r.Context(),
		projectID,
		relationID,
	)
	if err != nil {
		writeHandlerError(
			w,
			"Topology not found",
			"Failed to fetch topology",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Topology fetched successfully",
		topology,
	)
}
