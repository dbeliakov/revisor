package revisions

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	uuid "github.com/satori/go.uuid"
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
	Content  string
	Revision int
	ID       string
}

// Modification represents changes to get next revision of file
type Modification struct {
	Type       modificationType
	LineNumber int
	Content    string
}

// Patch is a set of modifications to get next revision
type Patch struct {
	Modifications []Modification
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
	Name     string
	Original File
	Patches  []Patch
}

// NewVersionedFile constructs versioned file from content of original file
func NewVersionedFile(name string, content []string) VersionedFile {
	file := VersionedFile{
		Name:     name,
		Original: File{make([]Line, 0, len(content))},
	}
	for _, line := range content {
		u, err := uuid.NewV4()
		if err != nil {
			panic(errors.New("Assert: " + err.Error()))
		}
		file.Original.Lines = append(file.Original.Lines, Line{Content: line, Revision: 0, ID: u.String()})
	}
	return file
}

// RevisionsCount returns count of revisions (includes original file)
func (file VersionedFile) RevisionsCount() int {
	return len(file.Patches) + 1
}

// GetRevision returns specified revision of file
func (file *VersionedFile) GetRevision(revision int) (File, error) {
	if revision > len(file.Patches) || revision < 0 {
		return File{}, errors.New("Bad revision")
	}
	result := File{
		Lines: make([]Line, len(file.Original.Lines)),
	}
	copy(result.Lines, file.Original.Lines)
	for i := 1; i <= revision; i++ {
		patch := file.Patches[i-1]
		removed := 0
		for _, mod := range patch.Modifications {
			if mod.Type != deleteModification {
				continue
			}
			result.Lines = append(
				result.Lines[:mod.LineNumber-removed],
				result.Lines[mod.LineNumber+1-removed:]...,
			)
			removed++
		}
		for _, mod := range patch.Modifications {
			if mod.Type != insertModification {
				continue
			}
			result.Lines = append(
				result.Lines[:mod.LineNumber+1],
				append(
					[]Line{Line{Content: mod.Content, Revision: i}},
					result.Lines[mod.LineNumber+1:]...)...,
			)
		}
	}
	return result, nil
}

// AddRevision adds new revision to the versioned file
func (file *VersionedFile) AddRevision(content []string) error {
	lastRevision, err := file.GetRevision(len(file.Patches))
	if err != nil {
		return err
	}
	lastRevisionContent := make([]string, 0, len(lastRevision.Lines))
	for _, line := range lastRevision.Lines {
		lastRevisionContent = append(lastRevisionContent, line.Content)
	}
	m := difflib.NewMatcher(lastRevisionContent, content)
	var patch Patch
	for _, g := range m.GetGroupedOpCodes(0) {
		for _, c := range g {
			if c.Tag == 'e' {
				continue
			}
			if c.Tag == 'd' || c.Tag == 'r' {
				for i := c.I1; i < c.I2; i++ {
					patch.Modifications = append(patch.Modifications, Modification{
						Type:       deleteModification,
						LineNumber: i,
					})
				}
			}
			if c.Tag == 'i' || c.Tag == 'r' {
				for j := c.J1; j < c.J2; j++ {
					patch.Modifications = append(patch.Modifications, Modification{
						Type:       insertModification,
						LineNumber: j - 1,
						Content:    content[j],
					})
				}
			}
		}
	}
	file.Patches = append(file.Patches, patch)
	return nil
}

// DiffLine represents line in diff output
type DiffLine struct {
	Type DiffType
	Old  *Line
	New  *Line
}

type diffRange struct {
	From int
	To   int
}

// DiffGroup represents part of final diff
type DiffGroup struct {
	OldRange diffRange
	NewRange diffRange
	Lines    []DiffLine
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
	FileName string
	Groups   []DiffGroup
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
