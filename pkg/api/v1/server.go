package api

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	generate "github.com/santoshkal/genval-poc/pkg/generate/dockerfile-gen"
	validate "github.com/santoshkal/genval-poc/pkg/validate/dockerfile-val"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var log = logrus.New()

func init() {
	// Open a file for logging
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Debug("Failed to log to file, using default stderr")
	}
}

func SetupAPI() *gin.Engine {

	r := gin.Default()

	r.POST("/generate", func(ctx *gin.Context) {
		// log.Infof("Remote Address: %s", ctx.Request.RemoteAddr)
		log.Println("Request Headers:")
		for key, value := range ctx.Request.Header {
			log.Printf("%s: %s\n", key, value)
		}
		// Read and parse the request body (JSON/YAML)
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			log.Infof("This is Request Header: %v", ctx.Request.Header)
			return
		}

		yamlContent := string(body)

		// Validate the input using OPA
		err = validate.ValidateInput(yamlContent, validate.InputPolicy)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "message": err.Error()})
			return
		}

		var data generate.DockerfileContent

		// ReadAndParseFile parses the body into DockerfileContent struct
		err = ParseRequestBody(body, &data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse YAML content", "message": err.Error()})
			return
		}

		// Generate the Dockerfile
		dockerfileContent := generate.GenerateDockerfileContent(&data)

		// Validate the generated Dockerfile using OPA/Rego policies
		err = validate.ValidateDockerfile(dockerfileContent, validate.DockerfilePolicy)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dockerfile validation failed", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Dockerfile validation succeeded!", "dockerfile": dockerfileContent})
	})

	return r
}

func ParseRequestBody(content []byte, result interface{}) error {
	// Unmarshal (decode) the content into the result interface.
	return yaml.Unmarshal(content, result)
}
