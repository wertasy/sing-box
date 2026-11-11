package schema

import (
	"context"
	"strconv"
	"strings"

	"github.com/sagernet/sing/service"
)

// FieldComment carries the documentation comment information of one struct
// field, provided to the generator by an optional CommentSource.
type FieldComment struct {
	// Description is the markdown description of the field.
	Description string
	// Default is the raw textual default value, coerced to the field node
	// type by the generator.
	Default string
	// Deprecated marks the field as deprecated.
	Deprecated bool
}

// CommentSource supplies documentation comments for struct fields while the
// JSON schema is generated. Implementations live outside this package
// because Go doc comments are not available through reflection.
type CommentSource interface {
	// StructFieldComment returns the comment information of the field with
	// the given JSON key (or Go name) in the named struct.
	StructFieldComment(structName string, jsonKey string) (FieldComment, bool)
}

// CommentSourceFromContext returns the optional comment source.
func CommentSourceFromContext(ctx context.Context) (CommentSource, bool) {
	commentSource := service.FromContext[CommentSource](ctx)
	return commentSource, commentSource != nil
}

// ContextWithCommentSource returns a context with the comment source set.
func ContextWithCommentSource(ctx context.Context, commentSource CommentSource) context.Context {
	return service.ContextWith(ctx, commentSource)
}

// applyFieldComment attaches description/default/deprecated information to a
// freshly built field node.
func applyFieldComment(node *Node, comment FieldComment) {
	if comment.Description != "" {
		node.Description = comment.Description
	}
	if comment.Deprecated {
		node.Deprecated = true
	}
	if comment.Default == "" {
		return
	}
	switch nodeTypeName(node) {
	case "integer":
		if value, err := strconv.ParseInt(strings.TrimSpace(comment.Default), 10, 64); err == nil {
			node.Default = value
		}
	case "boolean":
		if value, err := strconv.ParseBool(strings.TrimSpace(comment.Default)); err == nil {
			node.Default = value
		}
	default:
		node.Default = comment.Default
	}
}

func nodeTypeName(node *Node) string {
	typeValue, _ := node.Type.(string)
	return typeValue
}
