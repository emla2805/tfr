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
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"io"
	"math"
	"os"

	protobuf "github.com/emla2805/tfr/protobuf"
	"github.com/emla2805/tfr/utils"
)

var numberRecords int
var record string

var rootCmd = &cobra.Command{
	Use:   "tfr {file ... | -}",
	Short: "A lightweight commandline TFRecords processor",
	Long: `tfr is a lightweight command-line TFRecords processor that 
reads serialized .tfrecord files and outputs results as JSON on standard output.`,
	Example: `  $ tfr data_tfrecord-00000-of-00001
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
  }`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 || isInputFromPipe() {
			return nil
		}
		return errors.New("requires argument or stdin")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var readers []io.Reader
		if isInputFromPipe() {
			readers = append(readers, os.Stdin)
		}

		for _, path := range args {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			readers = append(readers, file)
		}
		multi := io.MultiReader(readers...)

		scanner := bufio.NewScanner(multi)
		scanner.Split(utils.ScanTFRecord)

		var example proto.Message
		switch record {
		case "example":
			example = &protobuf.Example{}
		case "sequence_example":
			example = &protobuf.SequenceExample{}
		default:
			example = &protobuf.Example{}
		}

		count := 0
		for scanner.Scan() {
			if count >= numberRecords {
				break
			}
			err := proto.Unmarshal(scanner.Bytes(), example)
			if err != nil {
				return err
			}

			jsonBytes, err := utils.Marshal(example)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(jsonBytes))

			count++
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&numberRecords, "number", "n", math.MaxInt32, "number of records to show")
	rootCmd.Flags().StringVarP(&record, "record", "r", "example", "record type { example | sequence_example }")
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}
