package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	macrolanguageInputURL  = "https://iso639-3.sil.org/sites/iso639-3/files/downloads/iso-639-3-macrolanguages.tab"
	macrolanguageTimeout   = 60 * time.Second
	macrolanguageSeparator = '\t'

	macrolanguageFilePrefix = `package iso639_3

// MacrolanguageMappings lookup table. Keys are individual language ISO 639-3 codes
var MacrolanguageMappings = map[string]MacrolanguageMapping{
`

	macrolanguageFileSuffix = `}
`
)

func main() {
	inputFile := flag.String("i", macrolanguageInputURL,
		fmt.Sprintf("Path or URL to input file in tab-separated iso639-3.sil.org format (default %s)", macrolanguageInputURL))
	outfile := flag.String("o", "", "Output file (default - standard output)")
	flag.Parse()

	rd := getMacrolanguageInput(*inputFile)
	tsvReader := csv.NewReader(rd)
	tsvReader.Comma = macrolanguageSeparator

	langInput, err := tsvReader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading input file '%s': %v", *inputFile, err)
	}

	langInput = langInput[1:] // skip header

	wr := os.Stdout
	if *outfile != "" {
		var err error
		wr, err = os.Create(*outfile)
		if err != nil {
			log.Fatalf("Can't create output file '%s': %v", *outfile, err)
		}
	}

	outputMacrolanguageLookup(wr, langInput)
}

func getMacrolanguageInput(uri string) io.Reader {
	parsedUrl, err := url.Parse(uri)
	if err != nil || parsedUrl.Scheme == "" {
		f, err := os.Open(uri)
		if err != nil {
			log.Fatalf("Can't open input file '%s': %v", uri, err)
		}
		return bufio.NewReader(f)
	}

	httpClient := &http.Client{
		Timeout: macrolanguageTimeout,
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalf("Can't create request for '%s': %v", uri, err)
	}

	req.Header.Set("User-Agent", "USER")

	r, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Can't download input file '%s': %v", uri, err)
	}
	defer r.Body.Close()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading response from '%s': %v", uri, err)
	}

	return bytes.NewReader(bs)
}

func outputMacrolanguageStruct(w io.Writer, record []string) error {
	if len(record) != 3 {
		log.Fatalf("outputMacrolanguageStruct got malformed record: %v", record)
	}

	_, err := fmt.Fprintf(w, `"%s": {MacrolanguageId: "%s", IndividualId: "%s", Status: "%s"},
`, record[1], record[0], record[1], record[2])
	return err
}

func outputMacrolanguageLookup(w io.Writer, records [][]string) {
	buf := bytes.Buffer{}

	_, err := fmt.Fprint(&buf, macrolanguageFilePrefix)
	if err != nil {
		log.Fatalf("Error generating: %v", err)
	}

	for _, record := range records {
		err = outputMacrolanguageStruct(&buf, record)
		if err != nil {
			log.Fatalf("Error generating: %v", err)
		}
	}

	_, err = fmt.Fprint(&buf, macrolanguageFileSuffix)
	if err != nil {
		log.Fatalf("Error generating: %v", err)
	}

	outBytes, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Error formatting generated code: %v", err)
	}

	_, err = w.Write(outBytes)
	if err != nil {
		log.Fatalf("Error writing to output: %v", err)
	}
}
