package data_provider_handler

import (
	"data-provider-service/internal/service"
	"encoding/json"
	api "github.com/5krotov/task-resolver-pkg/api/v1"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

type DataProviderHandler struct {
	service *service.DataProviderService
}

func NewDataProviderHandler(service *service.DataProviderService) *DataProviderHandler {
	return &DataProviderHandler{service: service}
}

func (h *DataProviderHandler) Register(r *mux.Router) {
	r.HandleFunc("/api/v1/task", h.HandleCreateTask).Methods("POST")
	r.HandleFunc("/api/v1/task/{id}", h.HandleGetTask).Methods("GET")
	r.HandleFunc("/api/v1/task/{id}/status", h.HandleUpdateStatus).Methods("PUT")
}

func (h *DataProviderHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Read request body failed\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request api.CreateTaskRequest

	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Unmarshal request body failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateTask(&request)
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

func (h *DataProviderHandler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIdString, ok := vars["id"]
	if !ok {
		log.Printf("Get url pararm 'id' failed\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskId, err := strconv.ParseInt(taskIdString, 10, 64)
	if err != nil {
		log.Printf("Convert url pararm 'id' to int64 failed: %v\n", err)
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Read request body failed\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request api.UpdateStatusRequest

	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Unmarshal request body failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.UpdateStatus(&request)
	if err != nil {
		log.Printf("Update task status failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
