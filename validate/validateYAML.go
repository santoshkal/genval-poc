package validate

import (
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func ValidateYAML() error {
	ctx := cuecontext.New()

	// load schema.cue
	schemaBytes, err := os.ReadFile("input1.cue")
	if err != nil {
		return err
	}
	schema := ctx.CompileBytes(schemaBytes, cue.Filename("input1.cue"))
	if err := schema.Err(); err != nil {
		return err
	}

	// load input.json
	dataBytes, err := os.ReadFile("input1.json")
	if err != nil {
		return err
	}
	data := ctx.CompileBytes(dataBytes, cue.Filename("input1.json"))
	if err := data.Err(); err != nil {
		return err
	}

	// use #Dockerfile from the schema
	schema = schema.LookupPath(cue.ParsePath("#Dockerfile"))
	if err := schema.Err(); err != nil {
		return err
	}

	// unify the schema with the input and validate, like `cue vet`
	v := schema.Unify(data)
	if err := v.Err(); err != nil {
		return err
	}
	if err := v.Validate(); err != nil {
		return err
	}

	return nil
}
