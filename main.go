package main

import (
	"flag"
	"fmt"
	"os"

	api "github.com/santoshkal/genval-poc/pkg/api/v1"
	generate "github.com/santoshkal/genval-poc/pkg/generate/dockerfile-gen"
	"github.com/santoshkal/genval-poc/pkg/parser"
	validate "github.com/santoshkal/genval-poc/pkg/validate/dockerfile-val"

	log "github.com/sirupsen/logrus"
)

var serverMode bool

func init() {
	flag.BoolVar(&serverMode, "server", false, "Run in server mode")
	flag.Parse()
}
func main() {
	if serverMode {
		runServer()
		return
	}
	if len(os.Args) != 3 {
		log.Debug("Usage: go run main.go input.json output.Dockerfile")
		return
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Use ParseInputFile to read and unmarshal the input file
	var data generate.DockerfileContent

	err := parser.ReadAndParseFile(inputPath, &data)
	if err != nil {
		log.Error("Error:", err)
		return
	}

	yamlContent, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Validate the YAML using OPA
	err = validate.ValidateInput(string(yamlContent), validate.InputPolicy)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
		return
	}

	dockerfileContent := generate.GenerateDockerfileContent(&data)

	outputData := []byte(dockerfileContent)
	err = os.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		log.Error("Error writing Dockerfile:", err)
		return
	}
	fmt.Printf("Generated Dockerfile saved to: %s\n", outputPath)

	err = validate.ValidateDockerfile(string(outputData), validate.DockerfilePolicy)
	// fmt.Printf("Dockerfile JSON: %s\n", generatedDockerfileContent)
	if err != nil {
		log.Error("Dockerfile validation failed:", err)
		return
	} else {
		fmt.Printf("Dockerfile validation succeeded!\n")
	}

}
func runServer() {
	r := api.SetupAPI()
	port := ":3333"
	serverAddr := "localhost" + port
	fmt.Printf("API server is running on port %s...\n", port)

	if err := r.Run(serverAddr); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
