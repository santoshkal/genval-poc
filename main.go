// main.go
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/santoshkal/genval-poc/generate"
	"github.com/santoshkal/genval-poc/validate"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go input.yaml output.Dockerfile")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error reading input YAML:", err)
		return
	}

	var data struct {
		Dockerfile []struct {
			Stage        int                      `yaml:"stage"`
			Instructions []map[string]interface{} `yaml:"instructions"`
		} `yaml:"dockerfile"`
	}

	err = yaml.Unmarshal(inputData, &data)
	if err != nil {
		fmt.Println("Error parsing YAML:", err)
		return
	}
	if err := validate.ValidateYAML(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("Validation successful.")
	}
	dockerfileContent := generate.GenerateDockerfileContent(&data)

	outputData := []byte(dockerfileContent)
	err = os.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		fmt.Println("Error writing Dockerfile:", err)
		return
	}

	output, err := validate.ValidateDockerfile(outputData)
	if err != nil {
		fmt.Println("Error validating Dockerfile:", err)
		return
	}
	// Process the OPA output as needed
	fmt.Println("Dockerfile generated and validated successfully!\n")
	fmt.Printf("OPA Evaluation:\n%v\n", output)
}
