package middleware

import (
	"context"
	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

type contentKey string

const (
	topologyPayloadKey contentKey = "topologyPayload"
)

func TopologyCreateValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.TopologyCreateDto

		if err := utils.ReadJSONRequest(r, &payload); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if err := payload.ValidateCreate(); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
			return
		}

		ctx := context.WithValue(r.Context(), topologyPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTopologyCreatePayload(r *http.Request) (dto.TopologyCreateDto, bool) {
	payload, ok := r.Context().Value(topologyPayloadKey).(dto.TopologyCreateDto)
	return payload, ok
}

func TopologyDeleteValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.TopologyDeleteDto

		if err := utils.ReadJSONRequest(r, &payload); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if err := payload.ValidateDelete(); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
			return
		}

		ctx := context.WithValue(r.Context(), topologyPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTopologyDeletePayload(r *http.Request) (dto.TopologyDeleteDto, bool) {
	payload, ok := r.Context().Value(topologyPayloadKey).(dto.TopologyDeleteDto)
	return payload, ok
}
