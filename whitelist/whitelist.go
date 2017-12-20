// Package whitelist deals with creating a whitelist of items that can be ignored
// by the diff checker for it's current scan.
//
// Expects a .whitelist file comprising [file] and [string] headings, each followed
// by a set of md5 hashes of the items to be whitelisted.
//
//   [files]
//	 550906dbdeec3751b4126e52fa57687f
//	 7e5a982291b1558b1e811870cced48ca
//   edce9423418a99dfc735ccce1138193d
//
//	 [strings]
//	 39cd0772d567a418464cd6b09ff5c912
//
package whitelist

import (
	"bufio"
	"fmt"
	"hash"
	"io"
)

type (
	// Signature is the unique identifier of the item. This is used to do the
	// identification of the object and ensure previously whitelisted items
	// (such as files) automatically become tainted if their content changes.
	Signature hash.Hash

	// WhiteList represents the whole set of items to be whitelisted
	WhiteList struct {
		Files   map[Signature]string
		Strings map[Signature]string
	}
)

// New creates a new whitelist from the given io.Reader stream
func New(r io.Reader) (*WhiteList, error) {
	b := bufio.NewReader(r)
	for {
		line, _, err := b.ReadLine() // TODO long lines
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// TODO
		fmt.Println(line) // REMOVE
	}

	return nil, nil
}

// HasSignature returns whether the given signature is in the current whitelist.
// If so, then success along with the type of item found is returned
func (w WhiteList) HasSignature(signature hash.Hash) (bool, string) {

	// Check files
	if _, ok := w.Files[signature]; ok {
		return true, "file"
	}

	// Check strings
	if _, ok := w.Strings[signature]; ok {
		return true, "string"
	}

	return false, ""
}
