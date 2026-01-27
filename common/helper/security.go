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
	// sqlKeywordRe     = regexp.MustCompile(`(?i)\b(select|union|insert|update|delete|drop|sleep|benchmark|or\s+1=1)\b`)
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

	jsExecRe      = regexp.MustCompile(`(?i)\b(alert|eval|prompt|confirm|settimeout|setinterval|function)\s*\(`)
	jsPrototypeRe = regexp.MustCompile(`(?i)\b(object|array|string|number|regexp)\.prototype\b`)
	domSinkRe     = regexp.MustCompile(`(?i)\b(location|document|window)\.(hash|href|cookie|write)\b`)
	// sqlTimeRe       = regexp.MustCompile(`(?i)\b(waitfor\s+delay|sleep\s*\(|benchmark\s*\()\b`)
	encodedJsCallRe = regexp.MustCompile(`(?i)(alert|eval|prompt|confirm)[^a-z0-9]*\(`)
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
		// match <tag or &lt;tag (case-insensitive)
		parts = append(parts, `(?i)<\s*`+regexp.QuoteMeta(tag)+`(\b|[^a-z0-9])`)
		parts = append(parts, `(?i)&lt;\s*`+regexp.QuoteMeta(tag)+`(\b|[^a-z0-9])`)
	}
	compiledTagRegex = regexp.MustCompile(strings.Join(parts, "|"))

	return func(value interface{}) error {
		s, _ := value.(string)
		if s == "" {
			return nil
		}

		// 1) pre-normalize decode
		unescaped := deepDecode(s)
		lower := strings.ToLower(unescaped)

		// 2) event attributes
		if eventAttrPattern.MatchString(unescaped) {
			return errors.New("contains event handler attribute (on...=)")
		}
		if jsExecRe.MatchString(lower) {
			return errors.New("javascript execution detected")
		}
		if jsPrototypeRe.MatchString(lower) {
			return errors.New("javascript prototype manipulation detected")
		}
		if domSinkRe.MatchString(lower) {
			return errors.New("dom sink detected")
		}
		// if sqlTimeRe.MatchString(lower) {
		// 	return errors.New("sql time-based injection detected")
		// }
		if encodedJsCallRe.MatchString(unescaped) {
			return errors.New("encoded javascript execution detected")
		}

		// 3) dangerous protocols
		// for _, p := range protoList {
		// 	if strings.Contains(lower, p) {
		// 		return errors.New("contains dangerous protocol: " + p)
		// 	}
		// }
		if dangerousProtoRe.MatchString(lower) {
			return errors.New("dangerous protocol detected")
		}
		// if sqlKeywordRe.MatchString(lower) {
		// 	return errors.New("sql keyword detected")
		// }
		if lfiRe.MatchString(lower) {
			return errors.New("lfi pattern detected")
		}
		if !asciiAllowRe.MatchString(unescaped) {
			return errors.New("non-ascii or disallowed character")
		}

		// 4) css constructs
		for _, cp := range cssPatterns {
			if matched, _ := regexp.MatchString(cp, unescaped); matched {
				return errors.New("contains dangerous CSS construct")
			}
		}

		// 5) special patterns
		for _, sp := range specialPatterns {
			if matched, _ := regexp.MatchString(sp, unescaped); matched {
				return errors.New("contains suspicious token or encoding")
			}
		}

		// 6) fallback generic tag-like (last resort)
		stripped := allowedTagsRe.ReplaceAllString(unescaped, "")
		if compiledTagRegex.MatchString(stripped) {
			return errors.New("contains disallowed HTML tag")
		}

		// 7) UTF-8 validity check
		if !utf8.ValidString(s) {
			return errors.New("contains invalid UTF-8")
		}

		// 8) Latin-only + safe punctuation
		if !latinSafeRe.MatchString(s) {
			return errors.New("contains non-latin or disallowed characters")
		}

		return nil
	}
}
