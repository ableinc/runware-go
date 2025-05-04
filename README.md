# Runware-go

A simple wrapper around Runware.ai API to generate images

## How to Use

- Add to your project

```bash
go get github.com/ableinc/runware-go
```

- Usage

```go
import "github.com/ableinc/runware-go

gen := runware.NewGenerateImagesV1("YOUR_API_KEY").
	Config("imageInference", "task-uuid-123", "A dragon flying over mountains", 512, 512, "civitai:102438@133677", 1)

resp, err := gen.Generate()
if err != nil {
	fmt.Println("Error:", err)
} else {
	fmt.Printf("Image URL: %s\n", resp.ImageUrl)
}
```