# runware-go

```runware-go``` is a (unofficial) Go client library for interacting with the Runware API, providing easy access to image generation, inference tasks, and related features.

This library abstracts away HTTP request handling and gives you a clean, configurable interface for generating images using the Runware platform.

## Features

- Easy integration with Runware API
- Supports configuration via a fluent Config method
- Handles request building, headers, and JSON decoding
- Supports flexible output options (Base64, Data URI, URL)
- Built-in error handling on HTTP status codes

## Installation

```bash
go get github.com/ableinc/runware-go
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/ableinc/runware-go"
)

func main() {
	client := runware.NewGenerateImagesV1("YOUR_API_KEY").Config(map[string]any{
		"taskType":    runware.ImageInference,
		"taskUUID":    "task-uuid-123",
		"prompt":      "A dragon flying over mountains",
		"width":       runware.SD_Width,
		"height":      runware.SD_Height,
		"model":       "runware:100@1",
		"results":     int8(1),
		"checkNSFW":   true,
		"includeCost": true,
		"outputType":  runware.URL,
		"outputFormat": runware.PNG,
	})

	resp, err := client.GenerateV1()
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	fmt.Printf("Generated Image URL: %s\n", resp.ImageUrl)
}
```

## Configuration Parameters

The Config() method accepts a map[string]any with the following keys:


|Key             |Type          |Description|
| --------------- | ----------- | --------- | 
|taskType       |TaskType      |Type of task (e.g., ImageInference)|
|taskUUID       |string        |Unique task ID|
|prompt         |string        |Positive prompt description|
|width         |int8          |Width of output image|
|height        |int8          |Height of output image|
|model         |string        |Model name (e.g., dalle3)|
|results        |int8         |Number of results to generate|
|uploadEndpoint |string       |Optional upload endpoint|
|checkNSFW      |bool         |Enable NSFW checking|
|includeCost    |bool         |Include cost information|
|outputType     |OutputType   |Output type (Base64Data, DataURI, URL)|
|outputFormat   |OutputFormat |Output format (PNG, JPG, WEBP)|

## Response Fields

The RunwareResponseBody struct contains:

|Field           |Type      |Description|
|---------       | -------- | ----------|
|TaskType        |string    |Task type|
|TaskUUID        |string    |Task UUID|
|ImageUUID       |string    |Image UUID|
|ImageUrl        |string    |Public image URL|
|ImageBase64Data |string    |Base64-encoded image data|
|ImageDataURI    |string    |Data URI of the image|
|Seed            |int8      |Random seed used|
|NSFWContent     |bool      |Indicates if content was NSFW|
|Cost            |float64   |Cost of generation (if enabled)|

## Authentication

Use your Runware API key when creating a client:

```go
client := runware.NewGenerateImagesV1("YOUR_API_KEY")
```

## Example With Minimal Config

```go
client := runware.NewGenerateImagesV1("YOUR_API_KEY").Config(map[string]any{
	"taskType":     runware.ImageInference,
	"taskUUID":     "065cb06a-41ef-4fb6-a6d6-63dc8b76f189", // if you do not provide a taskUUID and random UUID v4 will be generated automatically
	"prompt":       "A dragon flying over mountains",
	"width":        runware.SD_Landscape16_9Width,
	"height":       runware.SD_Landscape16_9Height,
	"model":        "runware:100@1",
	"results":      int8(1),
	"checkNSFW":    true,
	"includeCost":  true,
	"outputType":   runware.Base64Data,
	"outputFormat": runware.PNG,
})

resp, err := client.GenerateV1()
jsonResponse, err := json.MarshalIndent(resp, "", "  ")
if err != nil {
	log.Fatalf("Failed to marshal response: %v", err)
}
fmt.Printf(string(jsonResponse))
```

**Output**
Note: This script was ran 3 times, these are the 3 images it produced on each run.

![Dragon over mountain 1](https://able.sfo2.cdn.digitaloceanspaces.com/github/dragon_over_mountain_3.jpg "A dragon flying over mountains")

![Dragon over mountain 2](https://able.sfo2.cdn.digitaloceanspaces.com/github/dragon_over_mountain_2.jpg "A dragon flying over mountains")

![Dragon over mountain 3](https://able.sfo2.cdn.digitaloceanspaces.com/github/dragon_over_mountain_1.jpg "A dragon flying over mountains")


## Error Handling

- The library automatically checks for HTTP status codes >= 400.
- If a request fails, you will get an error from GenerateV1().

## Example:

```go
resp, err := client.GenerateV1()
if err != nil {
	log.Fatalf("API error: %v", err)
}
```

## Contributing

Feel free to open issues or PRs to improve this library!

## License

MIT License
