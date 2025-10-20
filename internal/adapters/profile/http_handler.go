package adapter_profile

import (
	"encoding/json"
	"fmt"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"golang-template/internal/domain/action"
	"golang-template/internal/domain/profile"
	"golang-template/pkg/ievent"
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
		Service           profile.Service
		EventOrchestrator *ievent.Orchestrator
	}
)

func (handler HttpHandler) RegisterRoutes(router *xhttp.Router) {
	router.RegisterHandler("GET /v1/profiles", handler.getProfiles)
	router.RegisterHandler("POST /v1/profiles", handler.createProfile)
	router.RegisterHandler("GET /v1/profiles/{owner}/events", handler.profileEvents)
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

func (handler HttpHandler) profileEvents(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	owner := r.PathValue("owner")
	if owner == "" {
		return xhttp.NewError("unknown owner", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	done := ctx.Done()

	dataChannel, err := handler.EventOrchestrator.RegisterListener(ctx)
	if err != nil {
		return err
	}
	defer handler.EventOrchestrator.UnregisterListener(ctx, dataChannel)
	rc := http.NewResponseController(w)
	w.WriteHeader(http.StatusOK)

	for {
		select {
		case <-done:
			return nil
		case event := <-dataChannel:
			ev, ok := event.(action.Event)
			if !ok || string(ev.Owner) != owner {
				break
			}
			_, err := fmt.Fprintf(w, "event: %s\n", "action-reward")
			if err != nil {
				slog.ErrorContext(r.Context(), err.Error())
				return err
			}
			jsonData, err := json.Marshal(ev)
			if err != nil {
				slog.ErrorContext(r.Context(), err.Error())
				return err
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
			if err != nil {
				slog.ErrorContext(r.Context(), err.Error())
				return err
			}
			err = rc.Flush()
			if err != nil {
				slog.ErrorContext(r.Context(), err.Error())
				return err
			}
		}
	}
}
