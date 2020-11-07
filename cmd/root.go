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
	"fmt"
	"github.com/spf13/cobra"
  "io"
	"os"
  "strings"
  
  "github.com/emla2805/tfr/utils"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfr",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	RunE: func(cmd *cobra.Command, args []string) error {
    if isInputFromPipe() {
        // if input is from a pipe, upper case the
        // content of stdin
        print("data is from pipe")
        return parseRecords(os.Stdin)
    } else {
        // ...otherwise get the file
        // file, e := getFile()
        // if e != nil {
        //     return e
        // }
        // defer file.Close()
        file, e := os.Open(args[0])
        if e != nil {
          return e
        }
        defer file.Close()
        return parseRecords(file)
    }
  },
}

func isInputFromPipe() bool {
    fileInfo, _ := os.Stdin.Stat()
    return fileInfo.Mode() & os.ModeCharDevice == 0
}


func toUppercase(r io.Reader, w io.Writer) error {
    scanner := bufio.NewScanner(bufio.NewReader(r))
    for scanner.Scan() {
        _, e := fmt.Fprintln(
            w, strings.ToUpper(scanner.Text()))
        if e != nil {
            return e
        }
    }
    return nil
}

func parseRecords(r io.Reader) error {
  reader := utils.NewReader(r)
  count := 0

  for {
    example, e := reader.Next()

    if e == io.EOF {
      break
    } else if e != nil {
      return e
    }

    id := string(utils.ExampleFeatureBytes(example, "id"))
    fmt.Fprintln(os.Stdout, id)

    count++
  }
  return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
