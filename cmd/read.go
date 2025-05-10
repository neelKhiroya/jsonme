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

func printJSON(obj map[string]interface{}, filterKeys []string) {
	if len(filterKeys) == 0 {
		// Print the entire JSON object
		jsonData, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	} else {
		// Print only the filtered keys
		for _, key := range filterKeys {
			if value, exists := obj[key]; exists {
				fmt.Printf("%s: %v\n", key, value)
			} else {
				fmt.Printf("Key '%s' not found in JSON object\n", key)
			}
		}
	}
}

func readJSON(decoder *json.Decoder) error {

	delim, err := checkJSON(decoder)
	if err != nil {
		return err
	}

	switch delim {
	case '{':
		// single json
		err = parseJSON(decoder)
		checkErr(err, "parsing JSON object")

	case '[':
		// json array
		err = parseArray(decoder)
		checkErr(err, "parsing JSON array")
	}

	return nil
}

func parseJSON(decoder *json.Decoder) error {
	var obj map[string]interface{}
	if err := decoder.Decode(&obj); err != nil {
		return err
	}
	printJSON(obj, filterkeys)

	return nil
}

func parseArray(decoder *json.Decoder) error {
	for decoder.More() {
		var obj map[string]interface{}
		if err := decoder.Decode(&obj); err != nil {
			return err
		}
		printJSON(obj, filterkeys)
	}

	return nil
}

func checkJSON(decoder *json.Decoder) (json.Delim, error) {

	// check if the first token is a JSON object or array
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
		checkErr(err, "opening file")
		defer file.Close()
		// check if file is empty

		decoder := json.NewDecoder(file)

		err = readJSON(decoder)
		checkErr(err, "parsing JSON")
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	readCmd.Flags().StringVarP(&filepath, "file", "f", "", "Path to JSON file")
	readCmd.Flags().StringSliceVarP(&filterkeys, "keys", "k", []string{}, "Keys to filter from JSON")

	readCmd.MarkFlagRequired("file")
}
