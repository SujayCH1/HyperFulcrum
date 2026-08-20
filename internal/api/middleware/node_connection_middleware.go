package middleware

import (
	"context"
	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

type contextey string

const (
	nodeConnectionPayloadkey contextey = "nodeConnectionPayload"
)

func NodeConnectionValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.NodeConnectionDto

		if err := utils.ReadJSONRequest(w, r, &payload); err != nil {
			utils.WriteJSONErrorResponse(w, JSONErrorStatus(err), "Invalid request body", err)
			return
		}

		if err := payload.Validate(); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
			return
		}

		ctx := context.WithValue(r.Context(), nodeConnectionPayloadkey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetNodeConnectionPayload(r *http.Request) (dto.NodeConnectionDto, bool) {
	payload, ok := r.Context().Value(nodeConnectionPayloadkey).(dto.NodeConnectionDto)
	return payload, ok
}
