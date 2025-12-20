package adapter_location

import (
	"fmt"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"golang-template/internal/application/location"
	"golang-template/internal/domain/location"
	"log/slog"
	"net/http"
)

type (
	Response struct {
		Id               string  `json:"id"`
		Name             string  `json:"name"`
		Latitude         int     `json:"latitude"`
		Longitude        int     `json:"longitude"`
		RewardMultiplier float64 `json:"rewardMultiplier"`
		Type             string  `json:"type"`
	}

	BulkResponse struct {
		Items []Response `json:"items"`
	}

	HttpHandler struct {
		Service usecase_location.GetLocationService
	}
)

func dto(loc location.Location) Response {
	return Response{
		Id:               string(loc.Id),
		Name:             loc.Name,
		Latitude:         loc.Latitude,
		Longitude:        loc.Longitude,
		RewardMultiplier: loc.RewardMultiplier,
		Type:             string(loc.Type),
	}
}

func (handler HttpHandler) RegisterRoutes(router *xhttp.Router) {
	router.RegisterHandler("GET /v1/locations", handler.getAll)
	router.RegisterHandler("GET /v1/locations/{locationID}", handler.getLocation)
}

func (handler HttpHandler) getLocation(w http.ResponseWriter, r *http.Request) error {
	locationId := r.PathValue("locationID")
	if locationId == "" {
		return xhttp.NewError("unknown location", http.StatusBadRequest)
	}

	loc, err := handler.Service.GetLocation(r.Context(), location.Id(locationId))
	if err != nil {
		slog.ErrorContext(r.Context(), fmt.Sprintf("unable to get location %s, with error %v", locationId, err.Error()))
		return xhttp.NewError("unable to get location", http.StatusInternalServerError)
	}

	return xhttp.WriteResponse(w, http.StatusOK, dto(loc))
}

func (handler HttpHandler) getAll(w http.ResponseWriter, r *http.Request) error {
	locations, err := handler.Service.GetAllLocation(r.Context())
	if err != nil {
		return err
	}

	items := make([]Response, 0)

	for _, loc := range locations {
		items = append(items, dto(loc))
	}

	return xhttp.WriteResponse(w, http.StatusOK, BulkResponse{Items: items})
}
