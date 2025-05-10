/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"fmt"

	"github.com/spf13/cobra"
)

var exportPath string
// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export filtered JSON data to a new file",
	Long: `Streams and exports filtered data from a JSON file to a new file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Open source file
		file, err := os.Open(filepath)
		checkErr(err, "opening input JSON")
		defer file.Close()

		decoder := json.NewDecoder(file)
		delim, err := checkJSON(decoder)
		checkErr(err, "validating JSON structure")

		// Default to ./out.json if no path provided
		if exportPath == "" {
			exportPath = "out.json"
		} else {
			// Handle if user gave a folder
			info, err := os.Stat(exportPath)
			if err == nil && info.IsDir() {
				exportPath = filepath.Join(exportPath, "exported_output.json")
			}
		}

		// Create output file
		outFile, err := os.Create(exportPath)
		checkErr(err, "creating export file")
		defer outFile.Close()

		encoder := json.NewEncoder(outFile)

		// Export logic
		if delim == '[' {
			outFile.Write([]byte("["))
			first := true
			for decoder.More() {
				var item map[string]interface{}
				err := decoder.Decode(&item)
				checkErr(err, "decoding array item")

				filtered := filterMap(item, filterkeys)
				if !first {
					outFile.Write([]byte(","))
				}
				encoder.Encode(filtered)
				first = false
			}
			outFile.Write([]byte("]"))
		} else {
			var obj map[string]interface{}
			err := decoder.Decode(&obj)
			checkErr(err, "decoding object")
			filtered := filterMap(obj, filterkeys)
			encoder.Encode(filtered)
		}

		fmt.Printf("✅ Exported data saved to: %s\n", exportPath)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&filepath, "file", "f", "", "Path to input JSON file")
	exportCmd.Flags().StringVarP(&exportPath, "output", "o", "", "Optional output file path (defaults to ./exported_output.json)")
	exportCmd.Flags().StringSliceVarP(&filterkeys, "filter", "k", []string{}, "Keys to filter from JSON data")
	exportCmd.MarkFlagRequired("file")
}

// filterMap returns only key:value pairs from `input` where key is in `keys`
func filterMap(input map[string]interface{}, keys []string) map[string]interface{} {
	if len(keys) == 0 {
		return input
	}
	filtered := make(map[string]interface{})
	for _, key := range keys {
		if val, ok := input[key]; ok {
			filtered[key] = val
		}
	}
	return filtered
}