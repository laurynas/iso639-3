# ISO 639-3

[![Go Reference](https://pkg.go.dev/badge/github.com/laurynas/iso639-3?status.svg)](https://pkg.go.dev/github.com/laurynas/iso639-3)
[![Test](https://github.com/laurynas/iso639-3/actions/workflows/test.yml/badge.svg)](https://github.com/laurynas/iso639-3/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/laurynas/iso639-3)](https://goreportcard.com/report/github.com/laurynas/iso639-3)

A database of ISO 639-3, ISO 639-2 and ISO 639-1 languages.

## Motivation

There's an excellent [Go library for ISO 639-1](https://github.com/emvi/iso-639-1), but it lacks ISO 639-2 and ISO 639-3 codes.

## Data source

Database is generated (see `cmd/generator.go`) from official ISO 639-3 data. See [official site of the ISO 639-3 Registration Authority](https://iso639-3.sil.org) for details.

The retired language codes database is generated (see `cmd/retirement-generator/main.go`) from the official ISO 639-3 Retirements data.

## Installation

```
go get github.com/laurynas/iso639-3
```

## Examples

```go
iso639_3.LanguagesPart3 // returns ISO 639-3 languages lookup table
iso639_3.LanguagesPart2 // returns ISO 639-2 languages lookup table
iso639_3.LanguagesPart1 // returns ISO 639-1 languages lookup table
iso639_3.RetiredCodes // returns retired ISO 639-3 codes lookup table

iso639_3.FromAnyCode("eng") // returns object representing English language looking through ISO 639-3, ISO 639-2 and ISO 639-1 codes
iso639_3.FromPart3Code("deu") // returns object representing German language looking by ISO 639-3 code
iso639_3.FromPart2Code("ger") // returns object representing German language looking by ISO 639-2 code
iso639_3.FromPart1Code("de") // returns object representing German language looking by ISO 639-1 code
iso639_3.FromName("English") // returns object representing English language looking by language name
iso639_3.GetRetired("sdi") // returns information about a retired code if it exists
```

# Update databases

Rebuild go db files from the latest online database.
> go generate


## Contribute

Feel free to open issues and send pull requests.

## License

MIT