package main

import (
	"fmt"
	"os"

	"github.com/santoshkal/genval-poc/generate"
	"github.com/santoshkal/genval-poc/validate"
)

// const DockerfilePolicy = "./policies/docker-file.rego"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go input.yaml output.Dockerfile")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Use ParseInputFile to read and unmarshal the input file
	var data struct {
		Dockerfile []struct {
			Stage        int                      `yaml:"stage"`
			Instructions []map[string]interface{} `yaml:"instructions"`
		} `yaml:"dockerfile"`
	}

	err := generate.ParseInputFile(inputPath, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// if err := validate.ValidateYAML(); err != nil {
	// 	fmt.Printf("Validation error: %v\n", err)
	// } else {
	// 	fmt.Println("Validation successful.")
	// }

	dockerfileContent := generate.GenerateDockerfileContent(&data)

	outputData := []byte(dockerfileContent)
	err = os.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		fmt.Println("Error writing Dockerfile:", err)
		return
	}
	fmt.Printf("Generated Dockerfile saved to: %s\n", outputPath)

	err = validate.ValidateDockerfile(string(outputData), validate.DockerfilePolicy)
	// fmt.Printf("Dockerfile JSON: %s\n", generatedDockerfileContent)
	if err != nil {
		fmt.Println("Dockerfile validation failed:", err)
		return
	} else {
		fmt.Println("Dockerfile validation succeeded!")
	}
}
