package docgen

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const durationPattern = `^[-+]?(((\d+(\.\d*)?|\.\d+)(ns|us|µs|μs|ms|s|m|h|d))+|0)$`

const (
	markerStart    = "<!-- AUTO-GENERATED-CONFIG-START -->"
	markerStartAlt = "<!-- 此区域由脚本自动生成，切勿手动修改 -->"
	markerEnd      = "<!-- AUTO-GENERATED-CONFIG-END -->"
)

var ipV4Sample = "192.0.2.1"
var ipV6Sample = "2001:db8::1"
var prefixV4Sample = "192.0.2.1/24"
var prefixV6Sample = "2001:db8::/32"

// collectEnums gathers all enum values of a node, including inside unions.
func collectEnums(node *SchemaNode, index *SchemaIndex) ([]string, error) {
	var values []string
	seen := make(map[string]bool)
	var walk func(node *SchemaNode, depth int) error
	walk = func(node *SchemaNode, depth int) error {
		if node == nil || depth > 8 {
			return nil
		}
		if node.Ref != "" {
			if target := index.DefNode(node.Ref); target != nil {
				return walk(target, depth+1)
			}
			return nil
		}
		for _, value := range node.Enum {
			text := fmt.Sprintf("%v", value)
			if !seen[text] {
				seen[text] = true
				values = append(values, text)
			}
		}
		for _, branch := range node.AnyOf {
			if err := walk(branch, depth+1); err != nil {
				return err
			}
		}
		for _, branch := range node.OneOf {
			if err := walk(branch, depth+1); err != nil {
				return err
			}
		}
		if node.Items != nil {
			return walk(node.Items, depth+1)
		}
		return nil
	}
	if err := walk(node, 0); err != nil {
		return nil, err
	}
	sort.Strings(values)
	return values, nil
}

func enumValuesEqual(collected []string, declared string) bool {
	var wanted []string
	seen := make(map[string]bool)
	for _, item := range strings.Split(declared, ",") {
		item = strings.TrimSpace(item)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		wanted = append(wanted, item)
	}
	sort.Strings(wanted)
	if len(collected) != len(wanted) {
		return false
	}
	for idx := range collected {
		if collected[idx] != wanted[idx] {
			return false
		}
	}
	return true
}

// renderTagLines renders @default/@enum/@deprecated/@see markers.
func renderTagLine(entries []FlattenEntry, node *SchemaNode, index *SchemaIndex, lang string) ([]string, error) {
	var tags FieldTags
	for _, entry := range entries {
		if entry.Tags.Default != "" {
			tags.Default = entry.Tags.Default
		}
		if entry.Tags.Enum != "" {
			tags.Enum = entry.Tags.Enum
		}
		if entry.Tags.Deprecated {
			tags.Deprecated = true
		}
		tags.See = append(tags.See, entry.Tags.See...)
	}
	var lines []string
	if tags.Default != "" {
		if lang == "zh" {
			lines = append(lines, "- 默认值：`"+tags.Default+"`")
		} else {
			lines = append(lines, "- Default: `"+tags.Default+"`")
		}
	}
	if tags.Enum != "" {
		collected, err := collectEnums(node, index)
		if err != nil {
			return nil, err
		}
		if len(collected) > 0 {
			if !enumValuesEqual(collected, tags.Enum) {
				return nil, fmt.Errorf("@enum %q does not match schema values %v", tags.Enum, collected)
			}
		} else {
			var rendered []string
			for _, item := range strings.Split(tags.Enum, ",") {
				item = strings.TrimSpace(item)
				if item == "" {
					continue
				}
				rendered = append(rendered, "`"+item+"`")
			}
			if lang == "zh" {
				lines = append(lines, "- 可能值："+strings.Join(rendered, "、"))
			} else {
				lines = append(lines, "- Possible values: "+strings.Join(rendered, ", "))
			}
		}
	}
	if tags.Deprecated {
		if lang == "zh" {
			lines = append(lines, "- 已废弃。")
		} else {
			lines = append(lines, "- Deprecated.")
		}
	}
	for _, see := range tags.See {
		if lang == "zh" {
			lines = append(lines, "- 参阅："+see)
		} else {
			lines = append(lines, "- See also: "+see)
		}
	}
	return lines, nil
}

// Documentation metadata markers accepted in Go doc comments and locale
// catalogs. They render back to mkdocs syntax so the published pages keep
// their original appearance:
//
//	@required                       -> ==Required== / ==必填==
//	@indicator: Server only         -> ==Server only==
//	@since 1.11.0                   -> !!! question "Since sing-box 1.11.0"
//	                                   (zh: !!! question "自 sing-box 1.11.0 起")
//	@note                           -> !!! note
//	@note: Some title               -> !!! note "Some title"
//	@note: ""                       -> !!! note ""
//	@admonition question "Free .."  -> !!! question "Free .." (generic form)
//	@tab "With auto_redirect"       -> === "With auto_redirect"
var (
	requiredMarker   = "@required"
	indicatorPrefix  = "@indicator: "
	admonitionPrefix = "@admonition "
	sincePrefix      = "@since "
	tabPrefix        = "@tab "
)

// admonitionTypes enumerates the mkdocs admonition types accepted by the
// shorthand markers (@note, @note: title).
var admonitionTypes = map[string]bool{
	"note": true, "abstract": true, "info": true, "tip": true,
	"success": true, "question": true, "warning": true, "failure": true,
	"danger": true, "bug": true, "example": true, "quote": true,
}

// shorthandAdmonition parses "@<type>" or "@<type>: <title>".
func shorthandAdmonition(line string) (admonitionType string, title string, titled bool, ok bool) {
	if !strings.HasPrefix(line, "@") {
		return "", "", false, false
	}
	rest := strings.TrimPrefix(line, "@")
	if index := strings.IndexByte(rest, ':'); index >= 0 {
		admonitionType = rest[:index]
		title = strings.TrimPrefix(rest[index+1:], " ")
		titled = true
	} else {
		admonitionType = rest
	}
	if !admonitionTypes[admonitionType] {
		return "", "", false, false
	}
	return admonitionType, title, titled, true
}

// renderAdmonition builds the mkdocs admonition line.
// A marker without a colon renders the default mkdocs title, `@<type>:`
// renders no title, and `@<type>: title` renders the given title.
func renderAdmonition(admonitionType string, title string, titled bool) string {
	if !titled {
		return "!!! " + admonitionType
	}
	if title == "" || title == `""` {
		return "!!! " + admonitionType + ` ""`
	}
	return "!!! " + admonitionType + ` "` + title + `"`
}

// markerRank orders the metadata blocks that documentation descriptions keep
// after the field description: required banner, status indicators, version
// hints, admonitions, then tabbed presentations. Zero means "not a marker".
func markerRank(line string) int {
	if line == "@required" {
		return 1
	}
	if strings.HasPrefix(line, "@indicator: ") {
		return 2
	}
	if strings.HasPrefix(line, "@since ") {
		return 3
	}
	if strings.HasPrefix(line, "@tab ") {
		return 5
	}
	if !strings.HasPrefix(line, "@") {
		return 0
	}
	name := strings.TrimPrefix(line, "@")
	if index := strings.IndexAny(name, ": "); index >= 0 {
		name = name[:index]
	}
	if name == "admonition" || admonitionTypes[name] {
		return 4
	}
	return 0
}

// reorderDescription moves leading marker blocks after the field description,
// so JSON schema consumers read the description first.
func reorderDescription(prose []string) []string {
	var groups [][]string
	index := 0
	for index < len(prose) {
		line := prose[index]
		if line == "" || line[0] == ' ' || line[0] == '\t' || markerRank(line) == 0 {
			break
		}
		block := []string{line}
		index++
		for index < len(prose) && (prose[index] == "" || prose[index][0] == ' ' || prose[index][0] == '\t') {
			block = append(block, prose[index])
			index++
		}
		for len(block) > 1 && block[len(block)-1] == "" {
			block = block[:len(block)-1]
		}
		groups = append(groups, block)
		for index < len(prose) && prose[index] == "" {
			index++
		}
	}
	if len(groups) == 0 || index >= len(prose) {
		return prose
	}
	rest := prose[index:]
	for len(rest) > 0 && rest[0] == "" {
		rest = rest[1:]
	}
	for len(rest) > 0 && rest[len(rest)-1] == "" {
		rest = rest[:len(rest)-1]
	}
	if len(rest) == 0 {
		return prose
	}
	// Fixed order: required, indicators, version hints, admonitions, tabs;
	// relative order inside one rank is preserved.
	ordered := make([][][]string, 6)
	for _, group := range groups {
		rank := markerRank(group[0])
		ordered[rank] = append(ordered[rank], group)
	}
	out := append([]string{}, rest...)
	for rank := 1; rank < len(ordered); rank++ {
		for _, block := range ordered[rank] {
			out = append(out, "")
			out = append(out, block...)
		}
	}
	return out
}

// hasCJK reports whether the text contains CJK characters, used to detect the
// authored language of localized prose.
func hasCJK(text string) bool {
	for _, ch := range text {
		if ch >= 0x4E00 && ch <= 0x9FFF || ch >= 0x3400 && ch <= 0x4DBF {
			return true
		}
	}
	return false
}

func expandProse(prose []string, pageLang string, proseLang string) []string {
	if len(prose) == 0 {
		return prose
	}
	expanded := make([]string, 0, len(prose))
	for _, line := range prose {
		trimmed := strings.TrimLeft(line, " \t")
		indent := line[:len(line)-len(trimmed)]
		switch {
		case trimmed == requiredMarker:
			if pageLang == "zh" {
				expanded = append(expanded, indent+"==必填==")
			} else {
				expanded = append(expanded, indent+"==Required==")
			}
		case strings.HasPrefix(trimmed, indicatorPrefix):
			expanded = append(expanded, indent+"=="+strings.TrimPrefix(trimmed, indicatorPrefix)+"==")
		case strings.HasPrefix(trimmed, sincePrefix):
			version := strings.TrimPrefix(trimmed, sincePrefix)
			if proseLang == "zh" {
				expanded = append(expanded, indent+`!!! question "自 sing-box `+version+` 起"`)
			} else {
				expanded = append(expanded, indent+`!!! question "Since sing-box `+version+`"`)
			}
		case strings.HasPrefix(trimmed, admonitionPrefix):
			// A quoted title renders verbatim; without quotes the marker is
			// emitted bare (mkdocs default title), distinguishing it from the
			// explicit empty title form `@admonition note ""`.
			expanded = append(expanded, indent+"!!! "+strings.TrimPrefix(trimmed, admonitionPrefix))
		case strings.HasPrefix(trimmed, tabPrefix):
			expanded = append(expanded, indent+"=== "+strings.TrimPrefix(trimmed, tabPrefix))
		case strings.HasPrefix(trimmed, "@") && !strings.HasPrefix(trimmed, "@default") &&
			!strings.HasPrefix(trimmed, "@enum") && !strings.HasPrefix(trimmed, "@deprecated") &&
			!strings.HasPrefix(trimmed, "@see"):
			if admonitionType, title, titled, ok := shorthandAdmonition(trimmed); ok {
				expanded = append(expanded, indent+renderAdmonition(admonitionType, title, titled))
			} else {
				expanded = append(expanded, line)
			}
		default:
			expanded = append(expanded, line)
		}
	}
	return expanded
}

// renderFieldBlock renders one `#### <key>` block in the requested language.
func renderFieldBlock(entryKey string, prose []string, tagLines []string) string {
	var block strings.Builder
	block.WriteString("#### " + entryKey + "\n\n")
	if len(prose) > 0 {
		block.WriteString(strings.Join(prose, "\n"))
		block.WriteString("\n")
	}
	if len(tagLines) > 0 {
		block.WriteString("\n")
		block.WriteString(strings.Join(tagLines, "\n"))
		block.WriteString("\n")
	}
	return strings.TrimSuffix(block.String(), "\n")
}

var headingPattern = regexp.MustCompile(`^### `)
var fieldHeadingPattern = regexp.MustCompile(`^#### `)

type sectionLocation struct {
	headingIndex int
	endIndex     int // exclusive
}

func locateSection(lines []string, heading string) (sectionLocation, error) {
	for index, line := range lines {
		if line == heading {
			end := len(lines)
			for inner := index + 1; inner < len(lines); inner++ {
				if headingPattern.MatchString(lines[inner]) {
					end = inner
					break
				}
			}
			return sectionLocation{headingIndex: index, endIndex: end}, nil
		}
	}
	return sectionLocation{}, fmt.Errorf("heading %q not found", heading)
}

type sectionParts struct {
	leading  []string
	trailing []string
	oldOrder []string
	oldProse map[string][]string
	replace  []int // line range [start,end) containing old generated blocks
}

func splitSection(lines []string, location sectionLocation, format string) (sectionParts, error) {
	parts := sectionParts{oldProse: make(map[string][]string)}
	content := lines[location.headingIndex+1 : location.endIndex]
	// Remove the blank line following the heading.
	if len(content) > 0 && content[0] == "" {
		content = content[1:]
	}
	start := -1
	end := -1
	inMarker := false
	for offset, line := range content {
		switch line {
		case markerStart:
			if start >= 0 {
				return parts, fmt.Errorf("nested auto-generated markers")
			}
			start = offset
			inMarker = true
			continue
		case markerStartAlt:
			if !inMarker {
				return parts, fmt.Errorf("stray marker %s", line)
			}
			continue
		case markerEnd:
			if !inMarker {
				return parts, fmt.Errorf("stray marker %s", line)
			}
			inMarker = false
			end = offset
			continue
		}
		if inMarker {
			// Field headings inside an existing generated region still
			// define the previous order, keeping re-runs stable.
			if fieldHeadingPattern.MatchString(line) {
				key := strings.TrimSpace(strings.TrimPrefix(line, "####"))
				if _, loaded := parts.oldProse[key]; !loaded {
					parts.oldOrder = append(parts.oldOrder, key)
					parts.oldProse[key] = nil
				}
			}
			continue
		}
		if fieldHeadingPattern.MatchString(line) {
			key := strings.TrimSpace(strings.TrimPrefix(line, "####"))
			if key == "" {
				return parts, fmt.Errorf("empty field heading at offset %d", offset)
			}
			if _, loaded := parts.oldProse[key]; !loaded {
				parts.oldOrder = append(parts.oldOrder, key)
				parts.oldProse[key] = nil
			}
		}
	}
	if inMarker {
		return parts, fmt.Errorf("unterminated auto-generated marker")
	}
	if start < 0 {
		// First generation: the generated region runs from the first field
		// heading to the end of the section. Prose before the first field
		// stays outside the markers.
		first := -1
		for offset, line := range content {
			if fieldHeadingPattern.MatchString(line) {
				first = offset
				break
			}
		}
		if first < 0 {
			parts.leading = trimBlankLines(content)
			return parts, nil
		}
		if first > 0 {
			parts.leading = trimBlankLines(content[:first])
		}
		start = first
		end = len(content) - 1
		for offset := first; offset <= end; offset++ {
			line := content[offset]
			if fieldHeadingPattern.MatchString(line) {
				key := strings.TrimSpace(strings.TrimPrefix(line, "####"))
				if _, loaded := parts.oldProse[key]; !loaded {
					parts.oldOrder = append(parts.oldOrder, key)
					parts.oldProse[key] = nil
				}
			}
		}
	} else {
		// Existing markers: leading prose sits before START, trailing prose
		// after END.
		if start > 0 {
			parts.leading = trimBlankLines(content[:start])
		}
		if end < len(content)-1 {
			parts.trailing = trimBlankLines(content[end+1:])
		}
	}
	// Offset math: +1 for the blank line removed after the heading, +1 for
	// the heading line itself.
	parts.replace = []int{location.headingIndex + 2 + start, location.headingIndex + 2 + end}
	return parts, nil
}

func trimBlankLines(lines []string) []string {
	result := append([]string{}, lines...)
	for len(result) > 0 && result[0] == "" {
		result = result[1:]
	}
	for len(result) > 0 && result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}
	return result
}

// renderSection produces the replacement lines for one section, preserving
// leading and trailing prose outside the generated markers.
func renderSection(parts sectionParts, blocks []string) []string {
	var lines []string
	if len(parts.leading) > 0 {
		lines = append(lines, parts.leading...)
		lines = append(lines, "")
	}
	lines = append(lines, markerStart, markerStartAlt, "")
	for _, block := range blocks {
		lines = append(lines, block, "")
	}
	lines = append(lines, markerEnd)
	if len(parts.trailing) > 0 {
		lines = append(lines, "")
		lines = append(lines, parts.trailing...)
	}
	return lines
}

// spliceSection replaces a section content region inside a line slice.
func spliceSection(lines []string, location sectionLocation, replacement []string) []string {
	head := append([]string{}, lines[:location.headingIndex+1]...)
	tail := append([]string{}, lines[location.endIndex:]...)
	result := append(head, "")
	result = append(result, replacement...)
	result = append(result, "")
	result = append(result, tail...)
	return result
}

// twinHeadingByOrdinal maps an English section heading to the localized twin
// heading with the same ordinal among level-3 headings.
func twinHeadingByOrdinal(enLines []string, enHeading string, zhLines []string) (string, error) {
	ordinal := 0
	seen := 0
	for _, line := range enLines {
		if !headingPattern.MatchString(line) {
			continue
		}
		seen++
		if line == enHeading {
			ordinal = seen
			break
		}
	}
	if ordinal == 0 {
		return "", fmt.Errorf("twin heading: %q not found in english page", enHeading)
	}
	seen = 0
	for _, line := range zhLines {
		if !headingPattern.MatchString(line) {
			continue
		}
		seen++
		if seen == ordinal {
			return line, nil
		}
	}
	return "", fmt.Errorf("twin heading: localized page has no level-3 heading #%d for %q", ordinal, enHeading)
}
