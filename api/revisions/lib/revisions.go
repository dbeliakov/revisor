package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type modificationType int

const (
	insertModification modificationType = iota
	deleteModification modificationType = iota
)

// DiffType represents type of operation for particular line in diff
type DiffType int

const (
	// NoOperation - content similar in both files
	NoOperation DiffType = iota
	// InsertOperation - new inserted line
	InsertOperation DiffType = iota
	// DeleteOperation - line is removed from old file
	DeleteOperation DiffType = iota
)

// Line is a line in file with number of revision
type Line struct {
	Content  string `json:"content"`
	Revision int    `json:"revision"`
	ID       string `json:"id"`
}

// File represents lines of file with revision numbers
type File struct {
	Lines []Line
}

func (file File) String() string {
	var buffer bytes.Buffer
	for _, line := range file.Lines {
		buffer.WriteString(strconv.Itoa(line.Revision) + "\t")
		buffer.WriteString(line.Content)
	}
	return buffer.String()
}

// Content returns content of revisioned file
func (file File) Content() string {
	var buffer bytes.Buffer
	for _, line := range file.Lines {
		buffer.WriteString(line.Content)
	}
	return buffer.String()
}

// VersionedFile represents file with all it's revisions
type VersionedFile struct {
	Name      string
	Revisions []File
}

// NewVersionedFile constructs versioned file from content of original file
func NewVersionedFile(name string, content []string) VersionedFile {
	startRevision := File{make([]Line, 0, len(content))}
	for _, line := range content {
		u, err := uuid.NewV4()
		if err != nil {
			logrus.Panicf("Cannot create uuid for line: %+v", err)
		}
		startRevision.Lines = append(startRevision.Lines, Line{Content: line, Revision: 0, ID: u.String()})
	}
	file := VersionedFile{
		Name:      name,
		Revisions: []File{startRevision},
	}
	return file
}

// RevisionsCount returns count of revisions (includes original file)
func (file VersionedFile) RevisionsCount() int {
	return len(file.Revisions)
}

// GetRevision returns specified revision of file
func (file *VersionedFile) GetRevision(revision int) (File, error) {
	if revision >= len(file.Revisions) || revision < 0 {
		return File{}, errors.New(
			"Bad revision: expected from " + string(0) +
				" to " + string(len(file.Revisions)-1) + ", got " + string(revision))
	}
	return file.Revisions[revision], nil
}

// AddRevision adds new revision to the versioned file
func (file *VersionedFile) AddRevision(content []string) error {
	revisionsCount := file.RevisionsCount()
	lastRevision, err := file.GetRevision(revisionsCount - 1)
	if err != nil {
		return err
	}
	lastRevisionContent := make([]string, 0, len(lastRevision.Lines))
	for _, line := range lastRevision.Lines {
		lastRevisionContent = append(lastRevisionContent, line.Content)
	}
	newRevision := File{
		Lines: make([]Line, len(lastRevision.Lines)),
	}
	copy(newRevision.Lines, lastRevision.Lines)

	m := difflib.NewMatcher(lastRevisionContent, content)
	origNumChanges := 0
	for _, g := range m.GetGroupedOpCodes(0) {
		for _, c := range g {
			if c.Tag == 'e' {
				continue
			}
			if c.Tag == 'd' || c.Tag == 'r' {
				newRevision.Lines = append(
					newRevision.Lines[:c.I1+origNumChanges],
					newRevision.Lines[c.I2+origNumChanges:]...,
				)
				origNumChanges -= (c.I2 - c.I1)
			}
			if c.Tag == 'i' || c.Tag == 'r' {
				for j := c.J1; j < c.J2; j++ {
					u, err := uuid.NewV4()
					if err != nil {
						logrus.Panicf("Cannot create uuid: %+v", err)
					}
					newRevision.Lines = append(
						newRevision.Lines[:j],
						append(
							[]Line{Line{Content: content[j], Revision: revisionsCount, ID: u.String()}},
							newRevision.Lines[j:]...)...,
					)
				}
				origNumChanges += (c.J2 - c.J1)
			}
		}
	}
	file.Revisions = append(file.Revisions, newRevision)
	return nil
}

// DiffLine represents line in diff output
type DiffLine struct {
	Type DiffType `json:"-"`
	Old  *Line    `json:"old"`
	New  *Line    `json:"new"`
}

// MarshalJSON with corrected diff byte
func (diffLine DiffLine) MarshalJSON() ([]byte, error) {
	var diffType string
	if diffLine.Type == NoOperation {
		diffType = "no"
	} else if diffLine.Type == InsertOperation {
		diffType = "insert"
	} else if diffLine.Type == DeleteOperation {
		diffType = "delete"
	}

	type Alias DiffLine
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  diffType,
		Alias: (*Alias)(&diffLine),
	})
}

type diffRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// DiffGroup represents part of final diff
type DiffGroup struct {
	OldRange diffRange  `json:"old_range"`
	NewRange diffRange  `json:"new_range"`
	Lines    []DiffLine `json:"lines"`
}

// From https://github.com/pmezard/go-difflib/blob/master/difflib/difflib.go
// Convert range to the "ed" format
func formatRangeUnified(start, stop int) string {
	// Per the diff spec at http://www.unix.org/single_unix_specification/
	beginning := start + 1 // lines start numbering with one
	length := stop - start
	if length == 1 {
		return fmt.Sprintf("%d", beginning)
	}
	if length == 0 {
		beginning-- // empty ranges begin at line just before the range
	}
	return fmt.Sprintf("%d,%d", beginning, length)
}

func (group DiffGroup) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("@@ -%s +%s @@\n",
		formatRangeUnified(group.OldRange.From, group.OldRange.To),
		formatRangeUnified(group.NewRange.From, group.NewRange.To)))
	for _, line := range group.Lines {
		if line.Type == NoOperation {
			buffer.WriteString(" ")
			buffer.WriteString(line.Old.Content)
		} else if line.Type == DeleteOperation {
			buffer.WriteString("-")
			buffer.WriteString(line.Old.Content)
		} else if line.Type == InsertOperation {
			buffer.WriteString("+")
			buffer.WriteString(line.New.Content)
		}
	}
	return buffer.String()
}

// Diff represents diff between revisions
type Diff struct {
	FileName string      `json:"filename"`
	Groups   []DiffGroup `json:"groups"`
}

// Diff returns diff between two revisions as string
func (file *VersionedFile) Diff(revision1, revision2 int) (Diff, error) {
	file1, err := file.GetRevision(revision1)
	if err != nil {
		return Diff{}, err
	}
	file2, err := file.GetRevision(revision2)
	if err != nil {
		return Diff{}, err
	}
	file1Content, file2Content :=
		difflib.SplitLines(strings.TrimSpace(file1.Content())),
		difflib.SplitLines(strings.TrimSpace(file2.Content()))
	m := difflib.NewMatcher(file1Content, file2Content)
	diff := Diff{
		FileName: file.Name,
	}
	if revision1 == revision2 {
		// Return content of file
		diff.Groups = append(diff.Groups, DiffGroup{
			OldRange: diffRange{
				From: 0,
				To:   len(file1Content),
			},
			NewRange: diffRange{
				From: 0,
				To:   len(file1Content),
			},
		})
		for i := range file1Content {
			diff.Groups[0].Lines = append(diff.Groups[0].Lines, DiffLine{
				Type: NoOperation,
				Old:  &file1.Lines[i],
				New:  &file2.Lines[i],
			})
		}
		return diff, nil
	}
	for _, g := range m.GetGroupedOpCodes(10000) {
		first, last := g[0], g[len(g)-1]
		group := DiffGroup{
			OldRange: diffRange{
				From: first.I1,
				To:   last.I2,
			},
			NewRange: diffRange{
				From: first.J1,
				To:   last.J2,
			},
		}
		for _, c := range g {
			if c.Tag == 'e' {
				for i := 0; i < c.I2-c.I1; i++ {
					group.Lines = append(group.Lines, DiffLine{
						Type: NoOperation,
						Old:  &file1.Lines[c.I1+i],
						New:  &file2.Lines[c.J1+i],
					})
				}
			}
			if c.Tag == 'd' || c.Tag == 'r' {
				for i := 0; i < c.I2-c.I1; i++ {
					group.Lines = append(group.Lines, DiffLine{
						Type: DeleteOperation,
						Old:  &file1.Lines[c.I1+i],
					})
				}
			}
			if c.Tag == 'i' || c.Tag == 'r' {
				for j := 0; j < c.J2-c.J1; j++ {
					group.Lines = append(group.Lines, DiffLine{
						Type: InsertOperation,
						New:  &file2.Lines[c.J1+j],
					})
				}
			}
		}
		diff.Groups = append(diff.Groups, group)
	}
	return diff, nil
}

func (diff Diff) String() string {
	if len(diff.Groups) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("--- %s\n", diff.FileName))
	buffer.WriteString(fmt.Sprintf("+++ %s\n", diff.FileName))
	for _, g := range diff.Groups {
		buffer.WriteString(g.String())
	}
	return buffer.String()
}
