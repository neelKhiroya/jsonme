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
var outputType string

type keyCount struct {
	Count int    `json:"count"`
	Key   string `json:"key"`
}

func checkErr(err error, errorMessage string) {
	if err != nil {
		fmt.Printf("Error %s: %v\n", errorMessage, err)
		os.Exit(1)
	}
}

func printJSON(obj map[string]interface{}, filterKeys []string) {
	var dataToPrint map[string]interface{}

	if len(filterKeys) == 0 {
		dataToPrint = obj
	} else {
		dataToPrint = make(map[string]interface{})
		for _, key := range filterKeys {
			if value, exists := obj[key]; exists {
				dataToPrint[key] = value
			}
		}

		// skip if none of the keys matched
		if len(dataToPrint) == 0 {
			return
		}
	}

	jsonData, err := json.MarshalIndent(dataToPrint, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
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

func countJSON(obj map[string]interface{}, filterKeys []string) []keyCount {
	var countToReturn []keyCount

	if len(filterKeys) != 0 {
		for _, key := range filterKeys {
			if _, exists := obj[key]; exists {
				tempCount := &keyCount{
					Key:   key,
					Count: 0,
				}
				tempCount.Count += 1
				countToReturn = append(countToReturn, *tempCount)
			}
		}
	}

	return countToReturn
}

func mergeArrays(returnArray []keyCount, mergeArray []keyCount) []keyCount {
	for _, keycount := range mergeArray {
		found := false
		for i, r := range returnArray {
			if r.Key == keycount.Key {
				returnArray[i].Count += 1
				found = true
				break
			}
		}
		if !found {
			returnArray = append(returnArray, keycount)
		}
	}
	return returnArray
}

func printCount(countsToPrint []keyCount) {
	// parse array as json
	data, err := json.MarshalIndent(countsToPrint, "", "  ")
	checkErr(err, "error reading count json")
	fmt.Println(string(data))
}

func parseArray(decoder *json.Decoder) error {
	if outputType == "c" {
		// count the keys passed
		var keycount []keyCount

		// iterate objs
		for decoder.More() {
			var obj map[string]interface{}
			if err := decoder.Decode(&obj); err != nil {
				return err
			}
			tempCount := countJSON(obj, filterkeys)
			keycount = mergeArrays(keycount, tempCount)
		}
		printCount(keycount)
	} else {
		for decoder.More() {
			var obj map[string]interface{}
			if err := decoder.Decode(&obj); err != nil {
				return err
			}
			printJSON(obj, filterkeys)
		}
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

func checkValidOutputType(outputType string) error {
	if outputType == "k" || outputType == "c" {
		return nil
	}

	return fmt.Errorf("invalid output type %s", outputType)
}

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

		// check if output type is valid
		err = checkValidOutputType(outputType)
		checkErr(err, "with output type")

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
	readCmd.Flags().StringVarP(&outputType, "output", "o", "k", "Type of output - \nc (count)\nj (json: default)\n")

	readCmd.MarkFlagRequired("file")
}
