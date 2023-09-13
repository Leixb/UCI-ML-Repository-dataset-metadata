package main

import (
	"encoding/json"
	"strings"

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

	missing_values := 0

	out, err := os.Create("metadata.toml")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	for dec.More() {
		var d Dataset
		if err := dec.Decode(&d); err != nil {
			log.Fatal(err)
		}
		println(d.ID, d.Name, d.DateDonated.Year, d.YearCreated)

		for _, attr := range d.Attributes {
			attr_types[attr.Type] += 1
			attr_roles[attr.Role] += 1
		}

		if has_missing_values(d.Attributes) {
			missing_values += 1
		}

		d.print_toml(out)
	}

	log.Println("Types:")
	for k, v := range attr_types {
		log.Println(k, v)
	}

	log.Println("Roles:")
	for k, v := range attr_roles {
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

func normalize_name(name string) string {
	// Remove all non-alphanumeric characters
	result := ""
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			result += "_"
		} else {
			result += string(c)
		}
	}
	return strings.ToLower(result)
}
