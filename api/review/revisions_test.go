package review

import (
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
)

const (
	original = `Line1
Line2
Line3
Line4`
	rev1 = `Line1
Line2.1
Line3.1
Line4
Line5`
	rev2 = `Line1
Line2
line3
Line4
Line5`

	diff0 = `--- test
+++ test
@@ -1,4 +1,4 @@
 Line1
 Line2
 Line3
 Line4`

	diff1 = `--- test
+++ test
@@ -1,4 +1,5 @@
 Line1
-Line2
-Line3
+Line2.1
+Line3.1
 Line4
+Line5`
	diff2 = `--- test
+++ test
@@ -1,4 +1,5 @@
 Line1
 Line2
-Line3
+line3
 Line4
+Line5`
)

var (
	originalLines = []int{0, 0, 0, 0}
	rev1Lines     = []int{0, 1, 1, 0, 1}
	rev2Lines     = []int{0, 2, 2, 0, 1}
)

func checkRevisions(t *testing.T, file File, revLines []int) {
	for i := 0; i < len(file.Lines); i++ {
		assert.Equal(t, file.Lines[i].Revision, revLines[i])
	}
}

func TestRevisionsCount(t *testing.T) {
	file := NewVersionedFile("test", difflib.SplitLines(original))
	assert.Equal(t, file.RevisionsCount(), 1)
	err := file.AddRevision(difflib.SplitLines(rev1))
	assert.Nil(t, err)
	assert.Equal(t, file.RevisionsCount(), 2)
	err = file.AddRevision(difflib.SplitLines(rev2))
	assert.Nil(t, err)
	assert.Equal(t, file.RevisionsCount(), 3)
}

func TestGetRevision(t *testing.T) {
	file := NewVersionedFile("test", difflib.SplitLines(original))
	file.AddRevision(difflib.SplitLines(rev1))
	file.AddRevision(difflib.SplitLines(rev2))

	o, err := file.GetRevision(0)
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(o.Content()), strings.TrimSpace(original))
	checkRevisions(t, o, originalLines)

	r1, err := file.GetRevision(1)
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(r1.Content()), strings.TrimSpace(rev1))
	checkRevisions(t, r1, rev1Lines)

	r2, err := file.GetRevision(2)
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(r2.Content()), strings.TrimSpace(rev2))
	checkRevisions(t, r2, rev2Lines)

	_, err = file.GetRevision(3)
	assert.NotNil(t, err)
}

func TestDiff(t *testing.T) {
	file := NewVersionedFile("test", difflib.SplitLines(original))
	file.AddRevision(difflib.SplitLines(rev1))
	file.AddRevision(difflib.SplitLines(rev2))

	diff, err := file.Diff(0, 0)
	assert.Nil(t, err)
	assert.Equal(t, diff.FileName, "test")
	assert.Equal(t, strings.TrimSpace(diff.String()), diff0)

	diff, err = file.Diff(0, 1)
	assert.Nil(t, err)
	assert.Equal(t, diff.FileName, "test")
	assert.Equal(t, strings.TrimSpace(diff.String()), diff1)

	diff, err = file.Diff(0, 2)
	assert.Nil(t, err)
	assert.Equal(t, diff.FileName, "test")
	assert.Equal(t, strings.TrimSpace(diff.String()), diff2)
}
