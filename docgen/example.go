package docgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	exampleMarkerStart = "<!-- AUTO-GENERATED-EXAMPLE-START -->"
	exampleMarkerEnd   = "<!-- AUTO-GENERATED-EXAMPLE-END -->"

	structureHeading = "### Structure"
)

// originExampleLabels maps the struct origin of shared option fields to the
// ellipsis label rendered at the end of example blocks when the origin is
// omitted from the page (fields documented on the shared page instead). An
// empty label renders no marker: the origin's fields are removed legacy
// fields documented nowhere. Origins absent from this map are an error:
// extend the map, never guess.
var originExampleLabels = map[string][2]string{
	// [0] = en label, [1] = zh label
	"DialerOptions":         {"Dial Fields", "拨号字段"},
	"AbstractDialerOptions": {"Dial Fields", "拨号字段"},
	"ListenOptions":         {"Listen Fields", "监听字段"},
	"HTTP2Options":          {"HTTP2 Fields", "HTTP2 字段"},
	"QUICOptions":           {"QUIC Fields", "QUIC 字段"},
	"InboundOptions":        {"", ""},
}

// groupZhLabels resolves explicit example group labels on localized pages
// when the label is not a heading of the page itself (e.g. a cross-page
// reference like the HTTP2 Fields marker on shared/quic.md).
var groupZhLabels = map[string]string{
	"HTTP2 Fields":    "HTTP2 字段",
	"QUIC Fields":     "QUIC 字段",
	"UDP NAT Fields":  "UDP NAT 字段",
	"Provider Fields": "提供商字段",
}

// renderExample renders the inner lines of the generated Structure example
// block of a single-definition page.
func (g *generator) renderExample(page *Page, sections []Section, spec *ExampleSpec, lang string) ([]string, error) {
	var defRef DefRef
	count := 0
	for sectionIndex := range sections {
		for _, def := range sections[sectionIndex].Defs {
			defRef = def
			count++
		}
	}
	if count != 1 {
		return nil, fmt.Errorf("%s: example requires exactly one definition, got %d", page.File, count)
	}
	resolved, rootNode, err := g.schema.ResolveDef(defRef)
	if err != nil {
		return nil, err
	}
	if resolved.Root == "" {
		return nil, fmt.Errorf("%s: example def %s has no root", page.File, resolved.Name)
	}
	labelIndex := 0
	if lang == "zh" {
		labelIndex = 1
	}
	// Ellipsis markers for origins omitted from the page.
	omitOrigins, shared, err := g.omittedOrigins(resolved)
	if err != nil {
		return nil, err
	}
	originSlug := make(map[string]string, len(shared))
	for _, origin := range shared {
		originSlug[origin.Origin] = origin.Slug
	}
	full, err := g.ast.Flatten(resolved.Root, nil)
	if err != nil {
		return nil, err
	}
	depth1CountByOrigin := make(map[string]int)
	var originOrder []string
	originSeen := make(map[string]bool)
	for _, entry := range full {
		if strings.Contains(entry.Key, ".") {
			continue
		}
		if !originSeen[entry.Origin] {
			originSeen[entry.Origin] = true
			originOrder = append(originOrder, entry.Origin)
		}
		depth1CountByOrigin[entry.Origin]++
	}
	var ellipses []string
	for _, origin := range originOrder {
		if origin == resolved.Root || !omitOrigins[origin] || depth1CountByOrigin[origin] <= 1 {
			continue
		}
		if slug, isShared := originSlug[origin]; isShared {
			info, known := originSections[slug]
			if !known {
				return nil, fmt.Errorf("%s: example origin %s has no shared page %s", page.File, origin, slug)
			}
			label := info.enHeading
			if labelIndex == 1 {
				label = info.zhHeading
			}
			if label != "" && (len(ellipses) == 0 || ellipses[len(ellipses)-1] != label) {
				ellipses = append(ellipses, label)
			}
			continue
		}
		labels, known := originExampleLabels[origin]
		if !known {
			return nil, fmt.Errorf("%s: example origin %s has no example label, extend originExampleLabels", page.File, origin)
		}
		if labels[labelIndex] == "" {
			continue
		}
		if len(ellipses) == 0 || ellipses[len(ellipses)-1] != labels[labelIndex] {
			ellipses = append(ellipses, labels[labelIndex])
		}
	}
	// Explicit group labels for keys the page omits individually (omit_keys)
	// and documents in hand-maintained sections.
	for _, label := range g.resolveExampleGroups(page, spec.Groups, lang) {
		if len(ellipses) == 0 || ellipses[len(ellipses)-1] != label {
			ellipses = append(ellipses, label)
		}
	}
	// Rendered entries: omitted origins are dropped from the block and
	// referenced by an ellipsis label, except origins contributing a single
	// depth-1 key (option containers like OutboundTLSOptionsContainer),
	// which render as a placeholder value.
	entries, err := g.ast.Flatten(resolved.Root, nil)
	if err != nil {
		return nil, err
	}
	omit := make(map[string]bool)
	for _, key := range spec.Omit {
		omit[key] = true
	}
	type renderedEntry struct {
		lines      []string
		deprecated bool
	}
	var main, deprecated []renderedEntry
	addEntry := func(entry FlattenEntry) error {
		if omit[entry.Key] {
			return nil
		}
		value, err := g.exampleValue(spec, rootNode, entry)
		if err != nil {
			return fmt.Errorf("%s: %w", entry.Key, err)
		}
		comment := ""
		if spec.Comments != nil {
			comment = spec.Comments[entry.Key]
		}
		if lang == "zh" && spec.CommentsZh != nil {
			comment = spec.CommentsZh[entry.Key]
		}
		lines := renderExampleEntry(entry.Key, value, comment)
		if entry.Tags.Deprecated {
			deprecated = append(deprecated, renderedEntry{lines: lines, deprecated: true})
		} else {
			main = append(main, renderedEntry{lines: lines})
		}
		return nil
	}
	for _, entry := range entries {
		if strings.Contains(entry.Key, ".") {
			continue
		}
		if entry.Origin != resolved.Root && omitOrigins[entry.Origin] && !isExampleContainer(g, rootNode, entry.Key, depth1CountByOrigin[entry.Origin]) {
			continue
		}
		if err := addEntry(entry); err != nil {
			return nil, err
		}
	}
	// Registry discriminator and tag lead the block.
	var header []renderedEntry
	if resolved.Registry != "" {
		for _, key := range []string{"type", "tag"} {
			if omit[key] {
				continue
			}
			value, known := spec.Values[key]
			if !known {
				switch key {
				case "type":
					name := resolved.Name
					if node, found := g.schema.NodeForKey(rootNode, "type"); found {
						if constant := node.ConstString(); constant != "" {
							name = constant
						}
					}
					value = quoteJSON(name)
				default:
					if sample, ok := g.ast.TagSample(resolved.Root); ok {
						value = quoteJSON(sample)
						break
					}
					// The tag sample follows the corpus convention
					// <name>-in / <name>-out; other registries stay empty.
					switch resolved.Registry {
					case "Inbound":
						value = quoteJSON(resolved.Name + "-in")
					case "Outbound":
						value = quoteJSON(resolved.Name + "-out")
					default:
						value = `""`
					}
				}
			}
			comment := spec.Comments[key]
			header = append(header, renderedEntry{lines: renderExampleEntry(key, value, comment)})
		}
	}
	// Assemble; key lines carry trailing commas except the very last one.
	var all []renderedEntry
	all = append(all, header...)
	all = append(all, main...)
	if len(ellipses) > 0 {
		all = append(all, renderedEntry{lines: []string{""}})
		for _, label := range ellipses {
			all = append(all, renderedEntry{lines: []string{"  ... // " + label}})
		}
	}
	if len(deprecated) > 0 {
		all = append(all, renderedEntry{lines: []string{""}})
		all = append(all, renderedEntry{lines: []string{"  // Deprecated"}})
		all = append(all, deprecated...)
	}
	lastKey := -1
	for index, entry := range all {
		if len(entry.lines) == 0 {
			continue
		}
		if strings.HasPrefix(entry.lines[0], `  "`) {
			lastKey = index
		}
	}
	var lines []string
	lines = append(lines, "{")
	for index, entry := range all {
		entryLines := append([]string{}, entry.lines...)
		if index == lastKey && strings.HasSuffix(entryLines[len(entryLines)-1], ",") {
			entryLines[len(entryLines)-1] = strings.TrimSuffix(entryLines[len(entryLines)-1], ",")
		}
		lines = append(lines, entryLines...)
	}
	lines = append(lines, "}")
	return lines, nil
}

// resolveExampleGroups resolves explicit example group labels for the page
// language: each " / " separated part is translated by the localized twin
// heading of the same ordinal, or the literal groupZhLabels fallback.
func (g *generator) resolveExampleGroups(page *Page, groups []string, lang string) []string {
	if lang != "zh" || len(groups) == 0 {
		return groups
	}
	enLines := make([]string, 0)
	if enContent, err := os.ReadFile(filepath.Join(g.docsDir, filepath.FromSlash(page.Twin))); err == nil {
		enLines = strings.Split(string(enContent), "\n")
	}
	zhLines := make([]string, 0)
	if zhContent, err := os.ReadFile(filepath.Join(g.docsDir, filepath.FromSlash(page.File))); err == nil {
		zhLines = strings.Split(string(zhContent), "\n")
	}
	resolved := make([]string, 0, len(groups))
	for _, label := range groups {
		parts := strings.Split(label, " / ")
		for index, part := range parts {
			if heading, err := twinHeadingByOrdinal(enLines, "### "+part, zhLines); err == nil {
				parts[index] = strings.TrimPrefix(heading, "### ")
			} else if zh, known := groupZhLabels[part]; known {
				parts[index] = zh
			}
		}
		resolved = append(resolved, strings.Join(parts, " / "))
	}
	return resolved
}

// objectishNode reports whether a schema node resolves to an object shape,
// following references and union branches.
func objectishNode(node *SchemaNode, index *SchemaIndex) bool {
	for hop := 0; node != nil && hop < 8; hop++ {
		if node.TypeString() == "object" || node.Properties != nil {
			return true
		}
		switch {
		case node.Ref != "":
			node = index.DefNode(node.Ref)
		case len(node.AnyOf) > 0:
			for _, branch := range node.AnyOf {
				if objectishNode(branch, index) {
					return true
				}
			}
			return false
		case len(node.OneOf) > 0:
			for _, branch := range node.OneOf {
				if objectishNode(branch, index) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
	return false
}

// unionZeroValue returns the zero-value placeholder of a union node: arrays
// first, then integers, booleans, strings and objects, so listable marks and
// listable string fields keep their natural empty form without an explicit
// example value.
func unionZeroValue(node *SchemaNode, index *SchemaIndex) (string, bool) {
	kinds := make(map[string]bool)
	var collect func(*SchemaNode, int)
	collect = func(current *SchemaNode, depth int) {
		if current == nil || depth > 8 {
			return
		}
		if current.Ref != "" {
			collect(index.DefNode(current.Ref), depth+1)
			return
		}
		if kind := current.TypeString(); kind != "" {
			kinds[kind] = true
		}
		for _, branch := range current.AnyOf {
			collect(branch, depth+1)
		}
		for _, branch := range current.OneOf {
			collect(branch, depth+1)
		}
	}
	collect(node, 0)
	switch {
	case kinds["array"]:
		return `[]`, true
	case kinds["integer"], kinds["number"]:
		return `0`, true
	case kinds["boolean"]:
		return `false`, true
	case kinds["string"]:
		return `""`, true
	case kinds["object"]:
		return `{}`, true
	}
	return "", false
}

// isExampleContainer reports whether an omitted origin contributes a single
// option container key (e.g. OutboundTLSOptionsContainer contributes `tls`),
// which renders as a placeholder value instead of an ellipsis marker.
func isExampleContainer(g *generator, rootNode *SchemaNode, key string, depth1Count int) bool {
	if depth1Count != 1 {
		return false
	}
	node, found := g.schema.NodeForKey(rootNode, key)
	if !found || node == nil {
		return false
	}
	return node.TypeString() == "object" || objectishNode(node, g.schema)
}

// exampleValue renders the raw JSON value text of one entry.
func (g *generator) exampleValue(spec *ExampleSpec, rootNode *SchemaNode, entry FlattenEntry) (string, error) {
	if value, known := spec.Values[entry.Key]; known {
		return value, nil
	}
	if entry.Tags.Example != "" {
		return entry.Tags.Example, nil
	}
	node, found := g.schema.NodeForKey(rootNode, entry.Key)
	if !found || node == nil {
		return `""`, nil
	}
	if constant := node.ConstString(); constant != "" {
		return quoteJSON(constant), nil
	}
	switch node.TypeString() {
	case "string":
		// String enums keep the neutral placeholder; curated values are
		// encoded in the example spec.
		return `""`, nil
	case "integer", "number":
		enums, err := collectEnums(node, g.schema)
		if err != nil {
			return "", err
		}
		if len(enums) > 0 {
			return enums[0], nil
		}
		return `0`, nil
	case "boolean":
		return `false`, nil
	case "array":
		return `[]`, nil
	case "object":
		return `{}`, nil
	default:
		if zero, ok := unionZeroValue(node, g.schema); ok {
			return zero, nil
		}
		if objectishNode(node, g.schema) {
			return `{}`, nil
		}
		return `""`, nil
	}
}

// renderExampleEntry renders the lines of one example entry: the key line
// plus any continuation lines of a multi-line raw value.
func renderExampleEntry(key string, value string, comment string) []string {
	suffix := ","
	if comment != "" {
		suffix += " // " + comment
	}
	valueLines := strings.Split(strings.TrimRight(value, "\n"), "\n")
	lines := make([]string, 0, len(valueLines))
	lines = append(lines, `  "`+key+`": `+valueLines[0])
	// Continuation lines carry the value's own 2-space structure; shift them
	// to the key's nesting level so the JSON stays aligned.
	for _, line := range valueLines[1:] {
		if strings.TrimSpace(line) == "" {
			lines = append(lines, "")
			continue
		}
		lines = append(lines, "  "+line)
	}
	lines[len(lines)-1] += suffix
	return lines
}

func quoteJSON(text string) string {
	return `"` + text + `"`
}

// spliceExample replaces the Structure example fence with the generated
// block. Pages without a Structure fence are skipped.
func (g *generator) spliceExample(page *Page, sections []Section, lang string, lines []string) ([]string, error) {
	spec := page.ResolveExample(g.manifest)
	if spec == nil {
		return lines, nil
	}
	headingText := structureHeading
	var enLines []string
	if lang == "zh" {
		enContent, err := os.ReadFile(filepath.Join(g.docsDir, filepath.FromSlash(page.Twin)))
		if err != nil {
			return nil, err
		}
		enLines = strings.Split(string(enContent), "\n")
		resolved, err := twinHeadingByOrdinal(enLines, structureHeading, lines)
		if err != nil {
			// No Structure heading on the localized page: skip.
			return lines, nil
		}
		headingText = resolved
	}
	headingIndex := -1
	for index, line := range lines {
		if line == headingText {
			headingIndex = index
			break
		}
	}
	if headingIndex < 0 {
		return lines, nil
	}
	fenceStart := -1
	fenceEnd := -1
	for inner := headingIndex + 1; inner < len(lines); inner++ {
		line := lines[inner]
		if atxHeadingPattern.MatchString(line) {
			break
		}
		if fenceStart < 0 {
			if line == "```json" {
				fenceStart = inner
			}
			continue
		}
		if line == "```" {
			fenceEnd = inner
			break
		}
	}
	if fenceStart < 0 || fenceEnd < 0 {
		return lines, nil
	}
	for fenceStart > 0 && lines[fenceStart-1] == exampleMarkerStart {
		fenceStart--
	}
	for fenceEnd+1 < len(lines) && lines[fenceEnd+1] == exampleMarkerEnd {
		fenceEnd++
	}
	body, err := g.renderExample(page, sections, spec, lang)
	if err != nil {
		return nil, err
	}
	// Preserve any prose between the heading and the fence.
	head := append([]string{}, lines[:fenceStart]...)
	for len(head) > 0 && head[len(head)-1] == "" {
		head = head[:len(head)-1]
	}
	head = append(head, "")
	replacement := []string{exampleMarkerStart, "```json"}
	replacement = append(replacement, body...)
	replacement = append(replacement, "```", exampleMarkerEnd)
	tail := lines[fenceEnd+1:]
	for len(tail) > 0 && tail[0] == "" {
		tail = tail[1:]
	}
	if len(tail) > 0 {
		replacement = append(replacement, "")
	}
	result := append(head, replacement...)
	result = append(result, tail...)
	return result, nil
}
