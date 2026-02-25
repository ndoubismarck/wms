package stringutil

import (
	cryptoRand "crypto/rand"
	"math"
	mathBig "math/big"
	mathRand "math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func InSlice(needle string, haystack []string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}

func ToInt64(str string) int64 {
	str = strings.TrimSpace(str)
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

func ToInt64Ptr(str string) *int64 {
	val := ToInt64(str)
	return &val
}

func ToUint64(str string) uint64 {
	str = strings.TrimSpace(str)
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return uint64(math.Abs(float64(result)))
}

func ToUint64Ptr(str string) *uint64 {
	val := ToUint64(str)
	return &val
}

func ToFloat64(str string) float64 {
	str = strings.TrimSpace(str)
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return result
}

func ToFloat64Ptr(str string) *float64 {
	val := ToFloat64(str)
	return &val
}

func ToUint32(str string) uint32 {
	str = strings.TrimSpace(str)
	result, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(math.Abs(float64(result)))
}

func ToUint32Ptr(str string) *uint32 {
	val := ToUint32(str)
	return &val
}

func Default(val string, valDefault string) string {
	if val != "" {
		return val
	}
	return valDefault
}

func DefaultRune(val rune, valDefault rune) rune {
	if val != 0 {
		return val
	}
	return valDefault
}

func IsEmpty(val string) bool {
	return len(strings.TrimSpace(val)) == 0
}

func IsHTML(s string) bool {
	regex := regexp.MustCompile(`<("[^"]*"|'[^']*'|[^'">])*>`)
	return regex.MatchString(s)
}

func IsNumber(val string) bool {
	regex := regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	return regex.MatchString(val)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ToSnakeCase(input string) string {
	regex := regexp.MustCompile(`([a-z])([A-Z])`)
	result := regex.ReplaceAllString(input, "$1 $2")
	result = strings.ToLower(result)
	result = strings.TrimSpace(result)
	return strings.ReplaceAll(result, " ", "_")
}

func SplitCamelCase(input string) string {
	regex := regexp.MustCompile(`([a-z])([A-Z])`)
	result := regex.ReplaceAllString(input, "$1 $2")
	result = strings.ToLower(result)
	result = strings.TrimSpace(result)
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}
	return result
}

func ToSlug(s string) string {
	// Normalize Unicode text to decompose accents
	s = norm.NFD.String(s)

	// Remove diacritics (accents) by filtering out non-spacing marks
	var sb strings.Builder
	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		sb.WriteRune(r)
	}
	s = sb.String()

	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with spaces
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, " ")

	// Trim leading and trailing spaces
	s = strings.TrimSpace(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove any leftover multiple hyphens
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")
	return s
}

func GenerateRandom(n int) string {
	ltr := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := cryptoRand.Int(cryptoRand.Reader, mathBig.NewInt(int64(len(ltr))))
		if err != nil {
			var lrn = []rune(ltr)
			b := make([]rune, n)
			for i := range b {
				b[i] = lrn[mathRand.Intn(len(lrn))]
			}
			return string(b)
		}
		ret[i] = ltr[num.Int64()]
	}
	return string(ret)
}

func ToPtr(str string) *string {
	return &str
}
