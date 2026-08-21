package middleware

import (
	"context"
	"net/http"

	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/pkg/utils"
)

type shardKeyContextKey string

const shardKeyPayloadKey shardKeyContextKey = "shardKeyPayload"

func ShardKeyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.ShardKeyCreateDto

		err := utils.ReadJSONRequest(w, r, &payload)
		if err != nil {
			utils.WriteJSONErrorResponse(
				w,
				JSONErrorStatus(err),
				"Invalid request body",
				err,
			)
			return
		}

		err = payload.Validate()
		if err != nil {
			utils.WriteJSONErrorResponse(
				w,
				http.StatusBadRequest,
				"Validation failed",
				err,
			)
			return
		}

		ctx := context.WithValue(r.Context(), shardKeyPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetShardKeyPayload(r *http.Request) (dto.ShardKeyCreateDto, bool) {
	payload, ok := r.Context().Value(shardKeyPayloadKey).(dto.ShardKeyCreateDto)
	return payload, ok
}
