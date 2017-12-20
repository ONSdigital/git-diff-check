	"fmt"
	"os/exec"
	OK              bool
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
	t.Logf("    %s should equal %s", field, expected)
	if got != expected {
		t.Errorf("Expected %s to be %s, but got %s", field, expected, got)
}
func shouldEqualInt(field string, got, expected int, t *testing.T) {
	t.Logf("    %s should equal %d", field, expected)
	if got != expected {
		t.Errorf("Expected %s to be %d, but got %d", field, expected, got)
	}
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
		OK:   false,