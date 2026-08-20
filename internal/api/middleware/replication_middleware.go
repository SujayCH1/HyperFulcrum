package middleware

import (
	"context"
	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

const (
	replicationPayloadKey contentKey = "replicationPayload"
)

func ReplicationCreateValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.CreateReplicationDto

		if err := utils.ReadJSONRequest(w, r, &payload); err != nil {
			utils.WriteJSONErrorResponse(
				w,
				JSONErrorStatus(err),
				"Invalid request body",
				err,
			)
			return
		}

		if err := payload.ValidateCreate(); err != nil {
			utils.WriteJSONErrorResponse(
				w,
				http.StatusBadRequest,
				"Validation failed",
				err,
			)
			return
		}

		ctx := context.WithValue(r.Context(), replicationPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetReplicationCreatePayload(r *http.Request) (dto.CreateReplicationDto, bool) {
	payload, ok := r.Context().Value(replicationPayloadKey).(dto.CreateReplicationDto)
	return payload, ok
}

func ReplicationDeleteValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.DeleteReplicationDto

		if err := utils.ReadJSONRequest(w, r, &payload); err != nil {
			utils.WriteJSONErrorResponse(
				w,
				JSONErrorStatus(err),
				"Invalid request body",
				err,
			)
			return
		}

		if err := payload.ValidateDelete(); err != nil {
			utils.WriteJSONErrorResponse(
				w,
				http.StatusBadRequest,
				"Validation failed",
				err,
			)
			return
		}

		ctx := context.WithValue(r.Context(), replicationPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetReplicationDeletePayload(r *http.Request) (dto.DeleteReplicationDto, bool) {
	payload, ok := r.Context().Value(replicationPayloadKey).(dto.DeleteReplicationDto)
	return payload, ok
}

func ReplicationPromoteValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.PromoteReplicaDto

		if err := utils.ReadJSONRequest(w, r, &payload); err != nil {
			utils.WriteJSONErrorResponse(
				w,
				JSONErrorStatus(err),
				"Invalid request body",
				err,
			)
			return
		}

		if err := payload.ValidatePromote(); err != nil {
			utils.WriteJSONErrorResponse(
				w,
				http.StatusBadRequest,
				"Validation failed",
				err,
			)
			return
		}

		ctx := context.WithValue(r.Context(), replicationPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetReplicationPromotePayload(r *http.Request) (dto.PromoteReplicaDto, bool) {
	payload, ok := r.Context().Value(replicationPayloadKey).(dto.PromoteReplicaDto)
	return payload, ok
}
