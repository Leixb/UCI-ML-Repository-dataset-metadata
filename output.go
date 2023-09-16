package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func output(opts Options, datasets []Dataset, report Report) error {
	if opts.DoJulia {
		err := export_dataset_list(datasets, "datasets.jl", "julia")
		if err != nil {
			return err
		}
	}
	if opts.DoToml {
		err := export_dataset_list(datasets, "metadata.toml", "toml")
		if err != nil {
			return err
		}
	}

	if opts.Verbose {
		log.Println("Datasets:", len(datasets))
		json, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
	}

	return nil
}

func export_dataset_list(datasets []Dataset, filename string, format string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Fprint(out, "#", "!!! This file was auto-generated.\n\n")
	for _, ds := range datasets {
		ds.Export(out, format)
	}
	return nil
}
