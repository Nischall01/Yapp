package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/text/unicode/norm"
)

const FileSize int64 = 10

var usernameRegex = regexp.MustCompile(`^[a-z0-9_.-]{3,32}$`)

var hallNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.\- ]{3,32}$`)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
var hexRegex = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)

var blockedExt = map[string]struct{}{
	".exe": {}, ".bat": {}, ".cmd": {}, ".msix": {},
	".scr": {}, ".pif": {}, ".dll": {}, ".jse": {}, ".vbs": {},
	".vbe": {}, ".wsf": {}, ".wsh": {}, ".ps1": {}, ".psm1": {}, ".reg": {},
	".jar": {}, ".dmg": {}, ".iso": {}, ".pkg": {}, ".sh": {},
	".virus": {},
}

// NAME SECTION
func SanitizeUsername(s string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if !usernameRegex.MatchString(s) {
		return "", ErrorInvalidUsername
	}
	return s, nil
}

func SanitizeHallname(s string) (string, error) {
	s = strings.TrimSpace(s)

	if !hallNameRegex.MatchString(s) {
		return "", ErrorInvalidHallName
	}
	return s, nil
}

func SanitizeDisplayName(s string) (string, error) {
	s = strings.TrimSpace(s)

	s = norm.NFKC.String(s)

	s = strings.Join(strings.Fields(s), " ")

	if s == "" {
		return "", ErrorInvalidDisplayName
	}

	length := utf8.RuneCountInString(s)
	if length < 3 || length > 32 {
		return "", ErrorInvalidDisplayName
	}

	return s, nil
}

// EMAIL SECTION
func SanitizeEmail(s string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	// very light check; rely on validator.v10 too
	if len(s) < 6 || len(s) > 254 || !strings.Contains(s, "@") || !emailRegex.MatchString(s) {
		return "", ErrorInvalidEmail
	}
	return s, nil
}

func SanitizePhoneE164(ptr *string) (*string, error) {
	if ptr == nil {
		return nil, nil
	}
	s := strings.TrimSpace(*ptr)
	if s == "" {
		return nil, nil
	}
	// If you use libphonenumber:
	// num, err := phonenumbers.Parse(s, "NP") // or your default region
	// if err != nil || !phonenumbers.IsValidNumber(num) { return nil, ErrInvalidPhone }
	// e164 := phonenumbers.Format(num, phonenumbers.E164)
	// return &e164, nil
	// If not using lib yet, minimally keep digits/+ and do a length check:
	s = keepPlusDigits(s)
	if len(s) < 7 || len(s) > 20 {
		return nil, ErrorInvalidPhoneNumber
	}
	return &s, nil
}

func keepPlusDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '+' || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// PASSWORD SECTION
const minEntropyBits = 60.0 // ~good baseline for online attacks; use 70–80 for higher risk

func SanitizePasswordPolicy(raw string) (string, error) {
	// Do NOT silently modify. Reject confusing whitespace at edges.
	if strings.TrimSpace(raw) != raw {
		return "", ErrorPasswordWhiteSpace
	}
	if err := passwordvalidator.Validate(raw, minEntropyBits); err != nil {
		return "", ErrorInvalidPassword
	}
	return raw, nil
}

// COLOR SECTION
func SanitizeColorFormat(colorHex *string) (*string, error) {

	if colorHex == nil {
		return nil, nil
	}

	s := strings.TrimSpace(*colorHex)
	if s == "" {
		return nil, nil
	}

	if !hexRegex.MatchString(s) {
		return nil, ErrorInvalidBannerColor
	}

	return &s, nil
}

// TEXT SECTION
func SanitizeText(text *string) (*string, error) {
	if text == nil {
		return nil, nil
	}

	s := strings.TrimSpace(*text)
	if s == "" {
		return nil, nil
	}

	// xss injection prevention
	p := bluemonday.UGCPolicy()
	s = p.Sanitize(s)

	return &s, nil
}

func SanitizeMessageContent(content *string) *string {

	// s = strings.Join(strings.Fields(s), " ")

	s := strings.TrimSpace(*content)
	s = norm.NFKC.String(s)

	p := bluemonday.UGCPolicy()
	s = p.Sanitize(s)

	return &s

}

// FILE SECTION
func ValidateFileName(fileName string) (string, error) {

	s := strings.TrimSpace(fileName)

	_, err := os.Stat(s)
	if err == nil && os.IsNotExist(err) {
		return s, nil
	}

	return "", ErrorInvalidFileName
}

func ValidateFileType(fileType *string, url string) (*string, error) {

	//		check the url for the filetype
	ext := strings.ToLower(filepath.Ext(url))

	if _, bad := blockedExt[ext]; bad {
		return nil, ErrorBadFileType
	}

	//	cross checking (condition, filetype != nil)
	if fileType != nil {
		if !strings.Contains(strings.ToLower(*fileType), ext) {

			// fileType contains diff file than url
			return nil, ErrorFileUnmatch

		}
	}

	return &ext, nil
}
