package diffcheck_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/ONSdigital/git-diff-check/diffcheck"
)

type testCase struct {
	Name            string
	ExpectedReports []diffcheck.Report
	OK              bool
	Patch           []byte
}

func TestSnoopPatch(t *testing.T) {

	for _, tc := range testCases {

		t.Logf("Given a patch containing %s", tc.Name)
		t.Logf("  When the patch is snooped")
		ok, reports, err := diffcheck.SnoopPatch(tc.Patch)

		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		if ok != tc.OK {
			t.Errorf("Got result %v, but got %v", ok, tc.OK)
		}

		if len(reports) != len(tc.ExpectedReports) {
			t.Errorf("Incorrect number of reports, got %d, expected %d", len(reports), len(tc.ExpectedReports))
		}

		for i, expected := range tc.ExpectedReports {
			gotReport := reports[i]

			shouldEqual("path", gotReport.Path, expected.Path, t)
			shouldEqual("old path", gotReport.OldPath, expected.OldPath, t)

			if len(expected.Warnings) != len(gotReport.Warnings) {
				t.Errorf("Incorrect number of warnings in report, got %d, expected %d", len(gotReport.Warnings), len(expected.Warnings))
			}

			for j, expWarning := range expected.Warnings {
				gotWarning := gotReport.Warnings[j]

				shouldEqual("type", gotWarning.Type, expWarning.Type, t)
				shouldEqualInt("line", gotWarning.Line, expWarning.Line, t)
				shouldEqual("description", gotWarning.Description, expWarning.Description, t)
			}
		}
	}
}

func shouldEqual(field, got, expected string, t *testing.T) {
	t.Logf("    %s should equal %s", field, expected)
	if got != expected {
		t.Errorf("Expected %s to be %s, but got %s", field, expected, got)
	}
}

func shouldEqualInt(field string, got, expected int, t *testing.T) {
	t.Logf("    %s should equal %d", field, expected)
	if got != expected {
		t.Errorf("Expected %s to be %d, but got %d", field, expected, got)
	}
}

var testCases = []testCase{
	testCase{
		Name:            "a totally fine file",
		OK:              true,
		ExpectedReports: nil,
		Patch: []byte(`
diff --git a/diffcheck/README.md b/diffcheck/README.md
new file mode 100644
index 0000000..e69de29
			`),
	},
	testCase{
		Name: "a potentially bad filename",
		OK:   false,
		ExpectedReports: []diffcheck.Report{
			{
				Path:    "diffcheck/key.pem",
				OldPath: "diffcheck/key.pem",
				Warnings: []diffcheck.Warning{
					{
						Type:        "file",
						Line:        -1,
						Description: "Potential cryptographic private key",
					},
				},
			},
		},
		Patch: []byte(`
diff --git a/diffcheck/key.pem b/diffcheck/key.pem
new file mode 100644
index 0000000..e69de29
			`),
	},
	testCase{
		Name: "a potential AWS key",
		OK:   false,
		ExpectedReports: []diffcheck.Report{
			{
				Path:    ".aws/credentials",
				OldPath: ".aws/credentials",
				Warnings: []diffcheck.Warning{
					{
						Type:        "file",
						Line:        -1,
						Description: "AWS CLI credentials file",
					},
					{
						Type:        "file",
						Line:        -1,
						Description: "Contains word: credential",
					},
					{
						Type:        "line",
						Line:        6,
						Description: "Possible AWS Access Key",
					},
					{
						Type:        "line",
						Line:        7,
						Description: "Possible key in high entropy string",
					},
				},
			},
		},
		Patch: []byte(`
diff --git a/.aws/credentials b/.aws/credentials
index e69de29..92251f8 100644
--- a/.aws/credentials
+++ b/.aws/credentials
@@ -0,0 +4,6 @@

# Shhh
aws=AKIA7362373827372737
secret=ZWVTjPQSdhwRgl204Hc51YCsritMIzn8B=/p9UyeX7xu6KkAGqfm3FJ+oObLDNEva
		`),
	},
}

func ExampleSnoopPatch() {
	patch, _ := exec.Command("git", "diff", "-U0", "--staged").CombinedOutput()

	ok, reports, err := diffcheck.SnoopPatch(patch)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("WARNING! Potential sensitive data found:")
		for _, r := range reports {
			fmt.Printf("Found in (%s)\n", r.Path)
			for _, w := range r.Warnings {
				fmt.Printf("\t> [%s] %s (line %d)\n", w.Type, w.Description, w.Line)
			}
		}
	}
}
