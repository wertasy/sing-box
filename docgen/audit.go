package docgen

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	linkPattern       = regexp.MustCompile(`\[([^\]]*)\]\(([^)]*)\)`)
	markPattern       = regexp.MustCompile(`[*]{1,2}|={1,2}`)
	whitespacePattern = regexp.MustCompile(`\s+`)
	tableKeyPattern   = regexp.MustCompile("^`([^`]+)`")
	// Dotted deep headings are nested field keys (e.g. zh "##### platform.http_proxy.enabled").
	nestedHeadingPattern = regexp.MustCompile(`^#{5,} [a-z0-9_]+(\.[a-z0-9_]+)+\s*$`)
)

// normalizeLine reduces a documentation line to its comparable form.
func normalizeLine(line string) string {
	line = linkPattern.ReplaceAllString(line, "$1 ($2)")
	line = strings.ReplaceAll(line, "`", "")
	// CJK and ASCII full stops are equivalent sentence endings.
	line = strings.ReplaceAll(line, "。", ".")
	line = markPattern.ReplaceAllString(line, "")
	line = whitespacePattern.ReplaceAllString(line, " ")
	return strings.TrimSpace(line)
}

// normalizeBlock joins non-empty normalized lines of a block.
func normalizeBlock(lines []string) string {
	var normalized []string
	for _, line := range lines {
		text := normalizeLine(line)
		if text == "" {
			continue
		}
		normalized = append(normalized, text)
	}
	return strings.Join(normalized, "\n")
}

type oldSection struct {
	heading string
	keys    []string
	prose   map[string][]string
	levels  map[string]int // heading level of each key extracted from a heading
	anyLine bool
}

// extractOldSection reads a section of a baseline document. The heading is
// located level-insensitively; within the section, headings deeper than the
// section heading are field keys (plain headings exactly one level below,
// plus dotted nested keys at any deeper level), and headings at the same or a
// shallower level cannot occur inside the section bounds.
func extractOldSection(lines []string, heading string, format string) (*oldSection, error) {
	location, err := locateSection(lines, heading)
	if err != nil {
		return nil, err
	}
	section := &oldSection{
		heading: heading,
		prose:   make(map[string][]string),
		levels:  make(map[string]int),
	}
	addKey := func(key string, level int) {
		if _, loaded := section.prose[key]; !loaded {
			section.keys = append(section.keys, key)
			section.prose[key] = nil
		}
		if level > 0 {
			section.levels[key] = level
		}
	}
	var current string
	for _, line := range lines[location.headingIndex+1 : location.endIndex] {
		if line == "" {
			continue
		}
		section.anyLine = true
		switch {
		case atxHeadingLevel(line) > location.level:
			// Deeper headings are field keys (plain headings one level below,
			// dotted nested keys at any deeper level); other deeper headings
			// are prose (e.g. "##### Request format").
			level := atxHeadingLevel(line)
			if level <= location.level+1 || nestedHeadingPattern.MatchString(line) {
				current = headingText(line)
				addKey(current, level)
			} else if current != "" {
				section.prose[current] = append(section.prose[current], line)
			}
		case format == "table" && strings.HasPrefix(line, "|"):
			row := strings.TrimPrefix(line, "|")
			cells := strings.SplitN(row, "|", 2)
			if len(cells) == 0 {
				continue
			}
			first := strings.TrimSpace(cells[0])
			if first == "" || strings.HasPrefix(first, "-") {
				continue // separator row
			}
			match := tableKeyPattern.FindStringSubmatch(first)
			if match == nil {
				continue
			}
			current = match[1]
			addKey(current, 0)
		case current != "":
			section.prose[current] = append(section.prose[current], line)
		default:
			// Section level prose without keys.
		}
	}
	return section, nil
}

// auditSectionCompare verifies that every baseline key and prose line is
// present in the regenerated section.
// defaultSentencePattern matches prose describing a field default value,
// which is generated from the @default marker instead of the description.
var defaultSentencePattern = regexp.MustCompile(`(?i)(\bis|\bwill be) used by default|used if empty|Defaults to |默认使用|默认为`)

var backtickValuePattern = regexp.MustCompile("`([^`]+)`")

func defaultValueRetained(line string, newBlock string, rawBlock string) bool {
	values := backtickValuePattern.FindAllStringSubmatch(line, -1)
	if len(values) == 0 {
		return true
	}
	normalized := normalizeBlock(strings.Split(newBlock, "\n"))
	for _, value := range values {
		// The generated side is normalized (code spans stripped, whitespace
		// collapsed), so compare the value in the same form. Token boundaries
		// use ASCII alphanumerics only: CJK text has no word separators, and
		// `80` must not match inside `8080`.
		v := normalizeLine(value[1])
		if v == "" {
			// Values like `*` normalize to nothing; compare the raw text.
			if strings.Contains(rawBlock, value[1]) {
				return true
			}
			continue
		}
		pattern, err := regexp.Compile(`(?:^|[^0-9A-Za-z])` + regexp.QuoteMeta(v) + `(?:$|[^0-9A-Za-z])`)
		if err != nil {
			continue
		}
		if pattern.MatchString(normalized) {
			return true
		}
	}
	return false
}

func auditSectionCompare(file string, heading string, old *oldSection, newKeys map[string][]string, lang string, removeKeys []string) []error {
	var issues []error
	if len(old.keys) == 0 && old.anyLine {
		issues = append(issues, fmt.Errorf("%s [%s]: baseline section has content but no extractable keys", file, heading))
		return issues
	}
	removed := make(map[string]bool)
	for _, key := range removeKeys {
		removed[key] = true
	}
	for _, key := range old.keys {
		if removed[key] {
			// Intentionally dropped legacy keys: present at the baseline but
			// no longer a struct field (typically a renamed zh-only key).
			continue
		}
		newProse, found := newKeys[key]
		if !found {
			issues = append(issues, fmt.Errorf("%s [%s] %s: key missing in generated section", file, heading, key))
			continue
		}
		oldBlock := normalizeBlock(old.prose[key])
		if oldBlock == "" {
			continue
		}
		newBlock := normalizeBlock(newProse)
		for _, line := range old.prose[key] {
			normalized := normalizeLine(line)
			if normalized == "" {
				continue
			}
			// Default-value sentences migrate into @default markers and are
			// rendered as tag lines instead of prose.
			if defaultSentencePattern.MatchString(line) {
				rawBlock := strings.Join(newProse, "\n")
				if !defaultValueRetained(line, newBlock, rawBlock) {
					issues = append(issues, fmt.Errorf("%s [%s] %s: default not migrated (%s): %s", file, heading, key, lang, line))
				}
				continue
			}
			if !strings.Contains(newBlock, normalized) {
				issues = append(issues, fmt.Errorf("%s [%s] %s: missing line (%s): %s", file, heading, key, lang, line))
			}
		}
		_ = oldBlock
	}
	return issues
}
