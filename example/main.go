package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/rego"
)

const ourEnforceOutput = `
	package talk_opa

	result = {
		"admin": admin,
		"allow": allow,
		"deny": deny,
	}

	default admin = false
	default allow = false
	default deny = false`

const clientLoginPolicy = `
	package talk_opa

	teams := input.session.teams

	admin { teams[_] == "DevOps" }
	allow { teams[_] == "Engineering" }
	deny  { not input.session.member }`

func main() {
	ctx := context.Background()

	// Do not forget about UnsafeBuiltins
	base := rego.New(rego.Query("data.talk_opa.result"),
		rego.Module("defaults", ourEnforceOutput),
		rego.Module("policy", clientLoginPolicy))

	preparedQuery, err := base.PrepareForEval(ctx)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	inputRaw := `{
		"session": {
			"teams": ["DevOps"],
			"member": true
		}
	}`

	var input interface{}
	err = json.NewDecoder(bytes.NewBufferString(inputRaw)).Decode(&input)
	if err != nil {
		os.Exit(2)
	}

	result, err := preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(3)
	}

	fmt.Printf("%v\n", result)
}
