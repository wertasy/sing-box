package docgen

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed manifests/*.yaml
var manifestFS embed.FS

// DefRef references a schema definition rendered into one section.
type DefRef struct {
	Name           string   `yaml:"name"`
	Root           string   `yaml:"root"`
	Registry       string   `yaml:"registry,omitempty"`
	OmitOrigins    []string `yaml:"omit_origins,omitempty"`
	IncludeOrigins []string `yaml:"include_origins,omitempty"`
	OmitKeys       []string `yaml:"omit_keys,omitempty"`
	IncludeKeys    []string `yaml:"include_keys,omitempty"`

	// DocsTag documents a real `tag` field of the options struct (wrapper
	// tags are omitted by default).
	DocsTag bool `yaml:"docs_tag,omitempty"`
}

// Extra declares a legacy documentation key without a schema counterpart.
type Extra struct {
	Prose     string `yaml:"prose"`
	LocaleKey string `yaml:"locale_key,omitempty"`
}

// Override replaces the English prose of a schema-backed key for this page
// (per-page variants of shared struct fields).
type Override struct {
	Prose string `yaml:"prose"`
}

// SectionPatch carries localized section deltas, keyed by the english section
// heading, so localized pages only declare sections that differ.
type SectionPatch struct {
	Defs       []DefRef            `yaml:"defs,omitempty"`
	Expand     []string            `yaml:"expand,omitempty"`
	Extras     map[string]Extra    `yaml:"extras,omitempty"`
	Overrides  map[string]Override `yaml:"overrides,omitempty"`
	Format     string              `yaml:"format,omitempty"`
	RemoveKeys []string            `yaml:"remove_keys,omitempty"`
}

// ExampleSpec customizes the generated `### Structure` example block of a
// single-definition page.
type ExampleSpec struct {
	Omit       []string          `yaml:"omit,omitempty"`        // json keys excluded from the block
	Values     map[string]string `yaml:"values,omitempty"`      // key → raw JSON value text, rendered verbatim
	Comments   map[string]string `yaml:"comments,omitempty"`    // key → inline JSON comment after the value
	CommentsZh map[string]string `yaml:"comments_zh,omitempty"` // zh overrides for Comments on localized pages
	Groups     []string          `yaml:"groups,omitempty"`      // extra ellipsis labels, e.g. "UDP NAT Fields"
}

// Section is one generated region below a heading.
type Section struct {
	Heading    string              `yaml:"heading"`
	Root       string              `yaml:"root,omitempty"`
	EnHeading  string              `yaml:"en_heading,omitempty"`
	Defs       []DefRef            `yaml:"defs,omitempty"`
	Expand     []string            `yaml:"expand,omitempty"`
	Extras     map[string]Extra    `yaml:"extras,omitempty"`
	Overrides  map[string]Override `yaml:"overrides,omitempty"`
	Format     string              `yaml:"format,omitempty"`
	RemoveKeys []string            `yaml:"remove_keys,omitempty"`
}

// Defaults provides per-file fallbacks so manifest entries only carry deltas.
type Defaults struct {
	Heading     string   `yaml:"heading,omitempty"`
	Registry    string   `yaml:"registry,omitempty"`
	Root        string   `yaml:"root,omitempty"`
	OmitOrigins []string `yaml:"omit_origins,omitempty"`
	OmitKeys    []string `yaml:"omit_keys,omitempty"`
	IncludeKeys []string `yaml:"include_keys,omitempty"`
}

// Page is one documentation file with its generated sections.
type Page struct {
	File      string                  `yaml:"file"`
	Twin      string                  `yaml:"twin,omitempty"`
	AutoTwin  bool                    `yaml:"-"`
	Defs      []DefRef                `yaml:"defs,omitempty"`
	Expand    []string                `yaml:"expand,omitempty"`
	Extras    map[string]Extra        `yaml:"extras,omitempty"`
	Overrides map[string]Override     `yaml:"overrides,omitempty"`
	Format    string                  `yaml:"format,omitempty"`
	Patches   map[string]SectionPatch `yaml:"patches,omitempty"`
	Sections  []Section               `yaml:"sections,omitempty"`
	Example   *ExampleSpec            `yaml:"example,omitempty"`
}

// ResolveExample returns the example spec of a page, inheriting it from the
// twin page for localized pages.
func (p *Page) ResolveExample(manifest *Manifest) *ExampleSpec {
	if p.Example != nil {
		return p.Example
	}
	if p.Twin == "" {
		return nil
	}
	for index := range manifest.Pages {
		if manifest.Pages[index].File == p.Twin {
			return manifest.Pages[index].Example
		}
	}
	return nil
}

// Manifest describes which documentation pages are generated. It is merged
// from all YAML files in docgen/manifests: pages entries are concatenated and
// manual lists are unioned. Each file's defaults apply to its own pages.
type Manifest struct {
	Defaults *Defaults `yaml:"defaults,omitempty"`
	Manual   []string  `yaml:"manual"`
	Pages    []Page    `yaml:"pages"`
}

//go:embed locale/zh_CN/*.yaml
var localeFS embed.FS

// localeDir holds the localized description catalogs of one language, one
// YAML file per maintenance group.
const localeDir = "locale/zh_CN"

// Locale holds Chinese descriptions keyed by "<Root>.<key>".
type Locale struct {
	entries map[string]string
}

// loadLocale merges all docgen/locale/zh_CN/*.yaml files. Each file maps
// locale keys to plain multi-line markdown descriptions (YAML block scalars).
func loadLocale() (*Locale, error) {
	locale := &Locale{entries: make(map[string]string)}
	paths, err := localeFS.ReadDir(localeDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range paths {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	for _, name := range names {
		content, err := localeFS.ReadFile(localeDir + "/" + name)
		if err != nil {
			return nil, err
		}
		var entries map[string]string
		if err := yaml.Unmarshal(content, &entries); err != nil {
			return nil, fmt.Errorf("docgen/%s/%s: %w", localeDir, name, err)
		}
		for key, description := range entries {
			if _, loaded := locale.entries[key]; loaded {
				return nil, fmt.Errorf("docgen/locale/%s: duplicated locale key %s", name, key)
			}
			locale.entries[key] = strings.TrimSuffix(description, "\n")
		}
	}
	return locale, nil
}

// Lookup returns the zh description for the first existing key.
func (l *Locale) Lookup(keys ...string) (string, bool) {
	for _, key := range keys {
		if description, found := l.entries[key]; found {
			return description, true
		}
	}
	return "", false
}

// LoadManifest decodes and merges docgen/manifests/*.yaml.
func LoadManifest() (*Manifest, error) {
	manifest := &Manifest{}
	paths, err := manifestFS.ReadDir("manifests")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range paths {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	for _, name := range names {
		content, err := manifestFS.ReadFile("manifests/" + name)
		if err != nil {
			return nil, err
		}
		var partial Manifest
		decoder := yaml.NewDecoder(bytes.NewReader(content))
		decoder.KnownFields(true)
		if err := decoder.Decode(&partial); err != nil {
			return nil, fmt.Errorf("docgen/manifests/%s: %w", name, err)
		}
		inferImplicitDefs(partial.Pages)
		applyPageDefaults(partial.Defaults, partial.Pages)
		applyDirectoryDefaults(partial.Pages)
		applyTwinConvention(partial.Pages)
		manifest.Manual = append(manifest.Manual, partial.Manual...)
		manifest.Pages = append(manifest.Pages, partial.Pages...)
	}
	sort.Strings(manifest.Manual)
	return manifest, nil
}

// inferImplicitDefs synthesizes the implicit definition of pages registered
// under a registry directory (inbound/, outbound/, endpoint/, service/,
// dns/server/) so manifests do not need placeholder `defs: [{}]` entries.
func inferImplicitDefs(pages []Page) {
	for index := range pages {
		page := &pages[index]
		if strings.HasSuffix(page.File, ".zh.md") {
			continue
		}
		if directoryRegistry(page.File) == "" {
			continue
		}
		if len(page.Sections) > 0 || len(page.Defs) > 0 {
			continue
		}
		page.Defs = []DefRef{{}}
	}
}

// applyPageDefaults merges per-file defaults into page entries so the YAML
// only carries deltas.
func applyPageDefaults(defaults *Defaults, pages []Page) {
	if defaults == nil {
		return
	}
	for pageIndex := range pages {
		page := &pages[pageIndex]
		if len(page.Sections) == 0 && (len(page.Defs) > 0 || len(page.Extras) > 0 || len(page.Overrides) > 0) {
			page.Sections = []Section{{
				Defs:      page.Defs,
				Expand:    page.Expand,
				Extras:    page.Extras,
				Overrides: page.Overrides,
				Format:    page.Format,
			}}
			page.Defs = nil
			page.Expand = nil
			page.Extras = nil
			page.Overrides = nil
			page.Format = ""
		}
		for sectionIndex := range page.Sections {
			section := &page.Sections[sectionIndex]
			if section.Heading == "" {
				section.Heading = defaults.Heading
			}
			for defIndex := range section.Defs {
				def := &section.Defs[defIndex]
				if def.Registry == "" {
					def.Registry = defaults.Registry
				}
				if def.Root == "" {
					def.Root = defaults.Root
				}
				if def.OmitOrigins == nil {
					def.OmitOrigins = defaults.OmitOrigins
				}
				if def.OmitKeys == nil {
					def.OmitKeys = defaults.OmitKeys
				}
				if def.IncludeKeys == nil {
					def.IncludeKeys = defaults.IncludeKeys
				}
			}
		}
	}
}

// applyDirectoryDefaults infers def registration from the page file path so
// manifests stay minimal: a def without an explicit registry that lives under
// inbound/, outbound/, endpoint/, service/ or dns/server/ registers with that
// directory's registry, and a registry def without an explicit name takes the
// page file stem. Explicit name/registry/root entries always win.
func applyDirectoryDefaults(pages []Page) {
	for pi := range pages {
		file := pages[pi].File
		registry := directoryRegistry(file)
		for si := range pages[pi].Sections {
			applyDirectoryDefaultsToDefs(&pages[pi].Sections[si].Defs, registry, file)
		}
		applyDirectoryDefaultsToDefs(&pages[pi].Defs, registry, file)
	}
}

func applyDirectoryDefaultsToDefs(defs *[]DefRef, registry string, file string) {
	stem := strings.TrimSuffix(filepath.Base(file), ".md")
	for di := range *defs {
		def := &(*defs)[di]
		if def.Root == "" && def.Registry == "" && registry != "" {
			def.Registry = registry
		}
		if def.Name == "" && def.Root == "" && def.Registry != "" {
			def.Name = stem
		}
	}
}

// directoryRegistry maps page directories to the registry their defs
// register with by default.
func directoryRegistry(file string) string {
	segments := strings.Split(file, "/")
	if len(segments) >= 3 && segments[0] == "dns" && segments[1] == "server" {
		return "DNSServer"
	}
	switch segments[0] {
	case "inbound":
		return "Inbound"
	case "outbound":
		return "Outbound"
	case "endpoint":
		return "Endpoint"
	case "service":
		return "Service"
	}
	return ""
}

// applyTwinConvention derives the twin reference of localized pages from the
// "<name>.zh.md" file name convention.
func applyTwinConvention(pages []Page) {
	for pageIndex := range pages {
		page := &pages[pageIndex]
		if page.Twin != "" || !strings.HasSuffix(page.File, ".zh.md") {
			continue
		}
		page.Twin = strings.TrimSuffix(page.File, ".zh.md") + ".md"
	}
}

// SynthesizeTwins registers localized pages for every generated page that has
// a .zh.md sibling and no explicit twin entry. Section headings of the twin
// are resolved by heading ordinal from the localized file.
func (m *Manifest) SynthesizeTwins(docsDir string) error {
	listed := make(map[string]bool)
	for _, page := range m.Pages {
		listed[page.File] = true
	}
	var synthesized []Page
	for _, page := range m.Pages {
		if page.Twin != "" || !strings.HasSuffix(page.File, ".md") || len(page.Sections) == 0 {
			continue
		}
		zhFile := strings.TrimSuffix(page.File, ".md") + ".zh.md"
		if listed[zhFile] {
			continue
		}
		if _, err := os.Stat(filepath.Join(docsDir, filepath.FromSlash(zhFile))); err != nil {
			continue
		}
		twin := Page{File: zhFile, Twin: page.File, AutoTwin: true}
		for _, section := range page.Sections {
			twin.Sections = append(twin.Sections, Section{EnHeading: section.Heading})
		}
		synthesized = append(synthesized, twin)
	}
	m.Pages = append(m.Pages, synthesized...)
	sort.SliceStable(m.Pages, func(i, j int) bool {
		if m.Pages[i].File == m.Pages[j].File {
			return m.Pages[i].Twin < m.Pages[j].Twin
		}
		return m.Pages[i].File < m.Pages[j].File
	})
	return nil
}

// Validate checks manifest internal consistency, and optionally completeness
// against the documentation directory.
func (m *Manifest) Validate(docsDir string, complete bool) error {
	pageSet := make(map[string]string) // file -> duplicate owner
	for index := range m.Pages {
		page := &m.Pages[index]
		if page.File == "" {
			return fmt.Errorf("pages[%d]: missing file", index)
		}
		if owner, loaded := pageSet[page.File]; loaded {
			return fmt.Errorf("pages[%d] %s: duplicated with %s", index, page.File, owner)
		}
		pageSet[page.File] = page.File
		if len(page.Sections) == 0 && page.Twin == "" {
			return fmt.Errorf("%s: missing sections", page.File)
		}
		if page.Example != nil {
			total := 0
			for sectionIndex := range page.Sections {
				total += len(page.Sections[sectionIndex].Defs)
			}
			if total != 1 {
				return fmt.Errorf("%s: example requires exactly one definition, got %d", page.File, total)
			}
		}
		for sectionIndex := range page.Sections {
			section := &page.Sections[sectionIndex]
			if section.Heading == "" && page.Twin == "" {
				return fmt.Errorf("%s sections[%d]: missing heading", page.File, sectionIndex)
			}
			if section.Format != "" && section.Format != "table" {
				return fmt.Errorf("%s sections[%d]: unknown format %q", page.File, sectionIndex, section.Format)
			}
		}
	}
	var twinSet []string
	for index := range m.Pages {
		page := &m.Pages[index]
		if page.Twin == "" {
			continue
		}
		var twin *Page
		for twinIndex := range m.Pages {
			if m.Pages[twinIndex].File == page.Twin {
				twin = &m.Pages[twinIndex]
				break
			}
		}
		if twin == nil {
			return fmt.Errorf("%s: twin %s not found in pages", page.File, page.Twin)
		}
		twinSet = append(twinSet, page.File)
	}
	var manualFiles []string
	for _, file := range m.Manual {
		if _, loaded := pageSet[file]; loaded {
			continue // pages entries win over shadowed manual entries
		}
		manualFiles = append(manualFiles, file)
	}
	if !complete {
		return nil
	}
	files, err := listDocumentationFiles(docsDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if _, loaded := pageSet[file]; loaded {
			continue
		}
		if strings.HasSuffix(file, ".zh.md") {
			enFile := strings.TrimSuffix(file, ".zh.md") + ".md"
			if _, loaded := pageSet[enFile]; loaded {
				return fmt.Errorf("%s: missing pages entry, twin of generated %s", file, enFile)
			}
		}
		if contains(manualFiles, file) {
			continue
		}
		return fmt.Errorf("%s: not listed in manifest pages or manual", file)
	}
	return nil
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func listDocumentationFiles(docsDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		relPath, err := filepath.Rel(docsDir, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(relPath))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// ResolveSections returns the effective sections of a page, inheriting
// definitions from the twin page for localized pages.
func (p *Page) ResolveSections(manifest *Manifest) ([]Section, error) {
	if p.Twin == "" {
		return p.Sections, nil
	}
	var twin *Page
	for index := range manifest.Pages {
		if manifest.Pages[index].File == p.Twin {
			twin = &manifest.Pages[index]
			break
		}
	}
	if twin == nil {
		return nil, fmt.Errorf("%s: twin %s not found", p.File, p.Twin)
	}
	if len(twin.Sections) == 0 {
		return nil, fmt.Errorf("%s: twin %s has no sections", p.File, p.Twin)
	}
	if len(p.Sections) == 0 {
		// Localized page entries may only carry patches: mirror the twin's
		// sections and resolve headings by ordinal at render time.
		for _, enSection := range twin.Sections {
			p.Sections = append(p.Sections, Section{EnHeading: enSection.Heading})
		}
	}
	sections := make([]Section, 0, len(p.Sections))
	for index := range p.Sections {
		section := p.Sections[index]
		var enSection *Section
		for twinIndex := range twin.Sections {
			candidate := &twin.Sections[twinIndex]
			if section.EnHeading != "" {
				if candidate.Heading == section.EnHeading {
					enSection = candidate
					break
				}
				continue
			}
			if len(p.Sections) == len(twin.Sections) {
				enSection = &twin.Sections[index]
				break
			}
		}
		if enSection == nil {
			if len(twin.Sections) == 1 && section.EnHeading == "" {
				enSection = &twin.Sections[0]
			} else {
				return nil, fmt.Errorf("%s section %q: cannot resolve en_heading against twin %s", p.File, section.Heading, p.Twin)
			}
		}
		section.EnHeading = enSection.Heading
		section.Root = enSection.Root
		section.Defs = enSection.Defs
		section.Expand = enSection.Expand
		section.Format = enSection.Format
		if section.Extras == nil {
			section.Extras = enSection.Extras
		}
		if patch, patched := p.Patches[section.EnHeading]; patched {
			if patch.Defs != nil {
				section.Defs = patch.Defs
			}
			if patch.Expand != nil {
				section.Expand = patch.Expand
			}
			if patch.Extras != nil {
				section.Extras = patch.Extras
			}
			if patch.Overrides != nil {
				section.Overrides = patch.Overrides
			}
			if patch.Format != "" {
				section.Format = patch.Format
			}
			if patch.RemoveKeys != nil {
				section.RemoveKeys = patch.RemoveKeys
			}
		}
		sections = append(sections, section)
	}
	return sections, nil
}
