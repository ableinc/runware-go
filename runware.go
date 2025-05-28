package runware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	TaskType        TaskType
	TaskUUID        string
	Prompt          string
	Model           string
	UploadEndpoint  string
	OutputType      OutputType
	OutputFormat    OutputFormat
	Width           Definition
	Height          Definition
	NumberOfResults uint8
	CheckNSFW       bool
	IncludeCost     bool
}

type RunwareSuccessResponseBody struct {
	TaskType        string
	TaskUUID        string
	ImageUUID       string
	ImageUrl        string
	ImageBase64Data string
	ImageDataURI    string
	Seed            int
	Cost            float64
	NSFWContent     bool
}

type RunwareErrorResponseBody struct {
	Code      string
	Message   string
	Parameter string
	Type      string
	TaskType  string
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
	apiKey  string
	options []RunwareOptions
}

func NewGenerateImagesV1(apiKey string) GenerateImagesV1 {
	return &generateImagesV1Impl{
		apiKey: apiKey,
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
		}
		if data["includeCost"] != nil {
			g.options[i].IncludeCost = data["includeCost"].(bool)
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
	return sendRequest(g.apiKey, g.options, v1Domain)
}

func buildClient(apiKey string, requests []RunwareOptions, url string) (*http.Client, *http.Request, error) {
	var payload []map[string]any = make([]map[string]any, 0)
	for _, request := range requests {
		width, err := getDimensionValue(request.Width)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid width: %w", err)
		}
		height, err := getDimensionValue(request.Height)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid height: %w", err)
		}
		payload = append(payload, skipEmptyOrNil(map[string]any{
			"taskType":        request.TaskType,
			"taskUUID":        request.TaskUUID,
			"positivePrompt":  request.Prompt,
			"width":           width,
			"height":          height,
			"model":           request.Model,
			"numberOfResults": request.NumberOfResults,
			"uploadEndpoint":  request.UploadEndpoint,
			"checkNSFW":       request.CheckNSFW,
			"includeCost":     request.IncludeCost,
			"outputType":      request.OutputType,
			"outputFormat":    request.OutputFormat,
		}))
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	return client, req, nil
}

func sendRequest(apiKey string, requests []RunwareOptions, url string) (*[]RunwareSuccessResponseBody, error) {
	client, req, err := buildClient(apiKey, requests, url)
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

func skipEmptyOrNil(option map[string]any) map[string]any {
	for key, value := range option {
		if value == "" || value == nil {
			delete(option, key)
		}
	}
	return option
}
