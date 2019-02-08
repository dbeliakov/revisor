package review

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
)

const (
	fileName = "main.cpp"
)

func readFile(name string) string {
	res, err := ioutil.ReadFile(path.Join("test_data", name))
	if err != nil {
		panic(errors.Wrap(err, "Cannot open test data file"))
	}
	return string(res)
}

var (
	revisions = []string{readFile("first.cpp"), readFile("second.cpp"), readFile("third.cpp")}
)

func TestRevisionsCount(t *testing.T) {
	file := NewVersionedFile(fileName, difflib.SplitLines(revisions[0]))
	assert.Equal(t, file.RevisionsCount(), 1)
	err := file.AddRevision(difflib.SplitLines(revisions[1]))
	assert.Nil(t, err)
	assert.Equal(t, file.RevisionsCount(), 2)
	err = file.AddRevision(difflib.SplitLines(revisions[2]))
	assert.Nil(t, err)
	assert.Equal(t, file.RevisionsCount(), 3)
}

func TestRevisionContent(t *testing.T) {
	file := NewVersionedFile(fileName, difflib.SplitLines(revisions[0]))
	checkRevisionsContent := func() {
		for i := 0; i < file.RevisionsCount(); i++ {
			content, err := file.GetRevision(i)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, strings.TrimSpace(content.Content()), strings.TrimSpace(revisions[i]))
		}
	}
	checkRevisionsContent()
	file.AddRevision(difflib.SplitLines(revisions[1]))
	checkRevisionsContent()
	file.AddRevision(difflib.SplitLines(revisions[2]))
	checkRevisionsContent()
}

func TestIncorrectRevisionNumber(t *testing.T) {
	file := NewVersionedFile(fileName, difflib.SplitLines(revisions[0]))
	_, err := file.GetRevision(0)
	assert.NoError(t, err)
	_, err = file.GetRevision(1)
	assert.Error(t, err)
	_, err = file.GetRevision(2)
	assert.Error(t, err)
}

func TestLineRevisionNumbers(t *testing.T) {
	lineRevisions := [][]int{
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 2, 2, 1, 2, 2, 2, 0, 2, 0, 0, 0, 0, 0, 0, 0},
	}
	file := NewVersionedFile(fileName, difflib.SplitLines(revisions[0]))
	checkLineRevisions := func() {
		rev, err := file.GetRevision(file.RevisionsCount() - 1)
		assert.NoError(t, err)
		expected := lineRevisions[file.RevisionsCount()-1]
		if !assert.Equal(t, len(expected), len(rev.Lines)) {
			return
		}
		for i := range rev.Lines {
			assert.Equal(t, expected[i], rev.Lines[i].Revision, fmt.Sprintln(rev.Lines))
		}
	}
	checkLineRevisions()
	err := file.AddRevision(difflib.SplitLines(revisions[1]))
	assert.NoError(t, err)
	checkLineRevisions()
	err = file.AddRevision(difflib.SplitLines(revisions[2]))
	assert.NoError(t, err)
	checkLineRevisions()
}

func TestDiff(t *testing.T) {
	diffs := map[int]map[int]string{
		0: {
			0: readFile("first-first.diff"),
			1: readFile("first-second.diff"),
			2: readFile("first-third.diff"),
		},
		1: {
			0: readFile("second-first.diff"),
			1: readFile("second-second.diff"),
			2: readFile("second-third.diff"),
		},
		2: {
			0: readFile("third-first.diff"),
			1: readFile("third-second.diff"),
			2: readFile("third-third.diff"),
		},
	}
	file := NewVersionedFile(fileName, difflib.SplitLines(revisions[0]))
	checkDiffs := func() {
		for i := 0; i < file.RevisionsCount(); i++ {
			for j := 0; j < file.RevisionsCount(); j++ {
				diff, err := file.Diff(i, j)
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, strings.TrimSpace(diffs[j][i]), strings.TrimSpace(diff.String()))
			}
		}
	}
	file.AddRevision(difflib.SplitLines(revisions[1]))
	checkDiffs()
	file.AddRevision(difflib.SplitLines(revisions[2]))
	checkDiffs()
}
