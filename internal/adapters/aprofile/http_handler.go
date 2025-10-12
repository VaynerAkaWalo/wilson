package aprofile

import (
	"encoding/json"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"golang-template/internal/domain/profile"
	"log/slog"
	"net/http"
)

type (
	Response struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Level int64  `json:"level"`
		Gold  int64  `json:"gold"`
	}

	Request struct {
		Name string `json:"name"`
	}

	HttpHandler struct {
		Service profile.Service
	}
)

func (handler HttpHandler) RegisterRoutes(router *xhttp.Router) {
	router.RegisterHandler("GET /v1/profiles", handler.getProfiles)
	router.RegisterHandler("POST /v1/profiles", handler.createProfile)
}

func (handler HttpHandler) getProfiles(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	owner, ok := ctx.Value(xhttp.UserId).(string)
	if !ok {
		slog.ErrorContext(ctx, "error while parsing identity ID from context")
		return xhttp.NewError("internal server error", http.StatusInternalServerError)
	}

	profiles, err := handler.Service.GetProfilesByOwner(ctx, profile.OwnerId(owner))
	if err != nil {
		return err
	}

	responseDtos := make([]Response, len(profiles))
	for i, prof := range profiles {
		responseDtos[i] = Response{
			Id:    string(prof.Id),
			Name:  prof.Name,
			Level: prof.Level,
			Gold:  prof.Gold,
		}
	}

	return xhttp.WriteResponse(w, http.StatusOK, responseDtos)
}

func (handler HttpHandler) createProfile(w http.ResponseWriter, r *http.Request) error {
	var request Request

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return xhttp.NewError("request does not satisfy specified schema", http.StatusBadRequest)
	}

	prof, err := handler.Service.CreateProfile(r.Context(), request.Name)
	if err != nil {
		return err
	}

	response := Response{
		Id:    string(prof.Id),
		Name:  prof.Name,
		Level: prof.Level,
		Gold:  prof.Gold,
	}

	return xhttp.WriteResponse(w, http.StatusCreated, response)
}
