package docgen

import (
	"fmt"
	"regexp"
)

// originInfo describes the boilerplate link section of one shared option
// origin: the en/zh heading texts and the target documentation paths.
type originInfo struct {
	enHeading string
	zhHeading string
	enLink    string
	zhLink    string
}

// originSections maps the struct origin of shared option fields to their
// link-section rendering. Secondary origins absent from this map keep their
// hand-maintained sections (Multiplex, HTTP Client, UDP NAT Fields, ...).
var originSections = map[string]originInfo{
	"dial": {
		enHeading: "Dial Fields",
		zhHeading: "拨号字段",
		enLink:    "/configuration/shared/dial/",
		zhLink:    "/zh/configuration/shared/dial/",
	},
	"listen": {
		enHeading: "Listen Fields",
		zhHeading: "监听字段",
		enLink:    "/configuration/shared/listen/",
		zhLink:    "/zh/configuration/shared/listen/",
	},
}

const (
	originMarkerStart = "<!-- AUTO-GENERATED-ORIGIN-START -->"
	originMarkerEnd   = "<!-- AUTO-GENERATED-ORIGIN-END -->"
)

// originExcludedFiles lists pages that keep hand-maintained origin sections:
// their sections carry extra notes beyond the link line.
var originExcludedFiles = map[string]bool{
	"endpoint/tailscale.md":         true,
	"endpoint/tailscale.zh.md":      true,
	"endpoint/masque-server.md":     true,
	"endpoint/masque-server.zh.md":  true,
	"endpoint/openvpn-server.md":    true,
	"endpoint/openvpn-server.zh.md": true,
	"inbound/tailcat.md":            true,
	"inbound/tailcat.zh.md":         true,
	"outbound/tailcat.md":           true,
	"outbound/tailcat.zh.md":        true,
	"service/usbip-client.md":       true,
	"service/usbip-client.zh.md":    true,
	"service/usbip-server.md":       true,
	"service/usbip-server.zh.md":    true,
}

var atxHeadingPattern = regexp.MustCompile(`^#{1,6} `)

// renderOriginBody renders the canonical link line of one origin section.
func renderOriginBody(info originInfo, lang string) string {
	if lang == "zh" {
		return fmt.Sprintf("参阅 [%s](%s) 了解详情。", info.zhHeading, info.zhLink)
	}
	return fmt.Sprintf("See [%s](%s) for details.", info.enHeading, info.enLink)
}

// pageOrigins returns the ordered unique shared-page slugs a page links to:
// the @shared origins omitted from the page definitions, filtered to the
// known link-section catalog.
func (g *generator) pageOrigins(sections []Section) ([]string, error) {
	var slugs []string
	seen := make(map[string]bool)
	for sectionIndex := range sections {
		for _, rawDef := range sections[sectionIndex].Defs {
			def, _, err := g.schema.ResolveDef(rawDef)
			if err != nil {
				return nil, err
			}
			_, shared, err := g.omittedOrigins(def)
			if err != nil {
				return nil, err
			}
			for _, origin := range shared {
				if _, known := originSections[origin.Slug]; !known {
					continue
				}
				if seen[origin.Slug] {
					continue
				}
				seen[origin.Slug] = true
				slugs = append(slugs, origin.Slug)
			}
		}
	}
	return slugs, nil
}

// spliceOriginSections replaces the content of generated origin link
// sections, appending the heading at EOF when the page documents the origin
// fields but omits the section.
func (g *generator) spliceOriginSections(page *Page, sections []Section, lang string, lines []string) ([]string, error) {
	if originExcludedFiles[page.File] {
		return lines, nil
	}
	origins, err := g.pageOrigins(sections)
	if err != nil {
		return nil, err
	}
	for _, origin := range origins {
		info := originSections[origin]
		headingText := info.enHeading
		if lang == "zh" {
			headingText = info.zhHeading
		}
		heading := "### " + headingText
		headingIndex := -1
		for index, line := range lines {
			if line == heading {
				headingIndex = index
				break
			}
		}
		replacement := []string{"", originMarkerStart, renderOriginBody(info, lang), originMarkerEnd}
		if headingIndex < 0 {
			lines = append(lines, "", heading, replacement[1], replacement[2], replacement[3])
			continue
		}
		endIndex := len(lines)
		for inner := headingIndex + 1; inner < len(lines); inner++ {
			if atxHeadingPattern.MatchString(lines[inner]) {
				endIndex = inner
				break
			}
		}
		if endIndex < len(lines) {
			// Keep one blank line before the following heading.
			replacement = append(replacement, "")
		}
		rebuilt := append([]string{}, lines[:headingIndex+1]...)
		rebuilt = append(rebuilt, replacement...)
		rebuilt = append(rebuilt, lines[endIndex:]...)
		lines = rebuilt
	}
	return lines, nil
}
