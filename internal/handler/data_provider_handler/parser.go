package data_provider_handler

import (
	"data-provider-service/internal/entity"
	"encoding/json"
	"fmt"
	api "github.com/5krotov/task-resolver-pkg/api/v1"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseSearchTask(r *http.Request) (*entity.SearchTaskParams, error) {
	queryValues := r.URL.Query()

	var err error

	var page int
	pageRequest, ok := queryValues["page"]
	if !ok {
		page = 0
	} else {
		page, err = strconv.Atoi(pageRequest[0])
		if err != nil {
			return nil, fmt.Errorf("page field must be int: %v\n", err)
		}
	}

	var perPage int
	perPageRequest, ok := queryValues["per_page"]
	if !ok {
		perPage = 10
	} else {
		perPage, err = strconv.Atoi(perPageRequest[0])
		if err != nil {
			return nil, fmt.Errorf("per_page field must be int: %v\n", err)
		}
	}

	return &entity.SearchTaskParams{Page: page, PerPage: perPage}, nil
}

func (p *Parser) ParseCreateTask(r *http.Request) (*api.CreateTaskRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read request body failed")
	}
	defer r.Body.Close()

	var request api.CreateTaskRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		return nil, fmt.Errorf("unmarshal request body failed")
	}

	return &request, nil
}

func (p *Parser) ParseGetTask(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	taskIdString, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("no one url pararm 'id'")
	}
	taskId, err := strconv.ParseInt(taskIdString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("convert to int64 failed: %v", err)
	}

	return taskId, nil
}

func (p *Parser) ParseUpdateStatus(r *http.Request) (*api.UpdateStatusRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read request body failed")
	}
	defer r.Body.Close()

	var request api.UpdateStatusRequest

	err = json.Unmarshal(body, &request)
	if err != nil {
		return nil, fmt.Errorf("unmarshal request body failed")
	}

	return &request, nil
}
