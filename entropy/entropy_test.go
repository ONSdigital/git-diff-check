// Package entropy_test is for testing the entropy package
//
// !IMPORTANT! - none of the keys or strings listed in this file are real keys. They
// 				 are generated purely to test this package and MUST NOT be used
//				 anywhere else as actual credentials
package entropy_test

import (
	"fmt"
	"testing"

	"github.com/ONSdigital/git-diff-check/entropy"
)

var (
	highBase64 = [][]byte{
		[]byte(`ZWVTjPQSdhwRgl204Hc51YCsritMIzn8B=/p9UyeX7xu6KkAGqfm3FJ+oObLDNEva`),
		[]byte(`hSXAQy9D1J0hkCQy0tKBCxnpcOQCPeM54RFXZLJE`),
		[]byte(`secret=ZWVTjPQSdhwRgl204Hc51YCsritMIzn8B=/p9UyeX7xu6KkAGqfm3FJ+oObLDNEva`),
		[]byte(`aws:ZWVTjPQSdhwRgl204Hc51YCsritMIzn8B=/p9UyeX7xu6KkAGqfm3FJ+oObLDNEva`),
	}
	highHex = [][]byte{
		[]byte(`b3A0a1FDfe86dcCE945B72`),
	}
)

func ExampleCalculateShannon() {
	password := []byte("verysecret")
	if entropy.CalculateShannon(password) < entropy.Base64Threshold {
		fmt.Println("Password not complex enough!")
	}
}

func TestCalculateEntropy(t *testing.T) {

	t.Log("Check base64 data")
	for _, b := range highBase64 {
		if e := entropy.CalculateShannon(b); e < entropy.Base64Threshold {
			t.Errorf("Got entropy %f, expected > %f", e, entropy.Base64Threshold)
		}
	}

	t.Log("Check hex data")
	for _, b := range highHex {
		if e := entropy.CalculateShannon(b); e < entropy.HexThreshold {
			t.Errorf("Got entropy %f, expected > %f", e, entropy.HexThreshold)
		}
	}
}

func TestCheck(t *testing.T) {

	exampleBlock := []byte(`+// CheckPatchLine takes a line from a patch hunk and tests it for naughty patterns
+func CheckPatchLine(line []byte) (bool, []Warning) {
+       warnings := []Warning{}
+		aws := []byte("hSXAQy9D1J0hkCQy0tKBCxnpcOQCPeM54RFXZLJE")
+
+ 		// Log in with secret: ZWVTjPQSdhwRgl204Hc51YCsritMIzn8B=/p9UyeX7xu6KkAGqfm3FJ+oObLDNEva
+
+       for _, rule := range linePatterns {
+               if found := rule.Regex.FindAll(line, -1); len(found) > 0 {
+                       for _, f := range found {
+                               // TODO Ignore exclusions
+                               if string(f) != "b3A0a1FDfe86dcCE945B72" {
+                                       warnings = append(warnings, Warning{Type: "line", Line: -1})
+                               }`)

	ok, n := entropy.Check(exampleBlock)
	if ok {
		t.Error("Expected 'not ok'")
	}
	if n != 3 {
		t.Errorf("Expected warnings, got %d, expected 3", n)
	}

	for _, b := range highBase64 {
		ok, _ := entropy.Check(b)
		if ok {
			t.Error("Expected failed base64 entropy check")
		}
	}

	for _, b := range highHex {
		ok, _ := entropy.Check(b)
		if ok {
			t.Error("Expected failed hex entropy check")
		}
	}

}
