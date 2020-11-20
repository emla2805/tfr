/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/spf13/cobra"

	tfr "github.com/emla2805/tfr/protobuf"
	"github.com/emla2805/tfr/utils"
)

var numberRecords int

var rootCmd = &cobra.Command{
	Use:   "tfr",
	Short: "A lightweight commandline TFRecords processor",
	Long: `tfr is a lightweight command-line TFRecords processor that 
reads serialized .tfrecord files and outputs to stdout in JSON format.

Example:
$ tfr data_tfrecord-00000-of-00001
{
  "features": {
    "feature": {
      "a": {
        "bytesList": {
          "value": [
            "some text"
          ]
        }
      },
      "label": {
        "int64List": {
          "value": [
            0
          ]
        }
      }
    }
  }
}
  `,

	RunE: func(cmd *cobra.Command, args []string) error {
		if isInputFromPipe() {
			return parseRecords(os.Stdin)
		} else {
			file, e := os.Open(args[0])
			if e != nil {
				return e
			}
			defer file.Close()
			return parseRecords(file)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Example struct {
	Features map[string]Feature `json:"features"`
}

func NewExample(example *tfr.Example) *Example {
	feature := make(Feature)
	ex := &Example{
		Features: make(map[string]Feature),
	}
	for k, v := range example.Features.Feature {
		feature[k] = v
	}
	ex.Features["feature"] = feature
	return ex
}

type Feature map[string]*tfr.Feature

func (f Feature) MarshalJSON() ([]byte, error) {
	output := make(map[string]interface{})
	for k, v := range f {
		var typ string
		value := make(map[string]interface{})
		typMap := make(map[string]interface{})
		switch t := v.Kind.(type) {
		case *tfr.Feature_BytesList:
			typ = "bytesList"
			stringList := []string{}
			for _, byts := range t.BytesList.Value {
				stringList = append(stringList, string(byts))
			}
			value["value"] = stringList
		case *tfr.Feature_FloatList:
			typ = "floatList"
			value["value"] = t.FloatList.Value
		case *tfr.Feature_Int64List:
			typ = "int64List"
			value["value"] = t.Int64List.Value
		}
		typMap[typ] = value
		output[k] = typMap
	}

	return json.Marshal(output)
}

func init() {
	rootCmd.Flags().IntVarP(&numberRecords, "number", "n", math.MaxInt32, "Number of records to show")
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func parseRecords(r io.Reader) error {
	reader := utils.NewReader(r)
	count := 0

	for {
		example, e := reader.Next()

		if count >= numberRecords {
			break
		}

		if e == io.EOF {
			break
		} else if e != nil {
			return e
		}

		// convert to local example
		ex := NewExample(example)
		j, err := json.Marshal(ex)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(j))
		count++
	}
	return nil
}
