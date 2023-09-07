package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/open-policy-agent/opa/rego"
)

const (
	DockerfilePolicy  = "./policies/docker-file.rego"
	DockerfilePackage = "data.dockerfile_validation"
)

// DockerfileInstruction represents a Dockerfile instruction with Cmd and Value.
type DockerfileInstruction struct {
	Cmd   string `json:"cmd"`
	Value string `json:"value"`
}

// ...

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

// ValidateDockerfileUsingRego validates a Dockerfile using Rego.
func ValidateDockerfile(dockerfileContent string, regoPolicyPath string) error {
	// Read Rego policy code from file
	regoPolicyCode, err := os.ReadFile(regoPolicyPath)
	if err != nil {
		return fmt.Errorf("error reading rego policy: %v", err)
	}

	// Prepare Rego input data
	dockerfileInstructions := ParseDockerfileContent(dockerfileContent)

	jsonData, err := json.Marshal(dockerfileInstructions)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return nil
	}

	var commands []map[string]string
	err = json.Unmarshal([]byte(jsonData), &commands)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Create regoQuery for evaluation
	regoQuery := rego.New(
		rego.Query(DockerfilePackage),
		rego.Module(DockerfilePolicy, string(regoPolicyCode)),
		rego.Input(commands),
	)

	// Evaluate the Rego query
	rs, err := regoQuery.Eval(context.Background())
	if err != nil {
		log.Fatal("Error evaluating query:", err)
	}

	// Iterate over the resultSet and print the result metadata
	for _, result := range rs {
		if len(result.Expressions) > 0 {
			keys := result.Expressions[0].Value.(map[string]interface{})
			for key, value := range keys {
				if value != true {
					fmt.Printf("Dockerfile validation policy: %s failed\n", key)
				} else {
					fmt.Printf("Dockerfile validation policy: %s passed\n", key)
				}
			}
		} else {
			fmt.Println("No policies passed")
		}
	}

	if err != nil {
		return fmt.Errorf("error evaluating Rego: %v", err)
	}

	return nil
}
