package iso639_3

import (
	"testing"
)

func TestFromPart3Code(t *testing.T) {
	tests := []struct {
		code         string
		expectedName string
	}{
		{"rus", "Russian"},
		{"deu", "German"},
		{"123", ""}, // doesn't exist
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			actual := FromPart3Code(tt.code)

			if tt.expectedName == "" {
				if actual != nil {
					t.Errorf("FromCode() = %v, expected nil", actual)
				}
			} else if actual == nil || actual.Name != tt.expectedName {
				t.Errorf("FromCode() = %v, expected Language with english name %v", actual, tt.expectedName)
			}
		})
	}
}

func TestFromPart2Code(t *testing.T) {
	tests := []struct {
		code         string
		expectedName string
	}{
		{"rus", "Russian"},
		{"ger", "German"},
		{"123", ""}, // doesn't exist
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			actual := FromPart2Code(tt.code)

			if tt.expectedName == "" {
				if actual != nil {
					t.Errorf("FromCode() = %v, expected nil", actual)
				}
			} else if actual == nil || actual.Name != tt.expectedName {
				t.Errorf("FromCode() = %v, expected Language with english name %v", actual, tt.expectedName)
			}
		})
	}
}

func TestFromPart1Code(t *testing.T) {
	tests := []struct {
		code         string
		expectedName string
	}{
		{"ru", "Russian"},
		{"de", "German"},
		{"12", ""}, // doesn't exist
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			actual := FromPart1Code(tt.code)

			if tt.expectedName == "" {
				if actual != nil {
					t.Errorf("FromCode() = %v, expected nil", actual)
				}
			} else if actual == nil || actual.Name != tt.expectedName {
				t.Errorf("FromCode() = %v, expected Language with english name %v", actual, tt.expectedName)
			}
		})
	}
}

func TestFromAnyCode(t *testing.T) {
	tests := []struct {
		code         string
		expectedName string
	}{
		{"rus", "Russian"},
		{"ru", "Russian"},
		{"de", "German"},
		{"ger", "German"},
		{"bgh", "Bugan"}, // retired code that was changed to "bbh"
		{"123", ""},      // doesn't exist
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			actual := FromAnyCode(tt.code)

			if tt.expectedName == "" {
				if actual != nil {
					t.Errorf("FromCode() = %v, expected nil", actual)
				}
			} else if actual == nil || actual.Name != tt.expectedName {
				t.Errorf("FromCode() = %v, expected Language with english name %v", actual, tt.expectedName)
			}
		})
	}
}

func TestFromName(t *testing.T) {
	tests := []struct {
		name          string
		expectedPart3 string
	}{
		{"Russian", "rus"},
		{"German", "deu"},
		{"Elvish", ""}, // doesn't exist (ouch)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FromName(tt.name)

			if tt.expectedPart3 == "" {
				if actual != nil {
					t.Errorf("FromCode() = %v, expected nil", actual)
				}
			} else if actual == nil || actual.Part3 != tt.expectedPart3 {
				t.Errorf("FromCode() = %v, expected Language with Alpha3 %v", actual, tt.expectedPart3)
			}
		})
	}
}

func TestIsRetired(t *testing.T) {
	tests := []struct {
		code            string
		expectedRetired bool
		expectedReason  string
	}{
		{"fri", true, "C"}, // Western Frisian - Changed to Frysk
		{"eng", false, ""}, // English - Not retired
		{"abc", false, ""}, // Non-existent code
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			actual := IsRetired(tt.code)

			if tt.expectedRetired {
				if actual == nil {
					t.Errorf("IsRetired() = nil, expected RetiredCode")
				} else if actual.RetReason != tt.expectedReason {
					t.Errorf("IsRetired() = %v, expected RetiredCode with reason %v", actual, tt.expectedReason)
				}
			} else if actual != nil {
				t.Errorf("IsRetired() = %v, expected nil", actual)
			}
		})
	}
}

func TestFromRetiredCode(t *testing.T) {
	tests := []struct {
		code         string
		expectedName string
	}{
		{"bgh", "Bugan"}, // Single retirement
		{"abc", ""},      // Not retired
		{"eng", ""},      // Not retired
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			actual := FromRetiredCode(tt.code)

			if tt.expectedName == "" {
				if actual != nil {
					t.Errorf("FromRetiredCode() = %v, expected nil", actual)
				}
			} else if actual == nil || actual.Name != tt.expectedName {
				t.Errorf("FromRetiredCode() = %v, expected Language with english name %v", actual, tt.expectedName)
			}
		})
	}
}

func BenchmarkFromAnyCode(b *testing.B) {
	benchmarks := []struct {
		name string
		code string
	}{
		{"ISO639-3", "eng"},
		{"ISO639-2", "ger"},
		{"ISO639-1", "en"},
		{"Retired", "bgh"},
		{"NonExistent", "xyz"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FromAnyCode(bm.code)
			}
		})
	}
}
