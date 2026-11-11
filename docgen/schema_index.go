package docgen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/schema"
)

// SchemaNode is a JSON schema node decoded with ordered properties and typed
// additionalProperties wrappers.
type SchemaNode struct {
	Ref                   string                `json:"$ref,omitempty"`
	Type                  json.RawMessage       `json:"type,omitempty"`
	Const                 json.RawMessage       `json:"const,omitempty"`
	Enum                  []any                 `json:"enum,omitempty"`
	Pattern               string                `json:"pattern,omitempty"`
	Items                 *SchemaNode           `json:"items,omitempty"`
	Properties            *SchemaNodeMap        `json:"properties,omitempty"`
	Required              []string              `json:"required,omitempty"`
	PropertyNames         *SchemaNode           `json:"propertyNames,omitempty"`
	AdditionalProperties  *AdditionalProperties `json:"additionalProperties,omitempty"`
	UnevaluatedProperties json.RawMessage       `json:"unevaluatedProperties,omitempty"`
	AllOf                 []*SchemaNode         `json:"allOf,omitempty"`
	AnyOf                 []*SchemaNode         `json:"anyOf,omitempty"`
	OneOf                 []*SchemaNode         `json:"oneOf,omitempty"`
	Deprecated            bool                  `json:"deprecated,omitempty"`
	TagReference          string                `json:"x-tag-reference,omitempty"`
	Defs                  *SchemaNodeMap        `json:"$defs,omitempty"`
}

// TypeString returns the schema type: "string", "integer", "boolean",
// "number" or "object". Array types return their first element.
func (n *SchemaNode) TypeString() string {
	if len(n.Type) == 0 {
		return ""
	}
	var value any
	if err := json.Unmarshal(n.Type, &value); err != nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		for _, item := range typed {
			if text, ok := item.(string); ok {
				return text
			}
		}
	}
	return ""
}

// ConstString returns the const value as string for discriminator keys.
func (n *SchemaNode) ConstString() string {
	if len(n.Const) == 0 {
		return ""
	}
	var value string
	if err := json.Unmarshal(n.Const, &value); err != nil {
		return ""
	}
	return value
}

// SchemaNodeMap preserves JSON object property order.
type SchemaNodeMap struct {
	keys   []string
	values map[string]*SchemaNode
}

func (m *SchemaNodeMap) Put(key string, value *SchemaNode) {
	if m.values == nil {
		m.values = make(map[string]*SchemaNode)
	}
	if _, loaded := m.values[key]; !loaded {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
}

func (m *SchemaNodeMap) Get(key string) (*SchemaNode, bool) {
	value, loaded := m.values[key]
	return value, loaded
}

func (m *SchemaNodeMap) Keys() []string {
	return m.keys
}

func (m *SchemaNodeMap) UnmarshalJSON(content []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if token != json.Delim('{') {
		return fmt.Errorf("expected JSON object, got %v", token)
	}
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return err
		}
		key, ok := keyToken.(string)
		if !ok {
			return fmt.Errorf("expected object key, got %v", keyToken)
		}
		var node SchemaNode
		if err := decoder.Decode(&node); err != nil {
			return err
		}
		m.Put(key, &node)
	}
	_, err = decoder.Token()
	return err
}

// AdditionalProperties wraps the boolean or schema form of JSON Schema
// additionalProperties.
type AdditionalProperties struct {
	Bool   bool
	Schema *SchemaNode
}

func (a *AdditionalProperties) UnmarshalJSON(content []byte) error {
	trimmed := bytes.TrimSpace(content)
	switch {
	case bytes.Equal(trimmed, []byte("true")):
		a.Bool = true
		return nil
	case bytes.Equal(trimmed, []byte("false")):
		return nil
	}
	var node SchemaNode
	if err := json.Unmarshal(trimmed, &node); err != nil {
		return err
	}
	a.Schema = &node
	return nil
}

// SchemaIndex indexes the generated JSON schema for documentation lookups.
type SchemaIndex struct {
	Root     *SchemaNode
	Defs     map[string]*SchemaNode
	Variants map[string]map[string]*SchemaNode
}

// LoadSchemaIndex generates the JSON schema and indexes its definitions.
func LoadSchemaIndex() (*SchemaIndex, error) {
	content, err := schema.Generate(include.Context(context.Background()), reflect.TypeFor[option.Options]())
	if err != nil {
		return nil, err
	}
	var root SchemaNode
	if err := json.Unmarshal(content, &root); err != nil {
		return nil, err
	}
	index := &SchemaIndex{
		Root:     &root,
		Defs:     make(map[string]*SchemaNode),
		Variants: make(map[string]map[string]*SchemaNode),
	}
	if root.Defs != nil {
		for _, key := range root.Defs.Keys() {
			defNode, _ := root.Defs.Get(key)
			index.Defs[key] = defNode
			if len(defNode.OneOf) == 0 {
				continue
			}
			variants := make(map[string]*SchemaNode)
			for _, branch := range defNode.OneOf {
				objectNode := firstObject(branch)
				if objectNode == nil || objectNode.Properties == nil {
					continue
				}
				constValue := variantDiscriminator(objectNode)
				if constValue == "" {
					continue
				}
				variants[constValue] = objectNode
			}
			if len(variants) > 0 {
				index.Variants[key] = variants
			}
		}
	}
	return index, nil
}

func firstObject(node *SchemaNode) *SchemaNode {
	if node == nil {
		return nil
	}
	if len(node.OneOf) > 0 {
		for _, branch := range node.OneOf {
			if inner := firstObject(branch); inner != nil {
				return inner
			}
		}
		return nil
	}
	if node.TypeString() == "object" || node.Properties != nil {
		if node.Properties == nil {
			// Composed variants (e.g. rule matches via allOf) carry their
			// properties in the first object branch.
			for _, branch := range node.AllOf {
				if branch.Properties != nil {
					return branch
				}
			}
		}
		return node
	}
	return nil
}

// variantDiscriminator returns the registry variant key of a oneOf branch
// object: the const (or single meaningful enum value) of its "type" or
// "action" property.
func variantDiscriminator(objectNode *SchemaNode) string {
	for _, key := range []string{"type", "action"} {
		valueNode, found := objectNode.Properties.Get(key)
		if !found {
			continue
		}
		if constValue := valueNode.ConstString(); constValue != "" {
			return constValue
		}
		for _, value := range valueNode.Enum {
			if text, ok := value.(string); ok && text != "" {
				return text
			}
		}
	}
	return ""
}

// ResolveDef normalizes a definition reference: the registry is inferred when
// unambiguous, the Go root struct name is derived from the registry when
// omitted, and registry variants default to omitting the tag key.
func (s *SchemaIndex) ResolveDef(def DefRef) (DefRef, *SchemaNode, error) {
	if def.Registry == "" {
		var matches []string
		for registryName, variants := range s.Variants {
			if _, found := variants[def.Name]; found {
				matches = append(matches, registryName)
			}
		}
		switch len(matches) {
		case 0:
		case 1:
			def.Registry = matches[0]
		default:
			sort.Strings(matches)
			return def, nil, fmt.Errorf("def %s: matches registry variants %s, set registry field", def.Name, strings.Join(matches, ", "))
		}
	}
	node, err := s.RootNode(def)
	if err != nil {
		if def.Root != "" {
			// The definition has no schema node (deprecated options excluded
			// from the schema); generation falls back to the struct source.
			return def, nil, nil
		}
		return def, nil, err
	}
	if def.Registry == "" {
		if def.Root == "" {
			return def, nil, fmt.Errorf("def %s: root required for non-registry definitions", def.Name)
		}
		return def, node, nil
	}
	// The type discriminator is never documented; the tag key is omitted by
	// default. An explicit omit_keys list replaces the default set.
	omitKeys := []string{"type"}
	if def.OmitKeys == nil {
		omitKeys = append(omitKeys, "tag")
	} else {
		omitKeys = append(omitKeys, def.OmitKeys...)
	}
	def.OmitKeys = omitKeys
	if def.Root == "" {
		root, err := registryRootName(def.Registry, def.Name)
		if err != nil {
			return def, nil, err
		}
		def.Root = root
	}
	return def, node, nil
}

// registryRootName derives the option struct name of a registry variant.
func registryRootName(registryName string, name string) (string, error) {
	var optionsValue any
	var found bool
	switch registryName {
	case "Inbound":
		optionsValue, found = include.InboundRegistry().CreateOptions(name)
	case "Outbound":
		optionsValue, found = include.OutboundRegistry().CreateOptions(name)
	case "Endpoint":
		optionsValue, found = include.EndpointRegistry().CreateOptions(name)
	case "Service":
		optionsValue, found = include.ServiceRegistry().CreateOptions(name)
	case "DNSServer", "V2RayTransport":
		optionsValue, found = include.DNSTransportRegistry().CreateOptions(name)
	case "CertificateProvider":
		optionsValue, found = include.CertificateProviderRegistry().CreateOptions(name)
	default:
		return "", fmt.Errorf("def %s: cannot derive root for registry %s, set root field", name, registryName)
	}
	if !found || optionsValue == nil {
		return "", fmt.Errorf("def %s: registry %s has no options type, set root field", name, registryName)
	}
	optionsType := reflect.TypeOf(optionsValue)
	for optionsType.Kind() == reflect.Pointer {
		optionsType = optionsType.Elem()
	}
	if optionsType.Name() == "" {
		return "", fmt.Errorf("def %s: cannot derive root for registry %s, set root field", name, registryName)
	}
	return optionsType.Name(), nil
}

// RootNode returns the schema node for a manifest definition reference.
func (s *SchemaIndex) RootNode(def DefRef) (*SchemaNode, error) {
	if def.Registry != "" {
		variants := s.Variants[def.Registry]
		if variants == nil {
			return nil, fmt.Errorf("def %s: unknown registry def %s", def.Name, def.Registry)
		}
		node := variants[def.Name]
		if node == nil {
			return nil, fmt.Errorf("def %s: registry %s has no variant %s", def.Name, def.Registry, def.Name)
		}
		return node, nil
	}
	if node := s.Defs[def.Name]; node != nil {
		return node, nil
	}
	var matches []string
	for registryName, variants := range s.Variants {
		if _, found := variants[def.Name]; found {
			matches = append(matches, registryName)
		}
	}
	if len(matches) == 1 {
		return s.Variants[matches[0]][def.Name], nil
	}
	if len(matches) > 1 {
		sort.Strings(matches)
		return nil, fmt.Errorf("def %s: matches registry variants %s, set registry field", def.Name, strings.Join(matches, ", "))
	}
	return nil, fmt.Errorf("def %s: not found in schema $defs", def.Name)
}

// DefNode resolves a #/$defs/<name> reference.
func (s *SchemaIndex) DefNode(ref string) *SchemaNode {
	return s.Defs[strings.TrimPrefix(ref, "#/$defs/")]
}

// NodeForKey resolves a dotted key path inside a definition node, following
// references and listable/map wrappers.
func (s *SchemaIndex) NodeForKey(root *SchemaNode, key string) (*SchemaNode, bool) {
	current := root
	for _, segment := range strings.Split(key, ".") {
		next, found := s.propertyOf(current, segment)
		if !found {
			return nil, false
		}
		current = next
	}
	return current, true
}

func (s *SchemaIndex) propertyOf(node *SchemaNode, key string) (*SchemaNode, bool) {
	for hop := 0; node != nil && hop < 8; hop++ {
		if node.Properties != nil {
			if found, ok := node.Properties.Get(key); ok {
				return found, true
			}
			return nil, false
		}
		switch {
		case node.Ref != "":
			node = s.DefNode(node.Ref)
		case len(node.AnyOf) == 1:
			node = node.AnyOf[0]
		case len(node.OneOf) == 1:
			node = node.OneOf[0]
		case len(node.AnyOf)+len(node.OneOf) > 1:
			// Union forms (e.g. bool-or-object custom marshalers) expose
			// their keys through the single object branch.
			branches := append(append([]*SchemaNode{}, node.AnyOf...), node.OneOf...)
			var structured *SchemaNode
			for _, branch := range branches {
				if branch.Properties != nil {
					if structured != nil {
						return nil, false
					}
					structured = branch
				}
			}
			if structured == nil {
				return nil, false
			}
			node = structured
		case len(node.AllOf) > 0:
			for _, branch := range node.AllOf {
				if branch.Properties != nil {
					node = branch
					break
				}
			}
		case node.Items != nil:
			node = node.Items
		case node.AdditionalProperties != nil && node.AdditionalProperties.Schema != nil:
			node = node.AdditionalProperties.Schema
		default:
			return nil, false
		}
	}
	return nil, false
}

// SubNode resolves the object node containing the properties of a nested key
// prefix, used by manifest key expansion.
func (s *SchemaIndex) SubNode(root *SchemaNode, key string) (*SchemaNode, error) {
	current := root
	segments := strings.Split(key, ".")
	for idx, segment := range segments {
		next, found := s.propertyOf(current, segment)
		if !found {
			return nil, fmt.Errorf("key %s: schema node not found (at %s)", key, strings.Join(segments[:idx+1], "."))
		}
		current = next
	}
	return current, nil
}
