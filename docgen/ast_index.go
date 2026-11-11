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
	Default    string
	Enum       string   // comma separated values from "@enum:"
	Deprecated bool     // "@deprecated: true" or godoc "Deprecated:" line
	See        []string // arguments of "@see:" lines, markdown fragments
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
	fset    *token.FileSet
	structs map[string]*ast.StructType
	aliases map[string]string   // defined type -> underlying named struct
	types   map[string]ast.Expr // type name -> declared type expression
	files   map[string]string
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

// LoadASTIndex parses all Go files in dir (package option).
func LoadASTIndex(dir string) (*ASTIndex, error) {
	idx := &ASTIndex{
		fset:    token.NewFileSet(),
		structs: make(map[string]*ast.StructType),
		aliases: make(map[string]string),
		types:   make(map[string]ast.Expr),
		files:   make(map[string]string),
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

// Flatten expands root following schema.Generator.FlattenStruct semantics.
// Embedded structs listed in omitOrigins are skipped entirely.
func (i *ASTIndex) Flatten(root string, omitOrigins []string) ([]FlattenEntry, error) {
	omit := make(map[string]bool)
	for _, origin := range omitOrigins {
		omit[origin] = true
	}
	var out []FlattenEntry
	index := make(map[string]int)
	var walk func(structName string, prefix []string) error
	walk = func(structName string, prefix []string) error {
		structName = i.resolve(structName)
		structType := i.structs[structName]
		if structType == nil {
			return fmt.Errorf("unknown struct option.%s", structName)
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
				if err := walk(embedded, prefix); err != nil {
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
	if err := walk(root, nil); err != nil {
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
	var walk func(structName string, prefix []string) error
	walk = func(structName string, prefix []string) error {
		resolved := i.resolve(structName)
		structType := i.structs[resolved]
		if structType == nil {
			return fmt.Errorf("unknown struct option.%s", structName)
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
						if err := walk(holder, prefix); err != nil {
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
				if err := walk(embedded, prefix); err != nil {
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
				if err := walk(name, segments); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := walk(startStruct, prefix); err != nil {
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
			if value, ok := strings.CutPrefix(line, "@default: "); ok {
				tags.Default = value
				continue
			}
			if value, ok := strings.CutPrefix(line, "@enum: "); ok {
				tags.Enum = value
				continue
			}
			if value, ok := strings.CutPrefix(line, "@deprecated: "); ok {
				if value != "true" {
					return nil, tags, fmt.Errorf("invalid @deprecated value %q", value)
				}
				tags.Deprecated = true
				continue
			}
			if value, ok := strings.CutPrefix(line, "@see: "); ok {
				tags.See = append(tags.See, value)
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
	return prose, tags, nil
}
