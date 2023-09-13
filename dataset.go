package main

import (
	"fmt"
	"io"
	"net/url"
	"time"
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

func (dataset *Dataset) zip_url() string {
	if dataset.URLLink != "" {
		return dataset.URLLink
	}

	link, _ := url.Parse(BASE_URL)
	return link.JoinPath("static", "public", fmt.Sprint(dataset.ID), fmt.Sprintf("%s.zip", dataset.Slug)).String()
}

func (dataset *Dataset) home_url() string {
	link, _ := url.Parse(BASE_URL)
	return link.JoinPath("dataset", fmt.Sprint(dataset.ID), dataset.Slug).String()
}

// Print the dataset metadata in TOML format.
func (dataset *Dataset) print_toml(io io.Writer) {
	key := normalize_name(dataset.Slug)
	fmt.Fprintf(io, "[%s]\n", key)
	fmt.Fprintf(io, "url = %q\n", dataset.zip_url())
	fmt.Fprintf(io, "sha256 = \"\"\n")
	fmt.Fprintf(io, "[%s.meta]\n", key)
	fmt.Fprintf(io, "title = %q\n", dataset.Name)
	fmt.Fprintf(io, "description = %q\n", dataset.Abstract)
	fmt.Fprintf(io, "doi = %q\n", dataset.DOI)
	fmt.Fprintf(io, "kind = %q\n", dataset.Types)
	fmt.Fprintf(io, "year = %d\n", dataset.YearCreated)
	fmt.Fprintf(io, "home = %q\n", dataset.home_url())
	fmt.Fprintf(io, "family = %q\n", dataset.Area)
	fmt.Fprintln(io)
}
