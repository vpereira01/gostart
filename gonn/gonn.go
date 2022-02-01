package gonn

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
)

func Train(recordsFileName string) {
	// Load in a dataset, with headers. Header attributes will be stored.
	// Think of instances as a Data Frame structure in R or Pandas.
	// You can also create instances from scratch.
	rawData, err := base.ParseCSVToInstances(recordsFileName, true)
	if err != nil {
		panic(err)
	}

	// Print a pleasant summary of your data.
	fmt.Println(rawData)
}
