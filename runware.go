package runware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TaskType string
type OutputType string
type OutputFormat string

const (
	ImageInference TaskType     = "imageInference"
	Base64Data     OutputType   = "base64Data"
	DataURI        OutputType   = "dataURI"
	URL            OutputType   = "URL"
	PNG            OutputFormat = "PNG"
	JPG            OutputFormat = "JPG"
	WEBP           OutputFormat = "WEBP"
)

type RunwareOptions struct {
	ApiKey          string
	TaskType        TaskType
	TaskUUID        string
	Prompt          string
	Width           int8
	Height          int8
	Model           string
	NumberOfResults int8
	UploadEndpoint  string
	CheckNSFW       bool
	IncludeCost     bool
	OutputType      OutputType
	OutputFormat    OutputFormat
}

type RunwareResponseBody struct {
	TaskType        string  `json:"taskType,omitempty"`
	TaskUUID        string  `json:"taskUUID,omitempty"`
	ImageUUID       string  `json:"imageUUID,omitempty"`
	ImageUrl        string  `json:"imageUrl,omitempty"`
	ImageBase64Data string  `json:"imageBase64Data,omitempty"`
	ImageDataURI    string  `json:"imageDataURI,omitempty"`
	Seed            int8    `json:"seed,omitempty"`
	NSFWContent     bool    `json:"NSFWContent,omitempty"`
	Cost            float64 `json:"cost,omitempty"`
}

// Interface definition
type GenerateImagesV1 interface {
	Config(data map[string]any) GenerateImagesV1
	GenerateV1() (*RunwareResponseBody, error)
}

// Struct implementing the interface
type generateImagesV1Impl struct {
	apiKey  string
	options RunwareOptions
}

func NewGenerateImagesV1(apiKey string) GenerateImagesV1 {
	return &generateImagesV1Impl{
		apiKey: apiKey,
		options: RunwareOptions{
			ApiKey: apiKey,
		},
	}
}

func (g *generateImagesV1Impl) Config(data map[string]any) GenerateImagesV1 {
	if data["taskType"] != nil {
		g.options.TaskType = data["taskType"].(TaskType)
	}
	if data["taskUUID"] != nil {
		g.options.TaskUUID = data["taskUUID"].(string)
	}
	if data["prompt"] != nil {
		g.options.Prompt = data["prompt"].(string)
	}
	if data["width"] != nil {
		g.options.Width = data["width"].(int8)
	}
	if data["height"] != nil {
		g.options.Height = data["height"].(int8)
	}
	if data["model"] != nil {
		g.options.Model = data["model"].(string)
	}
	if data["results"] != nil {
		g.options.NumberOfResults = data["results"].(int8)
	}
	if data["uploadEndpoint"] != nil {
		g.options.UploadEndpoint = data["uploadEndpoint"].(string)
	}
	if data["checkNSFW"] != nil {
		g.options.CheckNSFW = data["checkNSFW"].(bool)
	}
	if data["includeCost"] != nil {
		g.options.IncludeCost = data["includeCost"].(bool)
	}
	if data["outputType"] != nil {
		g.options.OutputType = data["outputType"].(OutputType)
	}
	if data["outputFormat"] != nil {
		g.options.OutputFormat = data["outputFormat"].(OutputFormat)
	}
	return g
}

func (g *generateImagesV1Impl) GenerateV1() (*RunwareResponseBody, error) {
	var v1Domain string = "https://api.runware.ai/v1"
	return sendRequest(g.options, v1Domain)
}

// Helper functions

func buildClient(request RunwareOptions, url string) (*http.Client, *http.Request, error) {
	payload := map[string]any{
		"taskType":        request.TaskType,
		"taskUUID":        request.TaskUUID,
		"positivePrompt":  request.Prompt,
		"width":           request.Width,
		"height":          request.Height,
		"model":           request.Model,
		"numberOfResults": request.NumberOfResults,
		"uploadEndpoint":  request.UploadEndpoint,
		"checkNSFW":       request.CheckNSFW,
		"includeCost":     request.IncludeCost,
		"outputType":      request.OutputType,
		"outputFormat":    request.OutputFormat,
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

func sendRequest(request RunwareOptions, url string) (*RunwareResponseBody, error) {
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
