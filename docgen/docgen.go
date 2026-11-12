// Package docgen regenerates the configuration reference under
// docs/configuration from Go struct comments and the JSON schema.
//
// Generated field sections are wrapped in AUTO-GENERATED-CONFIG markers and
// audited against the baseline documentation to guarantee zero regression.
package docgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func sortStrings(values []string) {
	sort.Strings(values)
}

// DefaultBaseSha is the baseline commit containing the legacy documentation.
const DefaultBaseSha = "a13dbe1f7"

const docsDirRelative = "docs/configuration"
const optionDirRelative = "option"

type generator struct {
	check    bool
	base     string
	repoRoot string
	docsDir  string
	manifest *Manifest
	ast      *ASTIndex
	schema   *SchemaIndex
	locale   *Locale
}

// Run regenerates or verifies the configuration documentation. When only is
// non-empty, only the listed page files are processed and the completeness
// check is skipped.
func Run(check bool, base string, only []string) error {
	if base == "" {
		base = DefaultBaseSha
	}
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}
	manifest, err := LoadManifest()
	if err != nil {
		return err
	}
	docsDir := filepath.Join(repoRoot, docsDirRelative)
	if err := manifest.SynthesizeTwins(docsDir); err != nil {
		return err
	}
	if err := manifest.Validate(docsDir, len(only) == 0); err != nil {
		return err
	}
	if len(only) > 0 {
		var filtered []Page
		for _, page := range manifest.Pages {
			for _, pattern := range only {
				if matched, _ := filepath.Match(pattern, page.File); matched {
					filtered = append(filtered, page)
					break
				}
			}
		}
		if len(filtered) == 0 {
			return fmt.Errorf("--only %s: no matching pages in manifest", strings.Join(only, ","))
		}
		manifest.Pages = filtered
	}
	astIndex, err := LoadASTIndex(filepath.Join(repoRoot, optionDirRelative))
	if err != nil {
		return err
	}
	schemaIndex, err := LoadSchemaIndex()
	if err != nil {
		return err
	}
	if len(only) == 0 {
		issues := auditIndexTables(schemaIndex, docsDir)
		if len(issues) > 0 {
			for _, issue := range issues {
				fmt.Fprintln(os.Stderr, "AUDIT:", issue.Error())
			}
			return fmt.Errorf("documentation audit failed: %d issue(s)", len(issues))
		}
	}
	locale, err := loadLocale()
	if err != nil {
		return err
	}
	engine := &generator{
		check:    check,
		base:     base,
		repoRoot: repoRoot,
		docsDir:  docsDir,
		manifest: manifest,
		ast:      astIndex,
		schema:   schemaIndex,
		locale:   locale,
	}
	return engine.run()
}

func (g *generator) run() error {
	var issues []error
	updates := make(map[string]string)
	for index := range g.manifest.Pages {
		page := &g.manifest.Pages[index]
		pageIssues, content, err := g.renderPage(page)
		if err != nil {
			issues = append(issues, fmt.Errorf("%s: %w", page.File, err))
			continue
		}
		issues = append(issues, pageIssues...)
		if len(pageIssues) == 0 && content != "" {
			updates[page.File] = content
		}
	}
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Fprintln(os.Stderr, "AUDIT:", issue.Error())
		}
		return fmt.Errorf("documentation audit failed: %d issue(s)", len(issues))
	}
	if g.check {
		return nil
	}
	files := make([]string, 0, len(updates))
	for file := range updates {
		files = append(files, file)
	}
	sortStrings(files)
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(g.docsDir, filepath.FromSlash(file)), []byte(updates[file]), 0o644); err != nil {
			return err
		}
		fmt.Println("written", file)
	}
	return nil
}

// renderPage renders one page and audits it against the baseline version.
func (g *generator) renderPage(page *Page) ([]error, string, error) {
	sections, err := page.ResolveSections(g.manifest)
	if err != nil {
		return nil, "", err
	}
	lang := "en"
	if page.Twin != "" {
		lang = "zh"
	}
	currentPath := filepath.Join(g.docsDir, filepath.FromSlash(page.File))
	currentContent, err := os.ReadFile(currentPath)
	if err != nil {
		return nil, "", err
	}
	lines := strings.Split(string(currentContent), "\n")
	var twinLines []string
	for sectionIndex := range sections {
		if sections[sectionIndex].Heading != "" {
			continue
		}
		if page.Twin == "" {
			return nil, "", fmt.Errorf("section without heading on non-twin page %s", page.File)
		}
		twinContent, err := os.ReadFile(filepath.Join(g.docsDir, filepath.FromSlash(page.Twin)))
		if err != nil {
			return nil, "", err
		}
		twinLines = strings.Split(string(twinContent), "\n")
		break
	}
	var allIssues []error
	for sectionIndex := range sections {
		section := &sections[sectionIndex]
		if section.Heading == "" {
			resolved, err := twinHeadingByOrdinal(twinLines, section.EnHeading, lines)
			if err != nil {
				return nil, "", err
			}
			section.Heading = resolved
		}
		location, err := locateSection(lines, section.Heading)
		if err != nil {
			return nil, "", fmt.Errorf("section %q: %w", section.Heading, err)
		}
		parts, err := splitSection(lines, location, section.Format)
		if err != nil {
			return nil, "", fmt.Errorf("section %q: %w", section.Heading, err)
		}
		blocks, err := g.sectionBlocks(page.File, section, lang)
		if err != nil {
			return nil, "", fmt.Errorf("section %q: %w", section.Heading, err)
		}
		blocks = orderBlocks(blocks, parts.oldOrder)
		replacement := renderSection(parts, blocks)
		// Audit against the baseline document. Pages newly entering the
		// pipeline have no baseline section under the new heading; the audit
		// falls back to an empty baseline (additions are not regressions).
		baseline, err := g.baselinePage(page.File)
		if err != nil {
			return nil, "", err
		}
		baselineSection, err := extractOldSection(strings.Split(baseline, "\n"), section.Heading, section.Format)
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				return nil, "", fmt.Errorf("section %q: %w", section.Heading, err)
			}
			baselineSection = &oldSection{heading: section.Heading, prose: make(map[string][]string)}
		}
		newKeys := make(map[string][]string)
		for _, block := range blocks {
			newKeys[blockKey(block)] = blockProse(block)
		}
		allIssues = append(allIssues, auditSectionCompare(page.File, section.Heading, baselineSection, newKeys, lang)...)
		lines = spliceSection(lines, location, replacement)
	}
	lines, err = g.spliceOriginSections(page, sections, lang, lines)
	if err != nil {
		return nil, "", err
	}
	lines, err = g.spliceExample(page, sections, lang, lines)
	if err != nil {
		return nil, "", err
	}
	content := strings.Join(lines, "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if g.check {
		current := strings.TrimSuffix(string(currentContent), "\n") + "\n"
		if current != content {
			allIssues = append(allIssues, fmt.Errorf("%s: generated documentation is not up to date, run 'go run ./cmd/internal/generate_docs'", page.File))
			return allIssues, "", nil
		}
	}
	return allIssues, content, nil
}

// omittedOrigins returns the origins excluded from a page definition: the
// @shared origins inferred from the Go struct comments (unless re-included
// via include_origins) plus the explicit omit_origins list. The inferred
// shared origins are returned in walk order for link-section generation.
func (g *generator) omittedOrigins(def DefRef) (map[string]bool, []SharedOrigin, error) {
	omitted := make(map[string]bool)
	var shared []SharedOrigin
	if def.Root != "" {
		inferred, err := g.ast.SharedOrigins(def.Root)
		if err != nil {
			return nil, nil, err
		}
		for _, origin := range inferred {
			if contains(def.IncludeOrigins, origin.Origin) {
				continue
			}
			omitted[origin.Origin] = true
			shared = append(shared, origin)
		}
	}
	for _, origin := range def.OmitOrigins {
		omitted[origin] = true
	}
	return omitted, shared, nil
}

// omittedOnPage reports whether a field marked @omit is hidden on the page:
// without a page list it is hidden everywhere, otherwise the list matches the
// page file (localized twins match their English page).
func omittedOnPage(tags FieldTags, pageFile string) bool {
	if !tags.Omit {
		return false
	}
	if len(tags.OmitPages) == 0 {
		return true
	}
	for _, candidate := range []string{pageFile, englishPageName(pageFile)} {
		for _, page := range tags.OmitPages {
			if candidate == page || candidate == englishPageName(page) {
				return true
			}
		}
	}
	return false
}

// englishPageName maps a localized page file to its English twin name.
func englishPageName(file string) string {
	if strings.HasSuffix(file, ".zh.md") {
		return strings.TrimSuffix(file, ".zh.md") + ".md"
	}
	return file
}

// sectionBlocks renders the generated field blocks of one section in
// generation order.
func (g *generator) sectionBlocks(pageFile string, section *Section, lang string) ([]string, error) {
	type entryInfo struct {
		key    string
		prose  []string
		tags   []string
		root   string
		origin string
		field  string
	}
	entriesByKey := make(map[string]*entryInfo)
	var order []string
	add := func(info entryInfo) {
		if existing, loaded := entriesByKey[info.key]; loaded {
			if len(existing.prose) == 0 {
				existing.prose = info.prose
				existing.root = info.root
				existing.origin = info.origin
				existing.field = info.field
			}
			if len(existing.tags) == 0 {
				existing.tags = info.tags
			}
			return
		}
		copied := info
		entriesByKey[info.key] = &copied
		order = append(order, info.key)
	}
	for _, rawDef := range section.Defs {
		def, rootNode, err := g.schema.ResolveDef(rawDef)
		if err != nil {
			return nil, err
		}
		// The Go struct is the source of truth for documented keys; the schema
		// node, when present, only supplies type labels and enum checks.
		flattened, err := g.ast.Flatten(def.Root, nil)
		if err != nil {
			return nil, err
		}
		omit := make(map[string]bool)
		for _, key := range def.OmitKeys {
			omit[key] = true
		}
		omitOrigins, _, err := g.omittedOrigins(def)
		if err != nil {
			return nil, err
		}
		include := make(map[string]bool)
		for _, key := range def.IncludeKeys {
			include[key] = true
		}
		for _, astEntry := range flattened {
			if omittedOnPage(astEntry.Tags, pageFile) || omit[astEntry.Key] || omitOrigins[astEntry.Origin] {
				continue
			}
			if len(include) > 0 && !include[astEntry.Key] {
				continue
			}
			var node *SchemaNode
			if rootNode != nil {
				node, _ = g.schema.NodeForKey(rootNode, astEntry.Key)
			}
			info := entryInfo{key: astEntry.Key, root: def.Root, origin: astEntry.Origin, field: astEntry.Field, prose: astEntry.Prose}
			tags, err := renderTagLine([]FlattenEntry{astEntry}, node, g.schema, lang)
			if err != nil {
				return nil, fmt.Errorf("%s.%s: %w", def.Root, astEntry.Key, err)
			}
			info.tags = tags
			if override, ok := section.Overrides[astEntry.Key]; ok && lang == "en" {
				info.prose = splitProse(override.Prose)
			}
			add(info)
		}
		for _, expandKey := range section.Expand {
			expanded, err := g.ast.Expand(flattened, expandKey)
			if err != nil {
				return nil, fmt.Errorf("expand %s: %w", expandKey, err)
			}
			for _, entry := range expanded {
				if omittedOnPage(entry.Tags, pageFile) || omit[entry.Key] || omitOrigins[entry.Origin] {
					continue
				}
				if len(include) > 0 && !include[entry.Key] {
					continue
				}
				var node *SchemaNode
				if rootNode != nil {
					node, _ = g.schema.NodeForKey(rootNode, entry.Key)
				}
				tags, err := renderTagLine([]FlattenEntry{entry}, node, g.schema, lang)
				if err != nil {
					return nil, fmt.Errorf("%s.%s: %w", def.Root, entry.Key, err)
				}
				prose := entry.Prose
				if override, ok := section.Overrides[entry.Key]; ok && lang == "en" {
					prose = splitProse(override.Prose)
				}
				add(entryInfo{key: entry.Key, prose: prose, tags: tags, root: def.Root, origin: entry.Origin, field: entry.Field})
			}
		}
	}
	// Manifest extras are reserved for keys removed from the code: a struct
	// field with the same key must be documented from its comment instead.
	for _, key := range sortedExtraKeys(section.Extras) {
		if generated, exists := entriesByKey[key]; exists {
			return nil, fmt.Errorf("extra key %q is defined as %s.%s in the option structs, remove it from extras", key, generated.origin, generated.field)
		}
	}
	for _, key := range sortedExtraKeys(section.Extras) {
		extra := section.Extras[key]
		prose := splitProse(extra.Prose)
		if lang == "zh" {
			localeKey := extra.LocaleKey
			if localeKey == "" {
				localeKey = key
			}
			if description, found := g.locale.Lookup(localeKey); found {
				prose = splitProse(description)
			}
		}
		add(entryInfo{key: key, prose: prose})
	}
	// Locale pass for zh pages.
	if lang == "zh" {
		for _, key := range order {
			info := entriesByKey[key]
			var candidates []string
			for _, name := range localeNameVariants(info.root) {
				candidates = append(candidates, name+"."+key)
			}
			if info.origin != info.root {
				for _, name := range localeNameVariants(info.origin) {
					candidates = append(candidates, name+"."+key)
				}
			}
			if description, found := g.locale.Lookup(candidates...); found {
				info.prose = splitProse(description)
			}
			// Per-section zh prose variants (e.g. an option whose description
			// differs between the en and zh page).
			if override, ok := section.Overrides[key]; ok {
				info.prose = splitProse(override.Prose)
			}
		}
	}
	var blocks []string
	for _, key := range order {
		info := entriesByKey[key]
		proseLang := "en"
		if lang == "zh" && hasCJK(strings.Join(info.prose, "\n")) {
			proseLang = "zh"
		}
		blocks = append(blocks, renderFieldBlock(key, expandProse(markersFirst(info.prose), lang, proseLang), info.tags))
	}
	return blocks, nil
}

func orderBlocks(blocks []string, oldOrder []string) []string {
	if len(oldOrder) == 0 {
		return blocks
	}
	byKey := make(map[string]string)
	var keys []string
	for _, block := range blocks {
		key := blockKey(block)
		byKey[key] = block
		keys = append(keys, key)
	}
	var ordered []string
	seen := make(map[string]bool)
	for _, key := range oldOrder {
		if block, found := byKey[key]; found && !seen[key] {
			ordered = append(ordered, block)
			seen[key] = true
		}
	}
	for _, key := range keys {
		if !seen[key] {
			ordered = append(ordered, byKey[key])
			seen[key] = true
		}
	}
	return ordered
}

// localeNameVariants returns the struct name spellings used in locale keys:
// the name itself plus its underscore-prefixed or unprefixed counterpart
// (defined types and their underlying structs differ by a leading underscore).
func localeNameVariants(name string) []string {
	if name == "" {
		return nil
	}
	if trimmed := strings.TrimPrefix(name, "_"); trimmed != name {
		return []string{name, trimmed}
	}
	return []string{name, "_" + name}
}

func sortedExtraKeys(extras map[string]Extra) []string {
	keys := make([]string, 0, len(extras))
	for key := range extras {
		keys = append(keys, key)
	}
	sortStrings(keys)
	return keys
}

func splitProse(text string) []string {
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// baselinePage reads a documentation file at the baseline revision.
func (g *generator) baselinePage(file string) (string, error) {
	relPath := docsDirRelative + "/" + file
	output, err := exec.Command("git", "-C", g.repoRoot, "show", g.base+":"+relPath).Output()
	if err != nil {
		return "", fmt.Errorf("git show %s:%s failed (use --base with a ref containing the baseline): %w", g.base, relPath, err)
	}
	return string(output), nil
}

func blockKey(block string) string {
	for _, line := range strings.Split(block, "\n") {
		if strings.HasPrefix(line, "#### ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "#### "))
		}
	}
	return ""
}

// blockProse returns the prose lines of a rendered block, excluding the type
// line and rendered tag lines, as used by the documentation audit.
func blockProse(block string) []string {
	var prose []string
	for index, line := range strings.Split(block, "\n") {
		if index == 0 && strings.HasPrefix(line, "#### ") {
			continue
		}
		if strings.HasPrefix(line, "- ") && isTagLine(line) {
			continue
		}
		prose = append(prose, line)
	}
	for len(prose) > 0 && prose[0] == "" {
		prose = prose[1:]
	}
	for len(prose) > 0 && prose[len(prose)-1] == "" {
		prose = prose[:len(prose)-1]
	}
	return prose
}

func isTagLine(line string) bool {
	switch {
	case strings.HasPrefix(line, "- Default: `"),
		strings.HasPrefix(line, "- 默认值：`"),
		strings.HasPrefix(line, "- Possible values: "),
		strings.HasPrefix(line, "- 可能值："),
		line == "- Deprecated.",
		line == "- 已废弃。",
		strings.HasPrefix(line, "- See also: "),
		strings.HasPrefix(line, "- 参阅："):
		return true
	}
	return false
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found: run from the sing-box repository")
		}
		dir = parent
	}
}
