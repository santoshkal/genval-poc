package validate

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
	"github.com/santoshkal/genval-poc/pkg/parser"
	log "github.com/sirupsen/logrus"
)

//go:embed dockerFilePolicies.rego
var dockerPolicy []byte

const (
	// DockerfilePolicy is the name of the Rego policy file used to validate Dockerfiles.
	DockerfilePolicy = "dockerFilePolicies.rego"

	// DockerfilePackage is the name of the Rego package containing Dockerfile validation policies.
	DockerfilePackage = "data.dockerfile_validation"
)

// ValidateDockerfileUsingRego validates a Dockerfile using Rego.
func ValidateDockerfile(dockerfileContent string, regoPolicyPath string) error {

	// Prepare Rego input data
	dockerfileInstructions := parser.ParseDockerfileContent(dockerfileContent)

	jsonData, err := json.Marshal(dockerfileInstructions)
	if err != nil {
		log.WithError(err).Error("Error converting to JSON:", err)
		return errors.New("error converting to JSON")
	}

	var commands []map[string]string
	err = json.Unmarshal([]byte(jsonData), &commands)
	if err != nil {
		errWithContext := fmt.Errorf("error converting JSON to map: %v", err)
		log.WithError(err).Error(errWithContext.Error())
		return errWithContext
	}

	ctx := context.Background()

	// Create regoQuery for evaluation
	regoQuery := rego.New(
		rego.Query(DockerfilePackage),
		rego.Module(DockerfilePolicy, string(dockerPolicy)),
		rego.Input(commands),
	)

	// Evaluate the Rego query
	rs, err := regoQuery.Eval(ctx)
	if err != nil {
		log.Fatal("Error evaluating query:", err)
	}

	var policyError error

	for _, result := range rs {
		if len(result.Expressions) > 0 {
			keys := result.Expressions[0].Value.(map[string]interface{})
			for key, value := range keys {
				// log.Infof("The Key is: %v of Type:%T and Value is: %v of Type %T", key, key, value, value)
				if value != "" {
					log.Printf("%v", value)
					fmt.Printf("Dockerfile validation policy: %s passed\n", key)
				}
				if value != true {
					log.Errorf("Dockerfile validation policy: %s failed\n", key)
					policyError = fmt.Errorf("dockerfile validation policy: %s failed", key)
				}
			}
		} else {
			log.Error("No policies passed")
		}
	}

	if err != nil {
		log.WithError(err).Error("Error evaluating Rego.")
		return errors.New("error evaluating Rego")
	}

	return policyError
}
