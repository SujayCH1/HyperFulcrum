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

		var payload dto.NodeDto
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

func GetNodePayload(r *http.Request) (dto.NodeDto, bool) {
	payload, ok := r.Context().Value(nodePayloadKey).(dto.NodeDto)
	return payload, ok
}

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}