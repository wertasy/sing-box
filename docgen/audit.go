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
	anyLine bool
}

// extractOldSection reads a section of a baseline document.
func extractOldSection(lines []string, heading string, format string) (*oldSection, error) {
	location, err := locateSection(lines, heading)
	if err != nil {
		return nil, err
	}
	section := &oldSection{heading: heading, prose: make(map[string][]string)}
	var current string
	for _, line := range lines[location.headingIndex+1 : location.endIndex] {
		if line == "" {
			continue
		}
		section.anyLine = true
		switch {
		case fieldHeadingPattern.MatchString(line):
			current = strings.TrimSpace(strings.TrimPrefix(line, "####"))
			if _, loaded := section.prose[current]; !loaded {
				section.keys = append(section.keys, current)
				section.prose[current] = nil
			}
		case nestedHeadingPattern.MatchString(line):
			// Deeper headings are prose by default (e.g. "##### Request format"),
			// but dotted names are nested field keys (e.g. zh "##### platform.http_proxy.enabled").
			current = strings.TrimSpace(strings.TrimLeft(line, "#"))
			if _, loaded := section.prose[current]; !loaded {
				section.keys = append(section.keys, current)
				section.prose[current] = nil
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
			if _, loaded := section.prose[current]; !loaded {
				section.keys = append(section.keys, current)
				section.prose[current] = nil
			}
		case strings.HasPrefix(line, "### "):
			current = ""
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

func auditSectionCompare(file string, heading string, old *oldSection, newKeys map[string][]string, lang string) []error {
	var issues []error
	if len(old.keys) == 0 && old.anyLine {
		issues = append(issues, fmt.Errorf("%s [%s]: baseline section has content but no extractable keys", file, heading))
		return issues
	}
	for _, key := range old.keys {
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
