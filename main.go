package main

import (
	"flag"
	"log"

	"github.com/iancoleman/strcase"
)

type Options struct {
	DoJulia bool
	DoToml  bool
	Verbose bool

	JuliaFilename string
	TomlFilename  string
}

func main() {
	var options Options
	var preview string
	var runAPI bool
	flag.StringVar(&preview, "preview", "", "Display preview of dataset")
	flag.BoolVar(&runAPI, "serve", false, "Display preview of dataset")
	flag.BoolVar(&options.DoJulia, "julia", false, "Export to Julia")
	flag.BoolVar(&options.DoToml, "toml", false, "Export to TOML")
	flag.BoolVar(&options.Verbose, "verbose", false, "Verbose output")
	flag.StringVar(&options.JuliaFilename, "julia-filename", "generated/datasets.jl", "Julia output filename")
	flag.StringVar(&options.TomlFilename, "toml-filename", "generated/metadata.toml", "TOML output filename")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("Usage: parse [options] <file> [datasets...]")
	}

	filename := flag.Args()[0]

	datasets, report, err := parse(filename, options.Verbose)
	if err != nil {
		log.Fatal(err)
	}

	selected := flag.Args()[1:]
	if len(selected) > 0 {
		datasets = filter(datasets, selected)
	}

	err = output(options, datasets, *report)
	if err != nil {
		log.Fatal(err)
	}

	if runAPI {
		serve(datasets)
	}

	log.Println("Done.")
}

func filter(datasets []Dataset, selected []string) []Dataset {
	dsMap := make(map[string]Dataset)
	for _, ds := range datasets {
		dsMap[strcase.ToCamel(ds.Name)] = ds
	}
	return filterMap(dsMap, selected)
}

func filterMap(dsMap map[string]Dataset, selected []string) []Dataset {
	var filtered []Dataset
	for _, name := range selected {
		ds, ok := dsMap[name]
		if !ok {
			log.Fatalf("Dataset %s not found", name)
		}
		filtered = append(filtered, ds)
	}
	return filtered
}
