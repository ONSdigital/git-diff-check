package diffcheck_test

import (
	"testing"

	"github.com/ONSdigital/git-diff-check/diffcheck"

	. "github.com/smartystreets/goconvey/convey"
)

type testCase struct {
	Name            string
	ExpectedReports []diffcheck.Report
	Patch           []byte
}

func TestSnoopPatch(t *testing.T) {

	var patch []byte

	for _, tc := range testCases {

		Convey("Given a patch containing "+tc.Name, t, func() {
			patch = tc.Patch
			Convey("When the patch is snooped", func() {
				ok, reports, err := diffcheck.SnoopPatch(patch)

				Convey("Reports should be returned", func() {
					So(err, ShouldBeNil)
					So(ok, ShouldBeFalse)
					So(reports, ShouldResemble, tc.ExpectedReports)
				})
			})
		})

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
