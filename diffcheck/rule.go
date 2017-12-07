package diffcheck

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

type ruleSetType string

const (
	fileType ruleSetType = "file"
	lineType ruleSetType = "line"
)

var (
	ruleSets map[ruleSetType][]Rule
)

// Precompiles the rulesets. Panics if the rulesets can't be parsed.
func init() {

	ruleSets = make(map[ruleSetType][]Rule)

	// Compile the filename level rules
	var gr []Rule
	if err := json.Unmarshal(gitrobJSON, &gr); err != nil {
		panic(err)
	}
	ruleSets[fileType] = gr

	// Compile the patch line level rules
	var pr []Rule
	if err := json.Unmarshal(lineRulesJSON, &pr); err != nil {
		panic(err)
	}
	ruleSets[lineType] = pr

	// For all regex types, precompile the regex
	for _, set := range []ruleSetType{fileType, lineType} {
		for i := range ruleSets[set] {
			if ruleSets[set][i].Type == "regex" {
				ruleSets[set][i].Regex = regexp.MustCompile(ruleSets[set][i].Pattern)
			}
		}
	}
}
