package docgen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// indexTable describes one registry index table under docs/configuration:
// its first column lists the registry type constants.
type indexTable struct {
	file     string
	registry string
	// legacy lists first-column values allowed without a schema variant
	// (non-schema legacy rows).
	legacy map[string]bool
}

var indexTables = []indexTable{
	{file: "outbound/index.md", registry: "Outbound", legacy: map[string]bool{"wireguard": true, "dns": true}},
	{file: "inbound/index.md", registry: "Inbound"},
	{file: "endpoint/index.md", registry: "Endpoint"},
	{file: "service/index.md", registry: "Service", legacy: map[string]bool{
		// oom-killer is a registered service without a documentation page.
		"oom-killer": true,
	}},
	{file: "dns/server/index.md", registry: "DNSServer", legacy: map[string]bool{"empty (default)": true}},
}

// auditIndexTables validates the first column of the registry index tables
// against the schema variant sets.
func auditIndexTables(schema *SchemaIndex, docsDir string) []error {
	var issues []error
	for _, table := range indexTables {
		variants := schema.Variants[table.registry]
		if variants == nil {
			issues = append(issues, fmt.Errorf("%s: unknown registry %s", table.file, table.registry))
			continue
		}
		listed, err := readIndexTableTypes(filepath.Join(docsDir, filepath.FromSlash(table.file)))
		if err != nil {
			issues = append(issues, fmt.Errorf("%s: %w", table.file, err))
			continue
		}
		listedSet := make(map[string]bool, len(listed))
		for _, name := range listed {
			listedSet[name] = true
		}
		for _, name := range sortedVariantKeys(variants) {
			if table.legacy[name] {
				continue
			}
			if !listedSet[name] {
				issues = append(issues, fmt.Errorf("%s: index table missing type %s", table.file, name))
			}
		}
		for _, name := range listed {
			if _, known := variants[name]; !known {
				if table.legacy[name] {
					continue
				}
				issues = append(issues, fmt.Errorf("%s: index table has unknown type %s", table.file, name))
			}
		}
	}
	return issues
}

func sortedVariantKeys(variants map[string]*SchemaNode) []string {
	names := make([]string, 0, len(variants))
	for name := range variants {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// readIndexTableTypes parses the first column of a markdown index table.
func readIndexTableTypes(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, line := range strings.Split(string(content), "\n") {
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 3 {
			continue
		}
		first := strings.TrimSpace(cells[1])
		if first == "" || first == "Type" || strings.HasPrefix(first, "-") || strings.Contains(first, "---") {
			continue
		}
		first = strings.TrimPrefix(first, "`")
		first = strings.TrimSuffix(first, "`")
		names = append(names, first)
	}
	return names, nil
}
