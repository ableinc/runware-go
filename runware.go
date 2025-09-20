package runware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/google/uuid"
)

type TaskType string
type OutputType string
type OutputFormat string
type Definition uint16

const (
	ImageInference TaskType     = "imageInference"
	Base64Data     OutputType   = "base64Data"
	DataURI        OutputType   = "dataURI"
	URL            OutputType   = "URL"
	PNG            OutputFormat = "PNG"
	JPG            OutputFormat = "JPEG"
	WEBP           OutputFormat = "WEBP"
	SD_Height      Definition   = 512
	SD_Width       Definition   = 512

	SD_Portrait3_4Height   Definition = 1024
	SD_Portrait3_4Width    Definition = 768
	SD_Portrait9_16Height  Definition = 1152
	SD_Portrait9_16Width   Definition = 640
	SD_Landscape4_3Height  Definition = 768
	SD_Landscape4_3Width   Definition = 1024
	SD_Landscape16_9Height Definition = 640
	SD_Landscape16_9Width  Definition = 1152

	HD_Height              Definition = 1024
	HD_Width               Definition = 1024
	HD_Portrait3_4Height   Definition = 1536
	HD_Portrait3_4Width    Definition = 1152
	HD_Portrait9_16Height  Definition = 1728
	HD_Portrait9_16Width   Definition = 960
	HD_Landscape4_3Height  Definition = 1152
	HD_Landscape4_3Width   Definition = 1536
	HD_Landscape16_9Height Definition = 960
	HD_Landscape16_9Width  Definition = 1728
)

type RunwareOptions struct {
	TaskType        TaskType     `json:"taskType"`
	TaskUUID        string       `json:"taskUUID"`
	Prompt          string       `json:"prompt"`
	Model           string       `json:"model"`
	UploadEndpoint  string       `json:"uploadEndpoint"`
	OutputType      OutputType   `json:"outputType"`
	OutputFormat    OutputFormat `json:"outputFormat"`
	Width           Definition   `json:"width"`
	Height          Definition   `json:"height"`
	NumberOfResults uint8        `json:"numberOfResults"`
	CheckNSFW       bool         `json:"checkNSFW"`
	IncludeCost     bool         `json:"includeCost"`
}

type RunwareSuccessResponseBody struct {
	TaskType        string  `json:"taskType"`
	TaskUUID        string  `json:"taskUUID"`
	ImageUUID       string  `json:"imageUUID"`
	ImageUrl        string  `json:"imageUrl"`
	ImageBase64Data string  `json:"imageBase64Data"`
	ImageDataURI    string  `json:"imageDataURI"`
	Seed            int     `json:"seed"`
	Cost            float64 `json:"cost"`
	NSFWContent     bool    `json:"nsfwContent"`
}

type RunwareErrorResponseBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Parameter string `json:"parameter"`
	Type      string `json:"type"`
	TaskType  string `json:"taskType"`
}

type RunwareResponseBody struct {
	Data   []RunwareSuccessResponseBody
	Errors []RunwareErrorResponseBody
}

// Interface definition
type GenerateImagesV1 interface {
	Config(data []map[string]any) GenerateImagesV1
	GenerateV1() (*[]RunwareSuccessResponseBody, error)
}

// Struct implementing the interface
type generateImagesV1Impl struct {
	apiKey        string
	options       []RunwareOptions
	omittedFields []string
}

func NewGenerateImagesV1(apiKey string) GenerateImagesV1 {
	return &generateImagesV1Impl{
		apiKey:        apiKey,
		omittedFields: []string{},
	}
}

func (g *generateImagesV1Impl) Config(options []map[string]any) GenerateImagesV1 {
	g.options = make([]RunwareOptions, len(options))
	for i, data := range options {
		if data["taskType"] != nil {
			g.options[i].TaskType = data["taskType"].(TaskType)
		}
		if data["taskUUID"] != nil {
			g.options[i].TaskUUID = data["taskUUID"].(string)
		} else {
			g.options[i].TaskUUID = uuid.New().String()
		}
		if data["prompt"] != nil {
			g.options[i].Prompt = data["prompt"].(string)
		}
		if data["width"] != nil {
			g.options[i].Width = data["width"].(Definition)
		}
		if data["height"] != nil {
			g.options[i].Height = data["height"].(Definition)
		}
		if data["model"] != nil {
			g.options[i].Model = data["model"].(string)
		}
		if data["results"] != nil {
			g.options[i].NumberOfResults = data["results"].(uint8)
		}
		if data["uploadEndpoint"] != nil {
			g.options[i].UploadEndpoint = data["uploadEndpoint"].(string)
		}
		if data["checkNSFW"] != nil {
			g.options[i].CheckNSFW = data["checkNSFW"].(bool)
		} else {
			g.omittedFields = append(g.omittedFields, "checkNSFW")
		}
		if data["includeCost"] != nil {
			g.options[i].IncludeCost = data["includeCost"].(bool)
		} else {
			g.omittedFields = append(g.omittedFields, "includeCost")
		}
		if data["outputType"] != nil {
			g.options[i].OutputType = data["outputType"].(OutputType)
		}
		if data["outputFormat"] != nil {
			g.options[i].OutputFormat = data["outputFormat"].(OutputFormat)
		}
	}
	return g
}

func (g *generateImagesV1Impl) GenerateV1() (*[]RunwareSuccessResponseBody, error) {
	var v1Domain string = "https://api.runware.ai/v1"
	return sendRequest(g, v1Domain)
}

func buildClient(g *generateImagesV1Impl, url string) (*http.Client, *http.Request, error) {
	var payload []map[string]any = make([]map[string]any, 0)
	for _, request := range g.options {
		width, err := getDimensionValue(request.Width)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid width: %w", err)
		}
		request.Width = Definition(width)
		height, err := getDimensionValue(request.Height)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid height: %w", err)
		}
		request.Height = Definition(height)
		payload = append(payload, skipEmptyOrNil(map[string]any{
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
		}, g.omittedFields))
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.apiKey))
	return client, req, nil
}

func sendRequest(g *generateImagesV1Impl, url string) (*[]RunwareSuccessResponseBody, error) {
	client, req, err := buildClient(g, url)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response RunwareResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		log.Printf("request failed with status %d", resp.StatusCode)
		jsonDataErrResponse, err := json.MarshalIndent(response, "", "  ")
		if err == nil {
			return nil, fmt.Errorf("%s", jsonDataErrResponse)
		}
	}
	return &response.Data, nil
}

func getDimensionValue(dim any) (int16, error) {
	switch v := dim.(type) {
	case Definition:
		return int16(v), nil
	default:
		return 0, fmt.Errorf("invalid dimension type, must be Definition or Definition")
	}
}

func skipEmptyOrNil(option map[string]any, omittedFields []string) map[string]any {
	for key, value := range option {
		if value == "" || value == nil || slices.Contains(omittedFields, key) {
			delete(option, key)
		}
	}
	return option
}
