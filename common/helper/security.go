package helper

import (
	"errors"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/text/unicode/norm"
)

// -----------------------------
// Lists (abridged but extensive)
// -----------------------------
var (
	blacklistTags = []string{
		"html", "head", "body", "title", "meta", "base", "link", "style", "script", "noscript", "template",
		"form", "input", "textarea", "select", "option", "button", "datalist",
		"img", "picture", "source", "video", "audio", "track", "canvas",
		"iframe", "frame", "frameset", "object", "embed", "param", "applet",
		"svg", "g", "path", "rect", "circle", "ellipse", "line", "polyline", "polygon", "use", "defs", "symbol", "image", "text", "tspan",
		"math", "mrow", "mi", "mn", "mo", "mtext", "mglyph", "ms", "mtable", "mtr", "mtd", "annotation",
		"iframe", "object", "embed", "isindex", "layer", "ilayer", "noframes", "blink", "xmp", "plaintext",
	}
	// protoList = []string{
	// 	"javascript:", "data:", "vbscript:", "file:", "filesystem:", "blob:",
	// 	"about:", "chrome:", "chrome-extension:", "moz-extension:", "view-source:",
	// }
	dangerousProtoRe = regexp.MustCompile(`(?i)\b(javascript|data|vbscript|file|filesystem|blob|about|chrome|chrome-extension|moz-extension|view-source):`)
	sqlKeywordRe     = regexp.MustCompile(
		`(?i)(
		(['"` + "`" + `]\s*(or|and)\b)|
		((or|and)\s*[\d'"]?\s*=\s*[\d'"])|
		((or|and)\s*(sleep|benchmark|replace)\s*\()|
		(['"` + "`" + `]\s*--)|
		(['"` + "`" + `]\s*#)
	)`,
	)

	lfiRe        = regexp.MustCompile(`(?i)(\.\./|\.\.\\|/etc/passwd|boot.ini|win.ini|.env)`)
	asciiAllowRe = regexp.MustCompile(
		`^[A-Za-z0-9 .,:;'"()\[\]{}+\-*/=<>&!?%#@_~^]*$`,
	)
	cssPatterns = []string{
		`(?i)expression\s*\(`,        // expression(
		`(?i)-moz-binding\s*:`,       // -moz-binding
		`(?i)url\s*\(\s*data:`,       // url(data:
		`(?i)url\s*\(\s*javascript:`, // url(javascript:
		`(?i)@import\s+`,             // @import
	}
	specialPatterns = []string{
		`(?i)<!doctype`,   // doctype
		`(?i)<!--`,        // comment
		`(?i)<!\[CDATA\[`, // cdata
		`\x00`,            // null byte
		`%00`,             // url-encoded null
		`\\u0000`,         // escaped null
		`%3c`,             // %3c == <
		`%3e`,             // %3e == >
		`(?i)utf-7`,       // UTF-7 marker attempts
	}
	eventAttrPattern = regexp.MustCompile(`(?i)\bon[a-z]+\s*=`)
	anyTagRe         = regexp.MustCompile(`(?i)<\s*/?\s*[a-z][a-z0-9]*(?:\s+[^>]+)?>`)
	hexEntityRe      = regexp.MustCompile(`&#x([0-9A-Fa-f]+);?`)
	decEntityRe      = regexp.MustCompile(`&#([0-9]+);?`)
	zeroWidthRe      = regexp.MustCompile(string([]rune{
		'\u200B',
		'\u200C',
		'\u200D',
		'\uFEFF',
		'\u2060',
	}))
	latinSafeRe = regexp.MustCompile(
		`^[A-Za-z0-9 .,;:_\-+*/=()!%&@#?$'"<>/\n\r\t]*$`,
	)
	allowedTagsRe = regexp.MustCompile(
		`(?i)</?(p|b|i|ul|ol|li)\s*>`,
	)

	jsExecRe        = regexp.MustCompile(`(?i)\b(alert|eval|prompt|confirm|settimeout|setinterval|function)\s*\(`)
	jsPrototypeRe   = regexp.MustCompile(`(?i)\b(object|array|string|number|regexp)\.prototype\b`)
	domSinkRe       = regexp.MustCompile(`(?i)\b(location|document|window)\.(hash|href|cookie|write)\b`)
	sqlTimeRe       = regexp.MustCompile(`(?i)\b(waitfor\s+delay|sleep\s*\(|benchmark\s*\()\b`)
	encodedJsCallRe = regexp.MustCompile(`(?i)(alert|eval|prompt|confirm)[^a-z0-9]*\(`)
	xmlDeclRe       = regexp.MustCompile(`(?i)<\?(xml|xsl|php)`)
)

// deprecated
var compiledTagRegex *regexp.Regexp

// -----------------------------
// Decoding helpers
// -----------------------------

// decodeNumericEntities converts both hex (&#xHH;) and decimal (&#DDD;) numeric entities to runes.
func decodeNumericEntities(s string) string {
	s = hexEntityRe.ReplaceAllStringFunc(s, func(m string) string {
		parts := hexEntityRe.FindStringSubmatch(m)
		if len(parts) < 2 {
			return m
		}
		v, err := strconv.ParseUint(parts[1], 16, 32)
		if err != nil {
			return m
		}
		return string(rune(v))
	})

	s = decEntityRe.ReplaceAllStringFunc(s, func(m string) string {
		parts := decEntityRe.FindStringSubmatch(m)
		if len(parts) < 2 {
			return m
		}
		v, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			return m
		}
		return string(rune(v))
	})

	return s
}

// small helpers to parse hex/dec into rune without importing strconv repeatedly
func fmtSscanfHex(hexStr string, out *rune) (int, error) {
	// parse hex
	var v uint64
	var err error
	v, err = parseUint(hexStr, 16)
	if err != nil {
		return 0, err
	}
	*out = rune(v)
	return 1, nil
}
func fmtSscanfDec(decStr string, out *rune) (int, error) {
	var v uint64
	var err error
	v, err = parseUint(decStr, 10)
	if err != nil {
		return 0, err
	}
	*out = rune(v)
	return 1, nil
}
func parseUint(s string, base int) (uint64, error) {
	// keep import light: use strconv
	return strconvParseUint(s, base)
}
func strconvParseUint(s string, base int) (uint64, error) {
	// wrapper for strconv.ParseUint so imports clear
	return strconv.ParseUint(s, base, 64)
}

func deepDecode(input string) string {

	// 1) HTML entities
	s := html.UnescapeString(input)

	// 2) URL encoding
	if u, err := url.QueryUnescape(s); err == nil {
		s = u
	}

	// 3) Numeric entities
	s = decodeNumericEntities(s)

	// 4) Zero-width chars
	s = zeroWidthRe.ReplaceAllString(s, "")

	// 5) Unicode normalization
	s = norm.NFKC.String(s)

	return s
}

// -----------------------------
// Rule: NoXSSFullScanWithDecode
// -----------------------------

// NoXSSFullScanWithDecode returns an ozzo-validation RuleFunc with deep decoding normalization
// and aggressive detection. It returns an error with a short reason.
func NoXSSFullScanWithDecode() validation.RuleFunc {
	var parts []string
	for _, tag := range blacklistTags {
		parts = append(parts,
			`(?i)<\s*/?\s*`+regexp.QuoteMeta(tag)+`(\b|[^a-z0-9])`,
			`(?i)&lt;\s*/?\s*`+regexp.QuoteMeta(tag)+`(\b|[^a-z0-9])`,
		)
	}
	compiledTagRegex = regexp.MustCompile(strings.Join(parts, "|"))

	return func(value interface{}) error {
		s, _ := value.(string)
		if s == "" {
			return nil
		}

		// =========================
		// 1. FULL NORMALIZATION
		// =========================
		decoded := deepDecode(s)
		lower := strings.ToLower(decoded)

		// =========================
		// 2. ACTIVE XSS
		// =========================
		if eventAttrPattern.MatchString(decoded) {
			return errors.New("xss: event handler detected")
		}
		if jsExecRe.MatchString(lower) {
			return errors.New("xss: javascript execution")
		}
		if jsPrototypeRe.MatchString(lower) {
			return errors.New("xss: prototype pollution")
		}
		if domSinkRe.MatchString(lower) {
			return errors.New("xss: dom sink")
		}
		if encodedJsCallRe.MatchString(decoded) {
			return errors.New("xss: encoded js call")
		}
		if xmlDeclRe.MatchString(decoded) {
			return errors.New("xss: xml execution")
		}
		if dangerousProtoRe.MatchString(lower) {
			return errors.New("xss: dangerous protocol")
		}

		// =========================
		// 3. SQL INJECTION (CONTEXTUAL)
		// =========================
		// ' " ` ALLOWED unless forming SQL pattern
		if sqlKeywordRe.MatchString(lower) {
			return errors.New("sqli: keyword pattern")
		}
		if sqlTimeRe.MatchString(lower) {
			return errors.New("sqli: time-based")
		}

		// =========================
		// 4. LFI / PATH TRAVERSAL
		// =========================
		if lfiRe.MatchString(lower) {
			return errors.New("lfi: path traversal")
		}

		// =========================
		// 5. CSS ATTACKS
		// =========================
		for _, cp := range cssPatterns {
			if matched, _ := regexp.MatchString(cp, decoded); matched {
				return errors.New("xss: css injection")
			}
		}

		// =========================
		// 6. SPECIAL ENCODING
		// =========================
		for _, sp := range specialPatterns {
			if matched, _ := regexp.MatchString(sp, decoded); matched {
				return errors.New("attack: suspicious encoding")
			}
		}

		// =========================
		// 7. HTML TAG FALLBACK
		// =========================
		stripped := allowedTagsRe.ReplaceAllString(decoded, "")
		if compiledTagRegex.MatchString(stripped) {
			return errors.New("xss: disallowed html tag")
		}

		// =========================
		// 8. UTF-8 VALIDITY
		// =========================
		if !utf8.ValidString(decoded) {
			return errors.New("encoding: invalid utf-8")
		}

		// =========================
		// 9. SAFE CHAR SET (LAST)
		// =========================
		// Quotes allowed here
		if !latinSafeRe.MatchString(decoded) {
			return errors.New("charset: disallowed characters")
		}

		return nil
	}
}
