/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var filepath string
var filterkeys []string

func checkErr(err error, errorMessage string) {
	if err != nil {
		fmt.Printf("Error %s: %v\n", errorMessage, err)
		os.Exit(1)
	}
}

func parseJSON(delim json.Delim) {
	switch delim {
	case '{':
		//single json
	case '[':
		//
	}
}

func checkJSON(decoder *json.Decoder) (json.Delim, error) {

	token, err := decoder.Token()
	if err != nil {
		return 0, err
	}

	delim, ok := token.(json.Delim)
	if !ok {
		return 0, errors.New("top-level JSON is not an object or array")
	}

	if delim != '{' && delim != '[' {
		return 0, fmt.Errorf("unsupported JSON delimiter: %v", delim)
	}

	return delim, nil
}

// func findKey(keys []string, decoder json.Decoder) {
// 	encoder := json.NewEncoder(os.Stdout)

// }

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read from a json",
	Long: `
	Parse throught a json file - Paramaters can be passed to filter data.

	Streams through the json to handle large files - only prints the specific filter passed (if any)
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// check file path
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Error opening JSON: %v\n", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)

		// check for json
		delim, err := checkJSON(decoder)
		checkErr(err, "reading from json")
		parseJSON(delim)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	readCmd.Flags().StringVarP(&filepath, "file", "f", "", "Path to JSON file")

	readCmd.MarkFlagRequired("file")
}
