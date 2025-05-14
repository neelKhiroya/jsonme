package util

import (
	"fmt"
	"os"
)

func CheckErr(err error, errorMessage string) {
	if err != nil {
		fmt.Printf("Error %s: %v\n", errorMessage, err)
		os.Exit(1)
	}
}
