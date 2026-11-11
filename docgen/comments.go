package docgen

import (
	"strings"

	"github.com/sagernet/sing-box/schema"
)

// SchemaCommentSource supplies the schema generator with field documentation
// from the option package AST index.
type SchemaCommentSource struct {
	index *ASTIndex
}

// NewSchemaCommentSource parses the option package and returns a comment
// source for schema generation.
func NewSchemaCommentSource() (schema.CommentSource, error) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return nil, err
	}
	index, err := LoadASTIndex(repoRoot + "/" + optionDirRelative)
	if err != nil {
		return nil, err
	}
	return SchemaCommentSource{index: index}, nil
}

// StructFieldComment implements schema.CommentSource.
func (s SchemaCommentSource) StructFieldComment(structName string, jsonKey string) (schema.FieldComment, bool) {
	structType := s.index.structs[structName]
	if structType == nil {
		return schema.FieldComment{}, false
	}
	field := s.index.field(structName, jsonKey)
	if field == nil {
		return schema.FieldComment{}, false
	}
	prose, tags, err := commentProse(field.Doc)
	if err != nil {
		return schema.FieldComment{}, false
	}
	comment := schema.FieldComment{
		Deprecated: tags.Deprecated,
		Default:    tags.Default,
	}
	if len(prose) > 0 {
		// Keep the metadata markers as written in the Go comments (@note,
		// @since, @required, ...); only the markdown pages render them.
		// The field description comes first, informational markers follow.
		description := strings.Join(reorderDescription(prose), "\n")
		for _, see := range tags.See {
			description += "\n\nSee also: " + see
		}
		comment.Description = description
	}
	if comment.Description == "" && comment.Default == "" && !comment.Deprecated {
		return schema.FieldComment{}, false
	}
	return comment, true
}
