package middleware

import (
	"context"
	"hyperfulcrum/internal/api/dto"
	"hyperfulcrum/pkg/utils"
	"net/http"
)

type contextKey string

const (
	projectPayloadKey contextKey = "projectPayload"
)

func ProjectValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.ProjectDto

		if err := utils.ReadJSONRequest(r, &payload); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if err := payload.Validate(); err != nil {
			utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
			return
		}

		ctx := context.WithValue(r.Context(), projectPayloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetProjectPayload(r *http.Request) (dto.ProjectDto, bool) {
	payload, ok := r.Context().Value(projectPayloadKey).(dto.ProjectDto)
	return payload, ok
}
