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
	retirementInputURL  = "https://iso639-3.sil.org/sites/iso639-3/files/downloads/iso-639-3_Retirements.tab"
	retirementTimeout   = 60 * time.Second
	retirementSeparator = '\t'

	retirementFilePrefix = `package iso639_3

// RetiredCodes lookup table. Keys are retired ISO 639-3 codes
var RetiredCodes = map[string]RetiredCode{
`

	retirementFileSuffix = `}
`
)

func main() {
	inputFile := flag.String("i", retirementInputURL,
		fmt.Sprintf("Path or URL to input file in tab-separated iso639-3.sil.org format (default %s)", retirementInputURL))
	outfile := flag.String("o", "", "Output file (default - standard output)")
	flag.Parse()

	rd := getRetirementInput(*inputFile)
	tsvReader := csv.NewReader(rd)
	tsvReader.Comma = retirementSeparator

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

	outputRetirementLookup(wr, langInput)
}

func getRetirementInput(uri string) io.Reader {
	parsedUrl, err := url.Parse(uri)
	if err != nil || parsedUrl.Scheme == "" {
		f, err := os.Open(uri)
		if err != nil {
			log.Fatalf("Can't open input file '%s': %v", uri, err)
		}
		return bufio.NewReader(f)
	}

	httpClient := &http.Client{
		Timeout: retirementTimeout,
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

func outputRetirementStruct(w io.Writer, record []string) error {
	if len(record) != 6 {
		log.Fatalf("outputRetirementStruct got malformed record: %v", record)
	}

	_, err := fmt.Fprintf(w, `"%s": {Id: "%s", RefName: "%s", RetReason: "%s", ChangeTo: "%s", RetRemedy: "%s", Effective: "%s"},
`, record[0], record[0], record[1], record[2], record[3], record[4], record[5])
	return err
}

func outputRetirementLookup(w io.Writer, records [][]string) {
	buf := bytes.Buffer{}

	_, err := fmt.Fprint(&buf, retirementFilePrefix)
	if err != nil {
		log.Fatalf("Error generating: %v", err)
	}

	for _, record := range records {
		err = outputRetirementStruct(&buf, record)
		if err != nil {
			log.Fatalf("Error generating: %v", err)
		}
	}

	_, err = fmt.Fprint(&buf, retirementFileSuffix)
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
