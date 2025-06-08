package iso639_3

// LanguageScope represents language scope as defined in ISO 639-3
type LanguageScope rune

// LanguageType represents language scope as defined in ISO 639-3
type LanguageType rune

const (
	LanguageTypeIndividual    LanguageScope = 'I'
	LanguageTypeSpecial       LanguageScope = 'S'
	LanguageTypeMacrolanguage LanguageScope = 'M'

	LanguageScopeLiving      LanguageType = 'L'
	LanguageScopeHistorical  LanguageType = 'H'
	LanguageScopeAncient     LanguageType = 'A'
	LanguageScopeExtinct     LanguageType = 'E'
	LanguageScopeConstructed LanguageType = 'C'
	LanguageScopeSpecial     LanguageType = 'S'
)

// Language holds language info - all ISO 639 codes along with name and some additional info
type Language struct {
	Part3        string // ISO639-3 code
	Part2B       string // ISO639-2 bibliographic code
	Part2T       string // ISO639-2 terminology code
	Part1        string // ISO639-1 code
	Scope        LanguageScope
	LanguageType LanguageType
	Name         string
	Comment      string
}

// RetiredCode represents a retired ISO 639-3 language code
type RetiredCode struct {
	Id        string // The retired language identifier
	RefName   string // Reference name of the retired language
	RetReason string // Reason for retirement (C=Change, D=Duplicate, N=Nonexistent, S=Split, M=Merge)
	ChangeTo  string // The new language identifier to use instead
	RetRemedy string // Instructions for updating an implementation
	Effective string // The date the retirement became effective
}

//go:generate go run cmd/generator.go -o lang-db.go
//go:generate go run cmd/retirement-generator/main.go -o retired-db.go

// FromPart3Code looks up language for given ISO639-3 three-symbol code.
// Returns nil if not found
func FromPart3Code(code string) *Language {
	if l, ok := LanguagesPart3[code]; ok {
		return &l
	}
	return nil
}

// FromPart2Code looks up language for given ISO639-2 (both bibliographic or terminology) three-symbol code.
// Returns nil if not found
func FromPart2Code(code string) *Language {
	if l, ok := LanguagesPart2[code]; ok {
		return &l
	}
	return nil
}

// FromPart1Code looks up language for given ISO639-1 two-symbol code.
// Returns nil if not found
func FromPart1Code(code string) *Language {
	if l, ok := LanguagesPart1[code]; ok {
		return &l
	}
	return nil
}

// FromRetiredCode looks up language for a retired code by following the chain of retirements.
// If the code was retired multiple times, it will follow the chain until it finds a non-retired code.
// Returns nil if the code is not retired or if the final code is not found.
func FromRetiredCode(code string) *Language {
	if retired := IsRetired(code); retired != nil {
		if retired.ChangeTo == "" {
			return nil
		}
		// Recursively check if the new code is also retired
		if newLang := FromRetiredCode(retired.ChangeTo); newLang != nil {
			return newLang
		}
		// If the new code is not retired, return it
		return FromPart3Code(retired.ChangeTo)
	}
	return nil
}

// FromAnyCode looks up language for given code.
// For three-symbol codes it tries ISO639-3 first, then ISO639-2.
// For two-symbol codes it tries ISO639-1.
// If no match is found and the code is retired, returns the language it was changed to.
// Returns nil if not found
func FromAnyCode(code string) *Language {
	codeLen := len(code)

	if codeLen == 3 {
		// Try ISO639-3 first
		if lang := FromPart3Code(code); lang != nil {
			return lang
		}

		// Then try ISO639-2
		if lang := FromPart2Code(code); lang != nil {
			return lang
		}

		// Finally check if it's a retired code
		return FromRetiredCode(code)
	}

	if codeLen == 2 {
		return FromPart1Code(code)
	}

	return nil
}

// FromName looks up language for given reference name.
// Returns nil if not found
func FromName(name string) *Language {
	for _, l := range LanguagesPart3 {
		if l.Name == name {
			return &l
		}
	}
	return nil
}

// IsRetired checks if a given ISO639-3 code is retired.
// Returns the RetiredCode if found, nil otherwise
func IsRetired(code string) *RetiredCode {
	if r, ok := RetiredCodes[code]; ok {
		return &r
	}
	return nil
}
