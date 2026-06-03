package middleware

import (
	"context"
	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

type nodeContextKey string
const nodePayloadKey nodeContextKey = "nodePayload"

func NodeValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var payload dto.CreateNodeDto
		if err := utils.ReadJSONRequest(r, &payload); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if err := payload.Validate(); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
			return
		}

		ctx := context.WithValue(r.Context(), nodePayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetNodePayload(r *http.Request) (dto.CreateNodeDto, bool) {
	payload, ok := r.Context().Value(nodePayloadKey).(dto.CreateNodeDto)
	return payload, ok
}
