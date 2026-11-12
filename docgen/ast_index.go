package docgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// FieldTags holds the extension markers parsed from a field doc comment.
type FieldTags struct {
	Default        string
	Enum           string   // comma separated values from "@enum:"
	Deprecated     bool     // "@deprecated: <note>" or godoc "Deprecated:" line
	DeprecatedNote string   // optional replacement note, kept for reference
	See            []string // arguments of "@see:" lines, markdown fragments
	Example        string   // "@example: <value>" or an indented block, raw JSON value text
	Omit           bool     // "@omit[: pages]": not documented (on the listed pages)
	OmitPages      []string // page files for "@omit: <pages>"; empty means all pages
}

// FlattenEntry describes a single JSON field of a root struct, mirroring
// schema.Generator.FlattenStruct semantics.
type FlattenEntry struct {
	Key      string   // JSON key, dotted for expanded nested fields
	Segments []string // key path segments
	Origin   string   // name of the struct declaring the field
	Prose    []string // markdown prose lines from the doc comment
	Tags     FieldTags
	EnumTag  string   // struct `enum:"..."` tag, comma separated
	RefTag   string   // struct `reference:"..."` tag
	Field    string   // Go field name
	TypeExpr ast.Expr // field type expression, used for schema-less labels
}

// ASTIndex indexes option package struct definitions with doc comments.
type ASTIndex struct {
	fset        *token.FileSet
	structs     map[string]*ast.StructType
	aliases     map[string]string   // defined type -> underlying named struct
	types       map[string]ast.Expr // type name -> declared type expression
	files       map[string]string
	sharedSlugs map[string]string // type name -> @shared slug of its doc comment
	tagSamples  map[string]string // type name -> @tag example sample of its doc comment
}

// SharedOrigin is one option origin whose fields are documented on a shared
// page instead of the embedding page.
type SharedOrigin struct {
	Origin string // struct name declaring the fields
	Slug   string // @shared slug, e.g. "dial"
}

// resolve follows defined-type aliases to the underlying struct name.
func (i *ASTIndex) resolve(name string) string {
	for hop := 0; hop < 8; hop++ {
		target, found := i.aliases[name]
		if !found {
			return name
		}
		name = target
	}
	return name
}

// sharedSlugOf returns the @shared slug declared in a doc comment group.
func sharedSlugOf(docs ...*ast.CommentGroup) string {
	for _, doc := range docs {
		if doc == nil {
			continue
		}
		for _, comment := range doc.List {
			text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
			rest, found := cutPrefixFold(text, "@Shared")
			if !found {
				continue
			}
			slug := strings.TrimSpace(rest)
			if slug == "" || strings.ContainsAny(slug, " \t") {
				continue
			}
			return slug
		}
	}
	return ""
}

// tagSampleOf returns the @tag example sample declared in a doc comment.
func tagSampleOf(docs ...*ast.CommentGroup) string {
	for _, doc := range docs {
		if doc == nil {
			continue
		}
		for _, comment := range doc.List {
			text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
			rest, found := cutPrefixFold(text, "@Tag ")
			if !found {
				continue
			}
			sample := strings.TrimSpace(rest)
			if len(sample) >= 2 && sample[0] == sample[len(sample)-1] && (sample[0] == '"' || sample[0] == '\'') {
				sample = sample[1 : len(sample)-1]
			}
			if sample == "" {
				continue
			}
			return sample
		}
	}
	return ""
}

// TagSample returns the @tag example sample of a type, accepting the defined
// type and its underscore-prefixed underlying struct names.
func (i *ASTIndex) TagSample(name string) (string, bool) {
	for _, candidate := range []string{name, i.resolve(name), "_" + strings.TrimPrefix(name, "_"), strings.TrimPrefix(name, "_")} {
		if sample, ok := i.tagSamples[candidate]; ok {
			return sample, true
		}
	}
	return "", false
}

// SharedOrigins returns the ordered origins whose fields are documented on a
// shared page: every struct embedding a type marked `@shared <slug>`,
// including the origins embedded underneath it.
func (i *ASTIndex) SharedOrigins(root string) ([]SharedOrigin, error) {
	var out []SharedOrigin
	seen := make(map[string]bool)
	var collect func(name string, slug string)
	collect = func(name string, slug string) {
		name = i.resolve(name)
		if seen[name] {
			return
		}
		seen[name] = true
		out = append(out, SharedOrigin{Origin: name, Slug: slug})
		structType := i.structs[name]
		if structType == nil {
			return
		}
		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 {
				continue
			}
			embedded, err := embeddedTypeName(field.Type)
			if err != nil {
				continue
			}
			collect(embedded, slug)
		}
	}
	var walk func(name string) error
	walk = func(name string) error {
		name = i.resolve(name)
		structType := i.structs[name]
		if structType == nil {
			return fmt.Errorf("unknown struct option.%s", name)
		}
		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 {
				continue
			}
			embedded, err := embeddedTypeName(field.Type)
			if err != nil {
				return err
			}
			resolved := i.resolve(embedded)
			if slug, marked := i.sharedSlugs[resolved]; marked {
				collect(resolved, slug)
				continue
			}
			if err := walk(embedded); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk(root); err != nil {
		return nil, err
	}
	return out, nil
}

// LoadASTIndex parses all Go files in dir (package option).
func LoadASTIndex(dir string) (*ASTIndex, error) {
	idx := &ASTIndex{
		fset:        token.NewFileSet(),
		structs:     make(map[string]*ast.StructType),
		aliases:     make(map[string]string),
		types:       make(map[string]ast.Expr),
		files:       make(map[string]string),
		sharedSlugs: make(map[string]string),
		tagSamples:  make(map[string]string),
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := dir + "/" + name
		file, err := parser.ParseFile(idx.fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if file.Name.Name != "option" {
			continue
		}
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				idx.types[typeSpec.Name.Name] = typeSpec.Type
				if slug := sharedSlugOf(typeSpec.Doc, genDecl.Doc); slug != "" {
					idx.sharedSlugs[typeSpec.Name.Name] = slug
				}
				if sample := tagSampleOf(typeSpec.Doc, genDecl.Doc); sample != "" {
					idx.tagSamples[typeSpec.Name.Name] = sample
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					// Defined type over a named type (e.g. `type X _X`).
					if ident, isIdent := typeSpec.Type.(*ast.Ident); isIdent {
						idx.aliases[typeSpec.Name.Name] = ident.Name
					}
					continue
				}
				if _, loaded := idx.structs[typeSpec.Name.Name]; loaded {
					return nil, fmt.Errorf("duplicate struct %s in option package", typeSpec.Name.Name)
				}
				idx.structs[typeSpec.Name.Name] = structType
				idx.files[typeSpec.Name.Name] = name
			}
		}
	}
	return idx, nil
}

// embedOverrides holds prose overrides declared on an embedded field's doc
// comment: byField entries come from `@field <key>:` blocks, base carries a
// plain prose comment that applies when the embedded struct declares exactly
// one JSON field.
type embedOverrides struct {
	base     []string
	byField  map[string][]string
	examples map[string]string // per-field @example values declared on the embed
}

// parseEmbedOverrides reads prose overrides from an embedded field doc
// comment. `@override:` starts an embed-level block applying to single-field
// containers; `@field <key>:` starts a per-field block. Comments without
// these markers carry no overrides.
func parseEmbedOverrides(doc *ast.CommentGroup) (*embedOverrides, error) {
	if doc == nil {
		return nil, nil
	}
	comments := make([]string, 0, len(doc.List))
	for _, comment := range doc.List {
		text := strings.TrimPrefix(comment.Text, "//")
		comments = append(comments, strings.TrimPrefix(text, " "))
	}
	overrides := &embedOverrides{byField: make(map[string][]string), examples: make(map[string]string)}
	inBase := false
	var current string
	var currentLines []string
	flush := func() {
		if current == "" {
			return
		}
		lines := stripBlankEdges(stripCommonIndent(currentLines))
		if len(lines) > 0 {
			overrides.byField[current] = lines
		}
		current = ""
		currentLines = nil
	}
	for _, line := range comments {
		trimmed := strings.TrimSpace(line)
		if strings.EqualFold(trimmed, "@Override:") || strings.EqualFold(trimmed, "@Override") {
			flush()
			inBase = true
			continue
		}
		if key, ok := cutPrefixFold(trimmed, "@Override "); ok {
			flush()
			inBase = false
			current = strings.TrimSuffix(strings.TrimSpace(key), ":")
			if current == "" || strings.ContainsAny(current, " \t") {
				return nil, fmt.Errorf("@Override marker without a key: %q", line)
			}
			continue
		}
		if value, ok := cutPrefixFold(trimmed, "@Example: "); ok && current != "" {
			overrides.examples[current] = value
			continue
		}
		if current != "" {
			currentLines = append(currentLines, line)
			continue
		}
		if inBase {
			overrides.base = append(overrides.base, line)
		}
	}
	flush()
	overrides.base = stripBlankEdges(stripCommonIndent(overrides.base))
	if len(overrides.base) == 0 {
		overrides.base = nil
	}
	if len(overrides.byField) == 0 {
		overrides.byField = nil
	}
	if overrides.base == nil && overrides.byField == nil && len(overrides.examples) == 0 {
		return nil, nil
	}
	return overrides, nil
}

// stripCommonIndent removes the common leading indentation of a block, so
// embed overrides can be written visually indented below their marker while
// the rendered prose stays flush.
func stripCommonIndent(lines []string) []string {
	minimum := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if minimum < 0 || indent < minimum {
			minimum = indent
		}
	}
	if minimum <= 0 {
		return lines
	}
	out := make([]string, len(lines))
	for index, line := range lines {
		if strings.TrimSpace(line) == "" {
			out[index] = ""
			continue
		}
		out[index] = line[minimum:]
	}
	return out
}

// stripBlankEdges trims leading and trailing blank lines.
func stripBlankEdges(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// jsonFieldCount counts the documented JSON fields a struct declares itself
// (direct fields only, embedded structs excluded).
func (i *ASTIndex) jsonFieldCount(structName string) int {
	structName = i.resolve(structName)
	structType := i.structs[structName]
	if structType == nil {
		return 0
	}
	count := 0
	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		if !field.Names[0].IsExported() {
			continue
		}
		jsonName := ""
		if field.Tag != nil && strings.HasPrefix(field.Tag.Value, "`") {
			if value, err := strconv.Unquote(field.Tag.Value); err == nil {
				jsonName = reflect.StructTag(value).Get("json")
			}
		}
		tagName, _, _ := strings.Cut(jsonName, ",")
		if tagName == "-" {
			continue
		}
		count++
	}
	return count
}

// Flatten expands root following schema.Generator.FlattenStruct semantics.
// Embedded structs listed in omitOrigins are skipped entirely.
func (i *ASTIndex) Flatten(root string, omitOrigins []string) ([]FlattenEntry, error) {
	omit := make(map[string]bool)
	for _, origin := range omitOrigins {
		omit[origin] = true
	}
	var out []FlattenEntry
	index := make(map[string]int)
	var walk func(structName string, prefix []string, overrides *embedOverrides) error
	walk = func(structName string, prefix []string, overrides *embedOverrides) error {
		structName = i.resolve(structName)
		structType := i.structs[structName]
		if structType == nil {
			return fmt.Errorf("unknown struct option.%s", structName)
		}
		// A plain embed comment applies only to single-field containers.
		base := []string(nil)
		if overrides != nil && overrides.base != nil {
			if i.jsonFieldCount(structName) != 1 {
				return fmt.Errorf("struct %s: embed prose needs @field markers for multi-field structs", structName)
			}
			base = overrides.base
		}
		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 && !field.Names[0].IsExported() {
				continue
			}
			if len(field.Names) == 0 && !isExportedIdent(field.Type) {
				continue
			}
			jsonName := ""
			if field.Tag != nil && strings.HasPrefix(field.Tag.Value, "`") {
				if value, err := strconv.Unquote(field.Tag.Value); err == nil {
					jsonName = reflect.StructTag(value).Get("json")
				}
			}
			tagName, _, _ := strings.Cut(jsonName, ",")
			if tagName == "-" {
				continue
			}
			schemaTag := ""
			if field.Tag != nil && strings.HasPrefix(field.Tag.Value, "`") {
				if value, err := strconv.Unquote(field.Tag.Value); err == nil {
					schemaTag = reflect.StructTag(value).Get("schema")
				}
			}
			if schemaTag != "" && schemaTag != "omit" {
				return fmt.Errorf("unknown schema tag %q on %s", schemaTag, structName)
			}
			if len(field.Names) == 0 && tagName == "" {
				embedded, err := embeddedTypeName(field.Type)
				if err != nil {
					return fmt.Errorf("struct %s: %w", structName, err)
				}
				if omit[embedded] {
					continue
				}
				nestedOverrides, err := parseEmbedOverrides(field.Doc)
				if err != nil {
					return fmt.Errorf("struct %s: %w", structName, err)
				}
				if err := walk(embedded, prefix, nestedOverrides); err != nil {
					return err
				}
				continue
			}
			if tagName == "" {
				if len(field.Names) == 0 {
					continue
				}
				tagName = field.Names[0].Name
			}
			segments := append(append([]string{}, prefix...), tagName)
			prose, tags, err := commentProse(field.Doc)
			if err != nil {
				return fmt.Errorf("struct %s field %s: %w", structName, tagName, err)
			}
			if overrides != nil {
				if override, ok := overrides.byField[tagName]; ok {
					prose = override
				} else if base != nil {
					prose = base
				}
			}
			if overrides != nil && overrides.examples != nil && tags.Example == "" {
				if example, ok := overrides.examples[tagName]; ok {
					tags.Example = example
				}
			}
			entry := FlattenEntry{
				Key:      strings.Join(segments, "."),
				Segments: segments,
				Origin:   structName,
				Prose:    prose,
				Tags:     tags,
				Field:    fieldName(field),
				TypeExpr: field.Type,
			}
			if enumTag, ok := structTagName(field.Tag, "enum"); ok {
				entry.EnumTag = enumTag
			}
			if refTag, ok := structTagName(field.Tag, "reference"); ok {
				entry.RefTag = refTag
			}
			if pos, loaded := index[entry.Key]; loaded {
				out[pos] = entry // same semantics as badjson TypedMap.Put
			} else {
				index[entry.Key] = len(out)
				out = append(out, entry)
			}
		}
		return nil
	}
	if err := walk(root, nil, nil); err != nil {
		return nil, err
	}
	return out, nil
}

// Expand walks the given top level key of a flatten result through nested
// struct fields, producing dotted entries like peers.address.
func (i *ASTIndex) Expand(entries []FlattenEntry, key string) ([]FlattenEntry, error) {
	var top *FlattenEntry
	for idx := range entries {
		if entries[idx].Key == key {
			top = &entries[idx]
			break
		}
	}
	var prefix []string
	var startStruct string
	if top != nil {
		if len(top.Segments) != 1 {
			return nil, fmt.Errorf("expand key %s: only top level keys can be expanded", key)
		}
		startStruct = i.fieldStructName(top.Origin, top.Segments[0])
		prefix = []string{key}
	} else {
		// Fall back to a Go field name, for variant holders excluded from the
		// JSON schema with `json:"-"` whose sub-fields are decoded through a
		// custom unmarshaller.
		if len(entries) == 0 {
			return nil, fmt.Errorf("expand key %s: not found in struct fields", key)
		}
		rootStruct := i.resolve(entries[0].Origin)
		structType := i.structs[rootStruct]
		if structType == nil {
			return nil, fmt.Errorf("expand key %s: unknown struct option.%s", key, entries[0].Origin)
		}
		for _, field := range structType.Fields.List {
			if len(field.Names) != 1 || field.Names[0].Name != key {
				continue
			}
			startStruct = structNameOf(field.Type)
			break
		}
		if startStruct == "" {
			return nil, fmt.Errorf("expand key %s: not found in struct fields", key)
		}
	}
	var out []FlattenEntry
	var walk func(structName string, prefix []string, overrides *embedOverrides) error
	walk = func(structName string, prefix []string, overrides *embedOverrides) error {
		resolved := i.resolve(structName)
		structType := i.structs[resolved]
		if structType == nil {
			return fmt.Errorf("unknown struct option.%s", structName)
		}
		base := []string(nil)
		if overrides != nil && overrides.base != nil {
			if i.jsonFieldCount(resolved) != 1 {
				return fmt.Errorf("struct %s: embed prose needs @field markers for multi-field structs", resolved)
			}
			base = overrides.base
		}
		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 && !field.Names[0].IsExported() {
				continue
			}
			jsonName := ""
			if field.Tag != nil && strings.HasPrefix(field.Tag.Value, "`") {
				if value, err := strconv.Unquote(field.Tag.Value); err == nil {
					jsonName = reflect.StructTag(value).Get("json")
				}
			}
			tagName, _, _ := strings.Cut(jsonName, ",")
			if tagName == "-" {
				// Variant holders are not part of the JSON object itself but
				// their fields are decoded through custom unmarshallers.
				if len(field.Names) == 1 {
					if holder := structNameOf(field.Type); holder != "" && i.structs[i.resolve(holder)] != nil {
						if err := walk(holder, prefix, nil); err != nil {
							return err
						}
					}
				}
				continue
			}
			if len(field.Names) == 0 && tagName == "" {
				embedded, err := embeddedTypeName(field.Type)
				if err != nil {
					return fmt.Errorf("struct %s: %w", structName, err)
				}
				nestedOverrides, err := parseEmbedOverrides(field.Doc)
				if err != nil {
					return fmt.Errorf("struct %s: %w", structName, err)
				}
				if err := walk(embedded, prefix, nestedOverrides); err != nil {
					return err
				}
				continue
			}
			if tagName == "" {
				tagName = field.Names[0].Name
			}
			segments := append(append([]string{}, prefix...), tagName)
			prose, tags, err := commentProse(field.Doc)
			if err != nil {
				return fmt.Errorf("struct %s field %s: %w", structName, tagName, err)
			}
			if overrides != nil {
				if override, ok := overrides.byField[tagName]; ok {
					prose = override
				} else if base != nil {
					prose = base
				}
			}
			if overrides != nil && overrides.examples != nil && tags.Example == "" {
				if example, ok := overrides.examples[tagName]; ok {
					tags.Example = example
				}
			}
			entry := FlattenEntry{
				Key:      strings.Join(segments, "."),
				Segments: segments,
				Origin:   structName,
				Prose:    prose,
				Tags:     tags,
				Field:    fieldName(field),
				TypeExpr: field.Type,
			}
			if enumTag, ok := structTagName(field.Tag, "enum"); ok {
				entry.EnumTag = enumTag
			}
			if refTag, ok := structTagName(field.Tag, "reference"); ok {
				entry.RefTag = refTag
			}
			out = append(out, entry)
			if name := structNameOf(field.Type); name != "" && i.structs[name] != nil {
				if err := walk(name, segments, nil); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := walk(startStruct, prefix, nil); err != nil {
		return nil, err
	}
	return out, nil
}

// field returns the field of a struct matching the JSON key (falling back to
// the Go field name).
func (i *ASTIndex) field(structName string, jsonKey string) *ast.Field {
	structType := i.structs[structName]
	if structType == nil {
		return nil
	}
	for _, field := range structType.Fields.List {
		if len(field.Names) != 1 {
			continue
		}
		name := field.Names[0].Name
		if field.Tag != nil && strings.HasPrefix(field.Tag.Value, "`") {
			if value, err := strconv.Unquote(field.Tag.Value); err == nil {
				if jsonName, _, _ := strings.Cut(reflect.StructTag(value).Get("json"), ","); jsonName != "" {
					name = jsonName
				}
			}
		}
		if name == jsonKey {
			return field
		}
	}
	return nil
}

func (i *ASTIndex) fieldStructName(origin string, fieldName string) string {
	field := i.field(origin, fieldName)
	if field == nil {
		return ""
	}
	return structNameOf(field.Type)
}

func fieldName(field *ast.Field) string {
	if len(field.Names) == 1 {
		return field.Names[0].Name
	}
	return ""
}

func isExportedIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.IsExported()
}

func embeddedTypeName(expr ast.Expr) (string, error) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, nil
	case *ast.StarExpr:
		return embeddedTypeName(t.X)
	case *ast.SelectorExpr:
		return "", fmt.Errorf("embedded type %s is not in package option", exprToString(expr))
	default:
		return "", fmt.Errorf("unsupported embedded type %s", exprToString(expr))
	}
}

func structNameOf(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return structNameOf(t.X)
	case *ast.ArrayType:
		return structNameOf(t.Elt)
	default:
		return ""
	}
}

func exprToString(expr ast.Expr) string {
	var b strings.Builder
	formatExpr(&b, expr)
	return b.String()
}

func formatExpr(b *strings.Builder, expr ast.Expr) {
	switch t := expr.(type) {
	case *ast.Ident:
		b.WriteString(t.Name)
	case *ast.StarExpr:
		b.WriteString("*")
		formatExpr(b, t.X)
	case *ast.SelectorExpr:
		formatExpr(b, t.X)
		b.WriteString(".")
		b.WriteString(t.Sel.Name)
	case *ast.ArrayType:
		if t.Len != nil {
			b.WriteString("[")
			b.WriteString(exprToString(t.Len))
			b.WriteString("]")
		} else {
			b.WriteString("[]")
		}
		formatExpr(b, t.Elt)
	default:
		b.WriteString(fmt.Sprintf("%T", expr))
	}
}

func structTagName(tag *ast.BasicLit, name string) (string, bool) {
	if tag == nil || !strings.HasPrefix(tag.Value, "`") {
		return "", false
	}
	value, err := strconv.Unquote(tag.Value)
	if err != nil {
		return "", false
	}
	structTag := reflect.StructTag(value)
	got, ok := structTag.Lookup(name)
	return got, ok
}

// commentProse converts a doc comment group into markdown prose lines and
// extension tags.
func commentProse(doc *ast.CommentGroup) ([]string, FieldTags, error) {
	var tags FieldTags
	var prose []string
	if doc == nil {
		return nil, tags, nil
	}
	var inExample bool
	var exampleLines []string
	for _, comment := range doc.List {
		text := comment.Text
		var lines []string
		switch {
		case strings.HasPrefix(text, "//"):
			lines = []string{strings.TrimPrefix(strings.TrimPrefix(text, "//"), " ")}
		case strings.HasPrefix(text, "/*"):
			body := strings.TrimSuffix(strings.TrimPrefix(text, "/*"), "*/")
			for _, line := range strings.Split(body, "\n") {
				line = strings.TrimSpace(line)
				line = strings.TrimPrefix(line, "* ")
				lines = append(lines, line)
			}
		default:
			continue
		}
		for _, line := range lines {
			if line == "" {
				if len(prose) > 0 && prose[len(prose)-1] != "" {
					prose = append(prose, "")
				}
				continue
			}
			if strings.HasPrefix(line, "nolint:") {
				continue
			}
			if value, ok := cutPrefixFold(line, "@Default: "); ok {
				tags.Default = value
				continue
			}
			if value, ok := cutPrefixFold(line, "@Enum: "); ok {
				tags.Enum = value
				continue
			}
			if value, ok := cutPrefixFold(line, "@Deprecated: "); ok {
				tags.Deprecated = true
				if value != "true" {
					tags.DeprecatedNote = value
				}
				continue
			}
			if value, ok := cutPrefixFold(line, "@See: "); ok {
				tags.See = append(tags.See, value)
				continue
			}
			if value, ok := cutPrefixFold(line, "@Example: "); ok {
				tags.Example = value
				continue
			}
			if strings.EqualFold(line, "@Example:") {
				inExample = true
				continue
			}
			if inExample {
				if line == "" {
					inExample = false
					continue
				}
				if line[0] == ' ' {
					// Indented lines belong to the example value block.
					exampleLines = append(exampleLines, line)
					continue
				}
				inExample = false
			}
			if strings.EqualFold(line, "@Omit") {
				tags.Omit = true
				continue
			}
			if value, ok := cutPrefixFold(line, "@Omit:"); ok {
				tags.Omit = true
				for _, page := range strings.Split(value, ",") {
					page = strings.TrimSpace(page)
					if page != "" {
						tags.OmitPages = append(tags.OmitPages, page)
					}
				}
				if len(tags.OmitPages) == 0 {
					return nil, tags, fmt.Errorf("empty @omit page list")
				}
				continue
			}
			if value, ok := strings.CutPrefix(line, "Deprecated:"); ok {
				tags.Deprecated = true
				_ = value
				continue
			}
			prose = append(prose, line)
		}
	}
	for len(prose) > 0 && prose[len(prose)-1] == "" {
		prose = prose[:len(prose)-1]
	}
	if len(exampleLines) > 0 {
		tags.Example = strings.Join(stripCommonIndent(exampleLines), "\n")
	}
	return prose, tags, nil
}
