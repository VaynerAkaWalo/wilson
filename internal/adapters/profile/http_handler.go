package adapter_profile

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VaynerAkaWalo/go-toolkit/xevent"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"golang-template/internal/domain/action"
	"golang-template/internal/domain/profile"
	"golang-template/internal/domain/transaction"
	"log/slog"
	"net/http"
	"time"
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
		Broker  *xevent.Broker
	}
)

func (handler HttpHandler) RegisterRoutes(router *xhttp.Router) {
	router.RegisterHandler("GET /v1/profiles", handler.getProfiles)
	router.RegisterHandler("POST /v1/profiles", handler.createProfile)
	router.RegisterHandler("GET /v1/profiles/{profileId}/events", handler.profileEvents)
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
	profileId := r.PathValue("profileId")
	if profileId == "" {
		return xhttp.NewError("unknown profileId", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	done := ctx.Done()

	actionChannel := xevent.RegisterListener[action.Event](handler.Broker, ctx)
	defer xevent.RemoveListener[action.Event](handler.Broker, ctx, actionChannel)

	goldChangeChannel := xevent.RegisterListener[transaction.GoldChangeEvent](handler.Broker, ctx)
	defer xevent.RemoveListener[transaction.GoldChangeEvent](handler.Broker, ctx, goldChangeChannel)

	w.WriteHeader(http.StatusOK)

	var err error
	for {
		select {
		case <-done:
			return nil
		case actionEvent := <-actionChannel:
			if profileId != string(actionEvent.ProfileId) {
				break
			}

			event := &ActionEvent{
				Id:   string(actionEvent.Id),
				Gold: actionEvent.GoldReward,
				Exp:  actionEvent.ExpReward,
				Date: time.Now().Unix(),
			}

			err = sendEvent(ctx, w, Action, event)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				return err
			}
		case goldChangeEvent := <-goldChangeChannel:
			if profileId != goldChangeEvent.Profile {
				break
			}

			event := &GoldChangeEvent{
				Id:   goldChangeEvent.Id,
				Gold: goldChangeEvent.GoldBalance,
			}

			err = sendEvent(ctx, w, GoldChange, event)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				return err
			}
		}
	}
}

func sendEvent(ctx context.Context, w http.ResponseWriter, eventType EventType, event interface{}) error {
	rc := http.NewResponseController(w)

	_, err := fmt.Fprintf(w, "event: %s\n", eventType)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}
	jsonData, err := json.Marshal(event)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}
	_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}
	err = rc.Flush()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}
