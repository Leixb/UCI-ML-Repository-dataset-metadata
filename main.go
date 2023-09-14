package main

import (
	"encoding/json"
	"fmt"

	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./main <filename>")
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(file)

	// read open bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	attr_types := make(map[string]int)
	attr_roles := make(map[string]int)
	tasks := make(map[string]int)

	missing_values := 0

	out, err := os.Create("metadata.toml")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	out_jl, err := os.Create("datasets.jl")
	if err != nil {
		log.Fatal(err)
	}
	defer out_jl.Close()

	fmt.Fprint(out, "#", "!!! This file was auto-generated.\n\n")
	fmt.Fprint(out_jl, "#", "!!! This file was auto-generated.\n\n")

	for dec.More() {
		var d Dataset
		if err := dec.Decode(&d); err != nil {
			log.Fatal(err)
		}
		println(d.ID, d.Name, d.DateDonated.Year, d.YearCreated)

		tasks[d.Task] += 1

		for _, attr := range d.Attributes {
			attr_types[attr.Type] += 1
			attr_roles[attr.Role] += 1
		}

		if has_missing_values(d.Attributes) {
			missing_values += 1
		}

		err = d.Toml(out)
		if err != nil {
			log.Fatal(err)
		}

		err = d.Julia(out_jl)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Types:")
	for k, v := range attr_types {
		log.Println(k, v)
	}

	log.Println("Roles:")
	for k, v := range attr_roles {
		log.Println(k, v)
	}
	log.Println("Tasks:")
	for k, v := range tasks {
		log.Println(k, v)
	}

	log.Println("Missing values:", missing_values)

}

func has_missing_values(attributes []Attribute) bool {
	for _, attr := range attributes {
		if attr.MissingValues {
			return true
		}
	}
	return false
}
