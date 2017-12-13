package main

import (
	"flag"
	"fmt"
	"log"
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

var (
	target = flag.String("p", "", "(optional) path to repository")
)

func main() {
	flag.Parse()

	if *target == "" {
		*target = "."
	}
	fmt.Printf("Running precommit diff check on '%s'\n", *target)

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
