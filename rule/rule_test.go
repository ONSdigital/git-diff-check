package rule_test

import (
	"testing"

	"github.com/ONSdigital/git-diff-check/rule"
)

func TestInitRules(t *testing.T) {
	// The rulesets should be populated into rule.Sets as part of the package
	// init() method. We won't test for specific rules but the amount loaded
	// should be non-zero
	if len(rule.Sets) == 0 {
		t.Error("Failed to initialise rulesets")
	}
}
