package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ONSdigital/git-diff-check/diffcheck"
)

const (
	// Return codes - a non-zero will cause the commit hook to reject the commit
	accepted = 0
	rejected = 1
)

// Repository defines the github repo where the source code is located
const Repository = "ONSdigital/git-diff-check"

var target = flag.String("p", "", "(optional) path to repository")
var showVersion bool
var showHelp bool

func init() {
	flag.BoolVar(&showVersion, "version", false, "show current version")
	flag.BoolVar(&showHelp, "help", false, "show usage")
}

// Version is injected at build time
var Version string

func main() {

	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if showVersion {
		if len(Version) == 0 {
			fmt.Println(errors.New("No version set in binary! You may have a broken release"))
			os.Exit(1)
		}
		fmt.Println(Version)
		os.Exit(0)
	}

	// Attempt to check for a new version and inform the user if this is so.
	// If we can't connect or get the version for some reason then this is non-fatal
	versionCheck()

	if *target == "" {
		*target = "."
	}
	fmt.Printf("Running precommit diff check on '%s'\n", *target)

	// Import environmental feature flags
	if useEntropyFeature := os.Getenv("DC_ENTROPY_EXPERIMENT"); useEntropyFeature == "1" {
		fmt.Println("i) Experimental entropy checking enabled")
		diffcheck.UseEntropy = true
	}

	// Get where we are so we can get back
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("Couldn't get current dir:", err)
	}
	here := filepath.Dir(ex)

	err = os.Chdir(*target)
	if err != nil {
		log.Fatal("Failed to change to target dir:", err)
	}
	patch, err := exec.Command("git", "diff", "-U0", "--staged").CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run git command: %v (%s)", err, patch)
	}
	os.Chdir(here)

	if len(patch) == 0 {
		fmt.Println("No changes to test - exiting")
		os.Exit(accepted)
	}

	ok, reports, err := diffcheck.SnoopPatch(patch)
	if err != nil {
		log.Fatal("Failed to snoop:", err)
	}

	if len(reports) > 0 {
		fmt.Println("WARNING! Potential sensitive data found:")

		for _, r := range reports {
			fmt.Printf("Found in (%s)\n", r.Path)
			for _, w := range r.Warnings {
				if w.Type == "line" {
					fmt.Printf("\t> [%s] %s (line %d)\n", w.Type, w.Description, w.Line)
				} else {
					fmt.Printf("\t> [%s] %s\n", w.Type, w.Description)
				}
			}
			fmt.Println()
		}
	}

	if ok {
		fmt.Println("Diff probably ok!")
		os.Exit(accepted)
	}
	fmt.Println("If you're VERY SURE these files are ok, rerun commit with --no-verify")
	os.Exit(rejected)
}

// VersionResponse is the response from the github verson call
type VersionResponse struct {
	TagName string `json:"tag_name"`
}

func versionCheck() bool {

	// TODO
	// Set default client timeouts etc

	resp, err := http.Get("https://api.github.com/repos/" + Repository + "/releases/latest")
	if err != nil {
		fmt.Println(errors.New("Failed to check for new versions" + err.Error()))
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(errors.New("Failed to read response from version check" + err.Error()))
		return false
	}

	var v VersionResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		fmt.Println(errors.New("Failed to read response from version check" + err.Error()))
		return false
	}

	if len(v.TagName) == 0 {
		fmt.Println(errors.New("Failed to parse version from github response"))
		return false
	}

	if v.TagName != Version {
		fmt.Printf("\n** Precommit: New version %s available (installed %s) **\n\n", v.TagName, Version)
		return true
	}

	return false
}
