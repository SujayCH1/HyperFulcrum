package middleware

import (
	"context"
	"net/http"

	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/pkg/utils"
)

type shardContextKey string

const shardPayloadKey shardContextKey = "shardPayload"

func ShardValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.CreateShardDto
		if err := utils.ReadJSONRequest(w, r, &payload); err != nil {
			utils.WriteJSONErrorResponse(w, JSONErrorStatus(err), "Invalid request body", err)
			return
		}
		if err := payload.Validate(); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
			return
		}
		ctx := context.WithValue(r.Context(), shardPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetShardPayload(r *http.Request) (dto.CreateShardDto, bool) {
	payload, ok := r.Context().Value(shardPayloadKey).(dto.CreateShardDto)
	return payload, ok
}
