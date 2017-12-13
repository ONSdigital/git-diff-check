package diffcheck_test

import (
	"testing"

	"github.com/ONSdigital/git-diff-check/diffcheck"
)

type testCase struct {
	Name            string
	ExpectedReports []diffcheck.Report
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
		if ok {
			t.Errorf("Expected warning to be flagged, but got: %v", ok)
		}

		if len(reports) != len(tc.ExpectedReports) {
			t.Errorf("Incorrect number of reports, got %d, expected %d", len(reports), len(tc.ExpectedReports))
		}

		for i, expected := range tc.ExpectedReports {
			gotReport := reports[i]

			shouldEqual("path", gotReport.Path, expected.Path, t)
			shouldEqual("old path", gotReport.OldPath, expected.OldPath, t)

			if len(expected.Warnings) != len(gotReport.Warnings) {
				t.Errorf("Incorrect number of warnings in report, got %d, expected %d", len(expected.Warnings), len(gotReport.Warnings))
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
	if got != expected {
		t.Errorf("Expected %s to be %s, but got %s", field, expected, got)
	}
}

func shouldEqualInt(field string, got, expected int, t *testing.T) {
	if got != expected {
		t.Errorf("Expected %s to be %d, but got %d", field, expected, got)
	}
}

var testCases = []testCase{
	testCase{
		Name: "a potentially bad filename",
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
		`),
	},
}
