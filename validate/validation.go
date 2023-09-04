package validate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// TODO: Update the path to the OPA binary using embed.FS
const opaBinaryPath = "./opa"

// const opaBinaryName = "opa"

type DockerfileInstruction struct {
	Cmd   string `json:"cmd"`
	Value string `json:"value"`
}

func ParseDockerfileContent(content string) []DockerfileInstruction {
	lines := strings.Split(content, "\n")
	var instructions []DockerfileInstruction

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		cmd := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		instructions = append(instructions, DockerfileInstruction{
			Cmd:   cmd,
			Value: value,
		})
	}

	return instructions
}

func ValidateDockerfile(dockerfileJSON []byte) (string, error) {
	dockerfileInstructions := ParseDockerfileContent(string(dockerfileJSON))

	dockerfileJSON, err := json.Marshal(dockerfileInstructions)
	if err != nil {
		return "", fmt.Errorf("error converting to JSON: %w", err)
	}
	tempFile, err := os.CreateTemp("", "dockerfile.json")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the Dockerfile JSON to the temporary file
	_, err = tempFile.Write(dockerfileJSON)
	if err != nil {
		return "", fmt.Errorf("error writing Dockerfile JSON to temporary file: %w", err)
	}
	tempFile.Close()

	query := "data.dockerfile_validation"

	cmd := exec.Command(opaBinaryPath, "eval", query, "--data", "./security.rego", "--format", "pretty", "--input", tempFile.Name())
	cmd.Stdin = io.Reader(bytes.NewReader(dockerfileJSON))
	cmd.Stderr = os.Stderr
	var output bytes.Buffer
	cmd.Stdout = &output

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error running OPA: %w", err)
	}
	// Process the JSON output as needed
	result := output.String()

	// Parse the OPA JSON output
	var opaOutput map[string]bool
	if err := json.Unmarshal([]byte(result), &opaOutput); err != nil {
		return "", fmt.Errorf("error parsing OPA output: %w", err)
	}

	// Define the list of policies
	// Add more policies as needed
	policies := []string{"latest_base_image", "untrusted_base_image", "deny_root_user", "deny_sudo", "deny_caching", "deny_add", "deny_image_expansion"}

	// Format the output indicating which policies passed and which failed
	var formattedOutput string
	for _, policy := range policies {
		passed, exists := opaOutput[policy]
		if !exists || !passed {
			formattedOutput += fmt.Sprintf("Policy '%s:' failed.\n", policy)
		} else {
			formattedOutput += fmt.Sprintf("Policy '%s:' passed.\n", policy)
		}
	}

	return formattedOutput, nil
}
