// Package rule contains the configurations for the available rulesets. These are
// pre-compiled on start up for efficiency
package rule

import (
	"encoding/json"
	"regexp"
)

// Rule represents a rule that can be run against part of a patch or a filename
type Rule struct {
	Part        string         `json:"part"` // Only applicable to file types
	Type        string         `json:"type"`
	Pattern     string         `json:"pattern"`
	Caption     string         `json:"caption"`
	Description string         `json:"description"`
	Regex       *regexp.Regexp `json:"regex,omitempty"`
}

var (
	// Sets contain the available rulesets
	Sets map[string][]Rule
)

// Precompiles the rulesets. Panics if the rulesets can't be parsed.
func init() {
	Sets = make(map[string][]Rule)

	// Compile the filename level rules
	var gr []Rule
	if err := json.Unmarshal(gitrobJSON, &gr); err != nil {
		panic(err)
	}
	Sets["file"] = gr

	// Compile the patch line level rules
	var pr []Rule
	if err := json.Unmarshal(lineRulesJSON, &pr); err != nil {
		panic(err)
	}
	Sets["line"] = pr

	// For all regex types, precompile the regex
	for _, set := range []string{"file", "line"} {
		for i := range Sets[set] {
			if Sets[set][i].Type == "regex" {
				Sets[set][i].Regex = regexp.MustCompile(Sets[set][i].Pattern)
			}
		}
	}
}
