package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/itchyny/gojq"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/yaml"
)

const dir = "./../examples/templates"
//const exp = `.[] | select(.metadata.name == $VALUE)`
const exp = `.[] | select(.metadata.name == $VALUE)`

//select( .author as $a | ["Gary", "Larry"] | index($a) )

//const exp = `.`

func main() {
	zlog := zap.New(zap.UseDevMode(true), zap.JSONEncoder())
	logger := logging.NewLogrLogger(zlog.WithName("lcnc runtime"))

	files, err := os.ReadDir(dir)
	if err != nil {
		logger.Debug("cannot read dir", "error", err)
		os.Exit(1)
	}

	templates := make([]interface{}, 0, len(files))
	for _, f := range files {
		fb, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			logger.Debug("cannot read file", "error", err)
			os.Exit(1)
		}
		j, err := yaml.YAMLToJSON(fb)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		ji := map[string]interface{}{}
		if err := json.Unmarshal(j, &ji); err != nil {
			logger.Debug("cannot unmarshal", "error", err)
			os.Exit(1)
		}
		templates = append(templates, ji)
	}
	//alltemplates := map[string]interface{}{
	//	"templates": templates,
	//}
	/*
		for _, t := range tmpls {
			logger.Debug("template", "t", t)
		}
	*/

	vals := []string{"tmpl1", "tmpl2"}
	for _, val := range vals {
		//inputVar1 := "$extraInput"
		inputVar := "$VALUE"

		q, err := gojq.Parse(exp)
		if err != nil {
			logger.Debug("cannot parse jq", "error", err)
			os.Exit(1)
		}
		code, err := gojq.Compile(q, gojq.WithVariables([]string{inputVar}))
		if err != nil {
			logger.Debug("cannot compile jq", "error", err)
		}
		iter := code.Run(templates, val)
		for {
			// we can check for 1 result or no result as this should produce
			// an error as we want these maps to be there
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
}
