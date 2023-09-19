package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
)

const BASE_URL = "https://archive.ics.uci.edu"

type Attribute struct {
	ID            int    `json:"ID"`
	DatasetID     int    `json:"datasetID"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	Type          string `json:"type"`
	Demographic   string `json:"demographic"`
	Description   string `json:"description"`
	Units         string `json:"units"`
	MissingValues bool   `json:"missingValues"`
}

type DescriptiveQuestions struct {
	DatasetID                int    `json:"datasetID"`
	Purpose                  string `json:"purpose"`
	Funding                  string `json:"funding"`
	Represent                string `json:"represent"`
	DataSplits               string `json:"dataSplits"`
	SensitiveInfo            string `json:"sensitiveInfo"`
	PreprocessingDescription string `json:"preprocessingDescription"`
	SoftwareAvailable        string `json:"softwareAvailable"`
	Used                     string `json:"used"`
	OtherInfo                string `json:"otherInfo"`
	DatasetCitation          string `json:"datasetCitation"`
}

type DatasetCreator struct {
	DatasetCreatorsID int      `json:"datasetCreatorsID"`
	DatasetID         int      `json:"datasetID"`
	CreatorID         int      `json:"creatorID"`
	Creators          Creators `json:"creators"`
}

type Creators struct {
	ID          int    `json:"ID"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Institution string `json:"institution"`
	Address     string `json:"address"`
}

type DatasetKeywords struct {
	DatasetKeywordsID int      `json:"datasetKeywordsID"`
	DatasetID         int      `json:"datasetID"`
	KeywordsID        int      `json:"keywordsID"`
	Keywords          Keywords `json:"keywords"`
}

type Keywords struct {
	ID      int    `json:"ID"`
	Keyword string `json:"keyword"`
	Status  string `json:"status"`
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	User      string `json:"user"`
}

type Dataset struct {
	ID           int `json:"ID"`
	UserID       int `json:"userID"`
	IntroPaperID int `json:"introPaperID"`

	Name        string    `json:"Name"`
	Abstract    string    `json:"Abstract"`
	Area        string    `json:"Area"`
	Task        string    `json:"Task"`
	Types       string    `json:"Types"`
	DOI         string    `json:"DOI"`
	DateDonated time.Time `json:"DateDonated"`
	YearCreated int       `json:"YearCreated"`

	IsTabular bool   `json:"isTabular"`
	URLFolder string `json:"URLFolder"`
	URLLink   string `json:"URLLink"`
	Graphics  string `json:"Graphics"`

	IsAvailablePython bool        `json:"isAvailablePython"`
	Status            string      `json:"Status"`
	NumHits           int         `json:"NumHits"`
	NumDownloads      int         `json:"NumDownloads"`
	NumInstances      int         `json:"NumInstances"`
	NumAttributes     int         `json:"NumAttributes"`
	AttributeTypes    string      `json:"AttributeTypes"`
	Slug              string      `json:"Slug"`
	Attributes        []Attribute `json:"Attributes"`

	DatasetCreators []DatasetCreator  `json:"dataset_creators"`
	DatasetKeywords []DatasetKeywords `json:"dataset_keywords"`

	User User `json:"user"`
}

func (dataset *Dataset) ZipURL() string {
	if dataset.URLLink != "" {
		return dataset.URLLink
	}

	link, _ := url.Parse(BASE_URL)
	return link.JoinPath("static", "public", fmt.Sprint(dataset.ID), fmt.Sprintf("%s.zip", dataset.Slug)).String()
}

func (dataset *Dataset) HomeURL() string {
	link, _ := url.Parse(BASE_URL)
	return link.JoinPath("dataset", fmt.Sprint(dataset.ID), dataset.Slug).String()
}

func (dataset *Dataset) NormSlug() string {
	return NormalizeName(dataset.Slug)
}

func (d *Dataset) Target() string {
	for _, attr := range d.Attributes {
		if attr.Role == "Target" {
			return attr.Name
		}
	}
	// last attribute name
	idx := len(d.Attributes) - 1
	if idx < 0 {
		return ""
	}
	last := d.Attributes[len(d.Attributes)-1].Name
	if last == "" || last == " " {
		return fmt.Sprintf("Column%d", idx+1)
	}
	return last
}

func (d *Dataset) TaskName() string {
	if strings.Contains(d.Task, "Regression") {
		return "Regression"
	}
	if strings.Contains(d.Task, "Classification") {
		return "Classification"
	}
	return "Other"
}

func (d *Dataset) AttrsAreAllOK() bool {
	if len(d.Attributes) == 0 {
		return false
	}
	for _, attr := range d.Attributes {
		if attr.Name == "" || attr.Type == "" || attr.Name == " " || attr.Type == " " {
			return false
		}
	}
	return true
}

func (d *Dataset) NeedsCoercion() bool {
	for _, attr := range d.Attributes {
		if attr.Type != "Continuous" {
			return true
		}
	}
	return false
}

func (d *Dataset) CoerceAttrs() []Attribute {
	result := make([]Attribute, 0, len(d.Attributes))
	for _, attr := range d.Attributes {
		if attr.Type != "Continuous" {
			result = append(result, attr)
		}
	}
	return result
}

//go:embed templates/*
var embededFS embed.FS

var metadataTemplate = template.Must(template.New("metadata.toml.tmpl").ParseFS(embededFS, "templates/metadata.toml.tmpl"))

func (d *Dataset) Toml(out io.Writer) error {
	return metadataTemplate.Execute(out, d)
}

func toSciType(t string) string {
	switch t {
	case "Binary":
		return "Finite{2}"
	case "Integer":
		return "Count"
	case "Categorical":
		return "Multiclass"
	case "Continuous":
		return "Continuous"
	default:
		log.Println("Unknown type:", t)
		return t
	}
}

var juliaTemplate = template.Must(template.New("datasets.jl.tmpl").Funcs(
	template.FuncMap{
		"toSciType": toSciType,
		"toCamel":   strcase.ToCamel,
	},
).ParseFS(embededFS, "templates/datasets.jl.tmpl"))

func (d *Dataset) Julia(out io.Writer) error {
	return juliaTemplate.Execute(out, d)
}

var markdownTemplate = template.Must(template.New("markdown.md.tmpl").Funcs(
	template.FuncMap{
		"toSciType": toSciType,
		"toCamel":   strcase.ToCamel,
	},
).ParseFS(embededFS, "templates/markdown.md.tmpl"))

func (d *Dataset) Export(out io.Writer, format string) error {
	switch format {
	case "julia":
		return d.Julia(out)
	case "toml":
		return d.Toml(out)
	case "markdown":
		return d.Markdown(out)
	default:
		return fmt.Errorf("Unknown format: %s", format)
	}
}

func (d *Dataset) Markdown(out io.Writer) error {
	return markdownTemplate.Execute(out, d)
}

func (d *Dataset) HasIDAttribute() bool {
	for _, attr := range d.Attributes {
		if attr.Role == "ID" {
			return true
		}
	}
	return false
}

func NormalizeName(name string) string {
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

func (d *Dataset) HasMissingAttributes() bool {
	for _, attr := range d.Attributes {
		if attr.MissingValues {
			return true
		}
	}
	return false
}

func (d *Dataset) String(method string) (string, error) {
	strWriter := new(strings.Builder)
	var function func(io.Writer) error
	if method == "julia" {
		function = d.Julia
	} else if method == "toml" {
		function = d.Toml
	} else {
		function = d.Markdown
	}
	err := function(strWriter)
	if err != nil {
		return "", err
	}
	return strWriter.String(), nil
}
