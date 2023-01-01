package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itchyny/gojq"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/yaml"
)

const yamlFile = "./../examples/td2-templ.yaml"
const exp = `$topoDef | [.spec.properties.templates[].templateRef.name]`

func main() {
	zlog := zap.New(zap.UseDevMode(true), zap.JSONEncoder())
	logger := logging.NewLogrLogger(zlog.WithName("lcnc runtime"))

	fb, err := os.ReadFile(yamlFile)
	if err != nil {
		logger.Debug("cannot read file", "error", err)
		os.Exit(1)
	}
	//logger.Debug("read file")

	j, err := yaml.YAMLToJSON(fb)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	//fmt.Println(string(j))

	def := map[string]interface{}{}
	if err := json.Unmarshal(j, &def); err != nil {
		logger.Debug("cannot unmarshal", "error", err)
		os.Exit(1)
	}
	//logger.Debug("unmarshal succeeded")
	//logger.Debug("definition", "def", def)

	inputVar := "$topoDef"
	inputVal := def

	q, err := gojq.Parse(exp)
	if err != nil {
		logger.Debug("cannot parse jq", "error", err)
		os.Exit(1)
	}
	code, err := gojq.Compile(q, gojq.WithVariables([]string{inputVar}))
	if err != nil {
		logger.Debug("cannot compile jq", "error", err)
	}
	iter := code.Run(nil, inputVal)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			logger.Debug("jq result error", "error", err)
		os.Exit(1)
		}
		fmt.Printf("%v\n", v)
	}
}
