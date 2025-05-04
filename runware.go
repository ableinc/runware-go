package runware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	Image    = "imageInference"
	v1Domain = "https://api.runware.ai/v1"
)

type TaskType string

type RunwareRequest struct {
	ApiKey          string
	TaskType        TaskType
	TaskUUID        string
	Prompt          string
	Width           int8
	Height          int8
	Model           string
	NumberOfResults int8
}

type RunwareResponseBody struct {
	TaskType  string `json:"taskType"`
	TaskUUID  string `json:"taskUUID"`
	ImageUUID string `json:"imageUUID"`
	ImageUrl  string `json:"imageUrl"`
}

// Interface definition
type GenerateImagesV1 interface {
	Config(taskType string, taskUUID string, prompt string, width int8, height int8, model string, results int8) GenerateImagesV1
	Generate() (*RunwareResponseBody, error)
}

// Struct implementing the interface
type generateImagesV1Impl struct {
	apiKey  string
	request RunwareRequest
}

func NewGenerateImagesV1(apiKey string) GenerateImagesV1 {
	return &generateImagesV1Impl{
		apiKey: apiKey,
		request: RunwareRequest{
			ApiKey: apiKey,
		},
	}
}

func (g *generateImagesV1Impl) Config(taskType string, taskUUID string, prompt string, width int8, height int8, model string, results int8) GenerateImagesV1 {
	g.request.TaskType = TaskType(taskType)
	g.request.TaskUUID = taskUUID
	g.request.Prompt = prompt
	g.request.Width = width
	g.request.Height = height
	g.request.Model = model
	g.request.NumberOfResults = results
	return g
}

func (g *generateImagesV1Impl) Generate() (*RunwareResponseBody, error) {
	return sendRequest(g.request, v1Domain)
}

// Helper functions

func buildClient(request RunwareRequest, url string) (*http.Client, *http.Request, error) {
	payload := map[string]any{
		"taskType":        request.TaskType,
		"taskUUID":        request.TaskUUID,
		"positivePrompt":  request.Prompt,
		"width":           request.Width,
		"height":          request.Height,
		"model":           request.Model,
		"numberOfResults": request.NumberOfResults,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.ApiKey))
	return client, req, nil
}

func sendRequest(request RunwareRequest, url string) (*RunwareResponseBody, error) {
	client, req, err := buildClient(request, url)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	var response RunwareResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
