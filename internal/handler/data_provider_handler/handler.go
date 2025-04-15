package data_provider_handler

import (
	"data-provider-service/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type DataProviderHandler struct {
	service *service.DataProviderService
	parser  *Parser
}

func NewDataProviderHandler(service *service.DataProviderService) *DataProviderHandler {
	return &DataProviderHandler{service: service, parser: NewParser()}
}

func (h *DataProviderHandler) Register(r *mux.Router) {
	r.HandleFunc("/api/v1/task", h.HandleCreateTask).Methods("POST")
	r.HandleFunc("/api/v1/task", h.HandleSearchTask).Methods("GET")
	r.HandleFunc("/api/v1/task/{id}", h.HandleGetTask).Methods("GET")
	r.HandleFunc("/api/v1/task/{id}/status", h.HandleUpdateStatus).Methods("PUT")
}

func (h *DataProviderHandler) HandleSearchTask(w http.ResponseWriter, r *http.Request) {

	request, err := h.parser.ParseSearchTask(r)
	if err != nil {
		log.Printf("Parse request failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := h.service.SearchTask(request)
	if err != nil {
		log.Printf("Сreate task failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal response failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)

}

func (h *DataProviderHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {

	request, err := h.parser.ParseCreateTask(r)
	if err != nil {
		log.Printf("Parse request failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateTask(request)
	if err != nil {
		log.Printf("Сreate task failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal response failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(jsonData)

}

func (h *DataProviderHandler) HandleGetTask(w http.ResponseWriter, r *http.Request) {

	taskId, err := h.parser.ParseGetTask(r)
	if err != nil {
		log.Printf("Parse request failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := h.service.GetTask(taskId)
	if err != nil {
		log.Printf("Get task failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal response failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}

func (h *DataProviderHandler) HandleUpdateStatus(w http.ResponseWriter, r *http.Request) {

	request, err := h.parser.ParseUpdateStatus(r)

	err = h.service.UpdateStatus(request)
	if err != nil {
		log.Printf("Update task status failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
