package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ONSdigital/git-diff-check/diffcheck"
)

const (
	// Return codes - a non-zero will cause the commit hook to reject the commit
	accepted = 0
	rejected = 1
)

func main() {

	fmt.Println("Running precommit diff check")

	patch, err := exec.Command("git", "diff", "-U0", "--staged").Output()
	if err != nil {
		log.Fatal("Failed to run git command:", err)
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
