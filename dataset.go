package main

import (
	"fmt"
	"io"
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

func (d *Dataset) Camel() string {
	return strcase.ToCamel(d.Name)
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

var metadataTemplate = template.Must(template.New("metadata.toml").ParseFiles("templates/metadata.toml"))

func (d *Dataset) Toml(out io.Writer) error {
	return metadataTemplate.Execute(out, d)
}

var juliaTemplate = template.Must(template.New("datasets.jl").ParseFiles("templates/datasets.jl"))

func (d *Dataset) Julia(out io.Writer) error {
	return juliaTemplate.Execute(out, d)
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
