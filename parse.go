package main

import (
	"encoding/json"
	"log"
	"os"
)

type Report struct {
	AttrTypes map[string]int
	AttrRoles map[string]int
	Tasks     map[string]int

	WithMissingValues int
}

const APPROX_NUM_DATASETS = 700

func parse(filename string, verbose bool) (datasets []Dataset, report *Report, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	dec := json.NewDecoder(file)

	_, err = dec.Token()
	if err != nil {
		return nil, nil, err
	}

	report = &Report{
		AttrTypes: make(map[string]int),
		AttrRoles: make(map[string]int),
		Tasks:     make(map[string]int),
	}
	datasets = make([]Dataset, 0, APPROX_NUM_DATASETS)

	for dec.More() {
		var d Dataset
		if err := dec.Decode(&d); err != nil {
			return nil, nil, err
		}
		datasets = append(datasets, d)

		if verbose {
			log.Println(d.ID, d.Name, d.YearCreated)
		}

		report.Tasks[d.Task] += 1

		for _, attr := range d.Attributes {
			report.AttrTypes[attr.Type] += 1
			report.AttrRoles[attr.Role] += 1
		}

		if d.HasMissingAttributes() {
			report.WithMissingValues += 1
		}
	}

	return datasets, report, nil
}
