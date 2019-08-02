package models

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// executors can be nested. for example, Copy is nested within SystemD
type Copy struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type SystemD struct {
	Copy        *Copy  `yaml:"copy,omitempty"` // the caller can also, just restart the service w/o copying.
	ServiceName string `yaml:"serviceName"`
}

// Target defines the plan target
type Target struct {
	Hosts              []string `yaml:"hosts"`
	Serial             bool     `yaml:"serial"`
	ExitOnFirstFailure bool     `yaml:"exitOnFirstFailure"`
}

type BaseStep struct {
	operation string `yaml:"op"`
}

type RestartOp struct {
	Restart string `yaml:"restart"`
}

type SystemDStep struct {
	BaseStep
	serviceName string
}

// Plan defines the procedure of how to update the given plan
type Plan struct {
	Name   string                 `yaml:"name"`           //name of the plan
	User   string                 `yaml:"user,omitempty"` // default root
	Target Target                 `yaml:"target"`
	Steps  []map[string]Operation `yaml:"steps"`
}

func EncodeToObj(obj interface{}, i interface{}) error {
	out, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(out, i)
	if err != nil {
		return err
	}
	return nil
}

// Parse parses file into obj
func Parse(file string, obj interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, obj)
	if err != nil {
		return err
	}
	return nil
}
