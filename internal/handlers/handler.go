package handlers

import (
	client "assignment-2/internal/client/restcountries"
	utils "assignment-2/internal/models"
	"assignment-2/internal/store"
	"context"
	"encoding/json"
	"net/http"

	models "assignment-2/internal/models"
)

// StoreInterface allows both *store.FireStore and *store.MockStore to be used interchangeably.
type StoreInterface interface {
	CreateRegistration(ctx context.Context, apiKey string, reg models.Registration) (string, error)
	GetRegistration(ctx context.Context, apiKey string, id string) (*models.Registration, error)
	GetAllRegistrations(ctx context.Context, apiKey string) ([]models.Registration, error)
	UpdateRegistration(ctx context.Context, apiKey string, id string, reg models.Registration) error
	DeleteRegistration(ctx context.Context, apiKey string, id string) error
	TweakRegistration(ctx context.Context, apiKey string, id string, patch models.RegistrationPatch) error
	APIKeyExists(ctx context.Context, keyHash string) bool
	ApiKeyExists(ctx context.Context, keyHash string) bool
	CountApiPerUser(ctx context.Context, email string) (int, error)
	CreateApiStorage(ctx context.Context, reg models.Authentication) error
	DeleteAPIkey(ctx context.Context, apiKey string) error
}

// CacheInterface allows both *store.Cache and mockCache to be used interchangeably.
type CacheInterface interface {
	RequestFromCache(req store.CacheExternalRequest) (store.CacheResponse, error)
}

type Handler struct {
	store               StoreInterface
	restCountriesClient client.RestCountriesClient
	cache               CacheInterface
}

func NewHandler(s StoreInterface, restCountriesClient client.RestCountriesClient) *Handler {
	return &Handler{
		store:               s,
		restCountriesClient: restCountriesClient,
	}
}

func NewFirestoreHandler(s *store.FireStore, cache CacheInterface) *Handler {
	return &Handler{
		store: s,
		cache: cache,
	}
}

func writeJSONError(w http.ResponseWriter, code int, errMsg string) {
	response := utils.ErrorResponse{
		Code:    code,
		Message: errMsg,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonBytes)
}
