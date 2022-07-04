package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// {
//     "name": "${method}-debug",
//     "type": "go",
//     "request": "launch",
//     "program": "${workspaceFolder}/{vscode-path}",
//     "env": {
//         "_LAMBDA_SERVER_PORT": "${port}"
//     },
//     "args": []
// },

type VSCodeLaunchConfiguration struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Request string            `json:"request"`
	Program string            `json:"program"`
	Env     map[string]string `json:"env"`
	Args    []string          `json:"args,omitempty"`
}

type VSCodeLaunchCompound struct {
	Name           string   `json:"name"`
	PreLaunchTask  string   `json:"preLaunchTask"`
	Configurations []string `json:"configurations,omitempty"`
}

type XLocalGW struct {
	Port          int64  `json:port,omitempty`
	VSCodePath    string `json:"vscode-path,omitempty"`
	VSCodeProject string `json:"vscode-project,omitempty"`
}

func GenerateVSCodeLaunch(f string) error {
	specDoc, err := loads.Spec(f)
	if err != nil {
		return err
	}
	var cfg []*VSCodeLaunchConfiguration

	for k, p := range specDoc.Spec().Paths.Paths {
		// fmt.Printf("%v\n", k)

		cfg = keep(cfg, specDoc.Spec(), k, &p, p.Delete)
		cfg = keep(cfg, specDoc.Spec(), k, &p, p.Get)
		cfg = keep(cfg, specDoc.Spec(), k, &p, p.Head)
		cfg = keep(cfg, specDoc.Spec(), k, &p, p.Patch)
		cfg = keep(cfg, specDoc.Spec(), k, &p, p.Post)
		cfg = keep(cfg, specDoc.Spec(), k, &p, p.Put)
	}

	comp := &VSCodeLaunchCompound{
		Name:          "Run Local Gateway",
		PreLaunchTask: "run-local-gw-task",
	}

	fmt.Println()
	for _, c := range cfg {
		comp.Configurations = append(comp.Configurations, c.Name)
		data, err := json.MarshalIndent(c, "", "   ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v,\n", string(data))
	}

	fmt.Println()
	fmt.Println()
	data, err := json.MarshalIndent(comp, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v,\n", string(data))

	return nil
}

func keep(cfg []*VSCodeLaunchConfiguration, swag *spec.Swagger, path string, item *spec.PathItem, op *spec.Operation) []*VSCodeLaunchConfiguration {
	if op == nil {
		return cfg
	}

	con := GenerateVSCodeLaunchConfiguration(swag, path, item, op)
	if con != nil {
		return append(cfg, con)
	}

	return cfg

}

func GenerateName(swag *spec.Swagger, path string, item *spec.PathItem, op *spec.Operation) string {
	// fmt.Printf("%15v - %30v - %v\n", swag.Info.Title, path, op.ID)
	title := FixName(swag.Info.Title)
	p := FixName(path)
	opName := FixName(op.ID)
	// fmt.Printf("%v-%v-%v\n", title, p, opName)
	return fmt.Sprintf("%v%v:%v", title, p, opName)
}

func FixName(name string) string {
	fixed := strings.ToLower(strings.TrimSpace(name))
	fixed = strings.ReplaceAll(fixed, " ", "-")
	return fixed
}

func GetExtension(op *spec.Operation) *XLocalGW {
	extension, found := op.Extensions["x-localgw"]
	if !found {
		fmt.Printf("localgw extension not found for  %v -- skipping\n", op.ID)
		return nil
	}

	// This is really a map[string]interface
	var ext XLocalGW

	data, _ := json.Marshal(extension)
	err := json.Unmarshal(data, &ext)

	if err != nil {
		fmt.Printf("Malformed x-localgw exension for %v  -- skipping, err : %v\n", op.ID, err)
		return nil
	}
	return &ext
}

func GenerateVSCodeLaunchConfiguration(swag *spec.Swagger, path string, item *spec.PathItem, op *spec.Operation) *VSCodeLaunchConfiguration {
	name := GenerateName(swag, path, item, op)

	extension, found := op.Extensions["x-localgw"]
	if !found {
		fmt.Printf("localgw extension not found for %v : %v : %v -- skipping\n", swag.Info.Title, path, op.ID)
		return nil
	}

	// This is really a map[string]interface
	var ext XLocalGW

	data, _ := json.Marshal(extension)
	err := json.Unmarshal(data, &ext)

	if err != nil {
		fmt.Printf("Malformed x-localgw exension for %v : %v : %v -- skipping, err : %v\n", swag.Info.Title, path, op.ID, err)
		return nil
	}

	program := ext.VSCodePath

	env := make(map[string]string)
	env["_LAMBDA_SERVER_PORT"] = fmt.Sprintf("%v", ext.Port)

	rtn := &VSCodeLaunchConfiguration{
		Name:    name,
		Type:    "go",
		Request: "launch",
		Program: program,
		Env:     env,
	}

	return rtn
}
