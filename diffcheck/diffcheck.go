// Package diffcheck provides functions for checking a git diff for potentially
// sensitive information.
package diffcheck

import (
	"bufio"
	"bytes"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ONSdigital/git-diff-check/entropy"
	"github.com/ONSdigital/git-diff-check/rule"
)

type (
	// Warning is a specific warning about a file in diff. One or more are compiled
	// into a `Report`
	Warning struct {
		// The ruleset type that triggered the warning. e.g. "file" or "line"
		Type string

		// Human compatible warning description
		Description string

		// Line number (if applicable) where the warning was triggered.
		// If no line then will be -1
		Line int
	}

	// Report is a collection of warnings for a particular file discovered in
	// a patch
	Report struct {
		// Current relative path of the file to which the report pertains
		Path string

		// Old path of the file - will be identical unless the file has been
		// moved/renamed as part of the changeset
		OldPath string

		// Set of warnings pertaining to this report
		Warnings []Warning
	}
)

var (
	// Matches the first offset in the old and new diff
	reOffset = regexp.MustCompile("^@@ -(\\d+).* \\+(\\d+).* @@")
)

const (
	noNewlineWarning = "\\ No newline at end of file"
)

// SnoopPatch takes a raw github patch byte array and tests it against the
// defined rulesets. Returns true if diff appears clean and false otherwise. In
// the case of a potentially unclean diff, a report set will also be returned
// detailing a set of warnings identified.
func SnoopPatch(patch []byte) (bool, []Report, error) {

	reader := bufio.NewReader(bytes.NewReader(patch))

	reports := []Report{}
	report := Report{}

	inHunk := false
	linePosition := 0

	var tmp []byte

	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix {
			// Long line. Temporarily store what we have and pick up the next
			// chunk on the next pass.
			if tmp == nil {
				tmp = make([]byte, len(line))
			}
			tmp = append(tmp, line...)
			continue
		}
		if err == io.EOF {
			if len(report.Warnings) > 0 {
				reports = append(reports, report)
			}
			break
		}
		if tmp != nil {
			// If we have temporarily stored long line data then dump it back out to the
			// line and continue
			line = append(tmp, line...)
			tmp = nil
		}

		// Check whether we're starting a new file block section of the patch or
		// if we've hit the end of the input.
		if bytes.HasPrefix(line, []byte("diff --git")) {
			inHunk = false

			// If we already have previous warnings then we output the existing
			// report and clear down.
			if len(report.Warnings) > 0 {
				reports = append(reports, report)
				report = Report{}
			}

			report.Path, report.OldPath = getFilePath(line)

			if ok, w := checkFile(report.Path); !ok {
				report.Warnings = append(report.Warnings, w...)
			}

			continue
		}

		// Check whether we're starting a new hunk
		if bytes.HasPrefix(line, []byte("@@ ")) {
			matches := reOffset.FindAllSubmatch(line, -1)

			// nb. If there are no matches then we just assume we're in a line
			// that _looks_ like a hunk start but isn't (This isn't infallible
			// but only works for line number convenience anyhow - actual scan
			// isn't reliant on it)
			if len(matches) > 0 {
				linePosition, err = strconv.Atoi(string(matches[0][2]))
				if err != nil {
					// TODO handle better!
					return false, nil, err
				}
			}

			inHunk = true
			continue
		}

		if inHunk && string(line) != noNewlineWarning {
			if ok, w := checkLineBytes(line, linePosition); !ok {
				report.Warnings = append(report.Warnings, w...)
			}
		}

		// Incremented at the end as the hunk for an offset begins with the index
		// for the line. Post incrementing here ensures the line read in the
		// hunk is identified with it's correct file position.
		linePosition++

	}

	if len(reports) > 0 {
		return false, reports, nil
	}

	// All ok!
	return true, nil, nil
}

// checkLineBytes runs rules against each line in the patch for a file to see whether
// they match potentially sensitive patterns
// Returns true with a set of Warning structs if found, otherwise false
func checkLineBytes(line []byte, position int) (bool, []Warning) {

	warnings := []Warning{}

	// Normal line rulesets
	for _, rule := range rule.Sets["line"] {
		if rule.Regex.Match(line) {
			warnings = append(warnings, Warning{Type: "line", Description: rule.Caption, Line: position})
		}
	}

	// Entropy check
	if ok, _ := entropy.Check(line); !ok {
		warnings = append(warnings, Warning{Type: "line", Description: "Possible key in high entropy string", Line: position})
	}

	if len(warnings) > 0 {
		return false, warnings
	}
	return false, nil
}

// checkFile runs gitrob rules against the file name to see whether they match
// Returns true with a set of Warning structs if found, otherwise false
func checkFile(path string) (bool, []Warning) {

	// Prep the three possible bits that could be examined:
	// - path (already have), just filename, and extension
	name := filepath.Base(path)

	// Ext returns with a prefix period whilst gitrob rules specify without
	// so we need to try as well
	extension := strings.TrimLeft(filepath.Ext(path), ".")

	warnings := []Warning{}

	for _, rule := range rule.Sets["file"] {
		// Determine which bit of the filename we want to test
		toTest := ""

		switch rule.Part {
		case "extension":
			toTest = extension
		case "path":
			toTest = path
		case "filename":
			toTest = name
		}

		switch rule.Type {
		case "regex":
			if rule.Regex.Match([]byte(toTest)) {
				warnings = append(warnings, Warning{Type: "file", Description: rule.Caption, Line: -1})
			}
		case "match":
			if rule.Pattern == toTest {
				warnings = append(warnings, Warning{Type: "file", Description: rule.Caption, Line: -1})
			}
		}
	}
	if len(warnings) > 0 {
		return false, warnings
	}
	return true, nil
}

// Returns the actual filename and previous filename (which may or may not
// be different)
func getFilePath(raw []byte) (string, string) {
	words := bytes.Split(raw, []byte(" "))

	new := string(bytes.TrimLeft(words[3], "b/"))
	old := string(bytes.TrimLeft(words[2], "a/"))

	return new, old
}
