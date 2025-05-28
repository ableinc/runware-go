package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	goenv "github.com/ableinc/go-env"
	"github.com/ableinc/runware-go"
)

func main() {
	// Load environment variables
	goenv.LoadEnv(".env", false)
	// Generate image with Runware
	apiKey := os.Getenv("RUNWARE_API_KEY")
	// Note: taskUUID will be generated automatically, if not explicitly provided
	client := runware.NewGenerateImagesV1(apiKey).Config([]map[string]any{
		{
			"taskType":     runware.ImageInference,
			"taskUUID":     "065cb06a-41ef-4fb6-a6d6-63dc8b76f189",
			"prompt":       "A dragon flying over mountains",
			"width":        runware.SD_Landscape16_9Width,
			"height":       runware.SD_Landscape16_9Height,
			"model":        "runware:100@1",
			"results":      uint8(1),
			"checkNSFW":    true,
			"includeCost":  true,
			"outputType":   runware.Base64Data,
			"outputFormat": runware.PNG,
		},
	})
	resp, err := client.GenerateV1()
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}
	jsonResponse, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	fmt.Printf("Response: %s\n", string(jsonResponse))
	for _, image := range *resp {
		if err := saveBase64Image(image.ImageBase64Data, "output.jpg"); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func saveBase64Image(base64Str, outputPath string) error {
	// Decode the base64 string
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %w", err)
	}

	// Write to file
	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
