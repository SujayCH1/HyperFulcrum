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

func NewTopoogyHandler(
	service *metadata.TopologyService,
) *TopologyHandler {

	return &TopologyHandler{
		service: service,
	}
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
		payload.ProjectID,
		payload.ReplicaNodeID,
		payload.ShardNodeID,
	)

	if err != nil {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to create topology",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusCreated,
		"TOpology created successfully",
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
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to delete toology",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusOK,
		"Project created successfully",
		nil,
	)
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
		utils.WriteJSONErrorResponse(
			w,
			http.StatusNotFound,
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
		utils.WriteJSONErrorResponse(
			w,
			http.StatusNotFound,
			"Topology not found",
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
