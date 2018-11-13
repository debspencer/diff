package diff

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	aBuf   = []byte("a\nb\nc\nd\ne\nf\ng\n")
	bBuf   = []byte("a\nb\nc\nf\ng\n")
	expect = "@@ -1,7 +1,5 @@\n a\n b\n c\n-d\n-e\n f\n g\n"
)

func TestDiffFile(t *testing.T) {
	a := assert.New(t)

	tmp1, err := createTmpFile(aBuf)
	a.NoError(err)

	tmp2, err := createTmpFile(bBuf)
	a.NoError(err)

	t.Run("Normal Diff", func(t *testing.T) {
		a := assert.New(t)

		diff, err := DiffFile(tmp1, tmp2)
		a.NoError(err)

		lines := bytes.SplitN(diff, []byte{'\n'}, 3)
		a.Len(lines, 3)

		a.True(bytes.HasPrefix(lines[0], []byte("--- ")))
		a.True(bytes.HasPrefix(lines[1], []byte("+++ ")))
		a.Equal(string(lines[2]), expect)
	})

	t.Run("No Differences", func(t *testing.T) {
		a := assert.New(t)

		diff, err := DiffFile(tmp1, tmp1)
		a.NoError(err)
		a.Equal(len(diff), 0)
	})

	t.Run("Diff Bad File", func(t *testing.T) {
		a := assert.New(t)

		diff, err := DiffFile(tmp1, "/this/file/does/not/exist")
		a.Error(err)
		a.NotEqual(len(diff), 0)
	})

	t.Run("Diff Bad Program", func(t *testing.T) {
		a := assert.New(t)

		// diff with stderr, but no error
		diff, err := DiffFileCommand(tmp1, tmp2, "dd", "count=0 if=%f1 of=%f2") // 'dd count=0' will write to stderr, but nothing to stdout
		a.Error(err)
		a.NotEqual(len(diff), 0)
	})

	t.Run("Diff alternate output", func(t *testing.T) {
		a := assert.New(t)

		// diff with stderr, but no error
		_, err := DiffFileCommand(tmp1, tmp2, "ls", "")
		a.NoError(err)
	})

}

func TestDiffBuffer(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		a := assert.New(t)

		diff, err := DiffBuffer(Buffer{Data: aBuf}, Buffer{Data: bBuf, Filename: "bfile"})
		a.NoError(err)

		lines := bytes.SplitN(diff, []byte{'\n'}, 3)
		a.Equal(len(lines), 3)

		a.True(bytes.HasPrefix(lines[0], []byte("--- ")))
		a.True(bytes.HasPrefix(lines[1], []byte("+++ ")))
		a.Equal(string(lines[2]), expect)
	})

	t.Run("Temp File Fail", func(t *testing.T) {
		a := assert.New(t)

		TempDir = "/dir/ectory/that/does/not/exist"
		_, err := DiffBuffer(Buffer{Data: aBuf}, Buffer{Data: bBuf, Filename: "bfile"})
		a.Error(err)

		TempDir = ""
	})

	t.Run("Write Fail", func(t *testing.T) {
		a := assert.New(t)

		writeData = writeDataFail
		_, err := DiffBuffer(Buffer{Data: aBuf}, Buffer{Data: bBuf, Filename: "bfile"})
		a.Error(err)

		writeData = fileWriteData
	})

	t.Run("Diff Fail", func(t *testing.T) {
		a := assert.New(t)

		DiffProgram = "bad/diff/command"
		_, err := DiffBuffer(Buffer{Data: aBuf}, Buffer{Data: bBuf})
		a.Error(err)
	})

	t.Run("Bad Output", func(t *testing.T) {
		a := assert.New(t)

		DiffProgram = "ls"

		_, err := DiffBuffer(Buffer{Data: aBuf}, Buffer{Data: bBuf})
		a.NoError(err)
	})
}
func writeDataFail(f *os.File, data []byte) (int, error) {
	return 0, fmt.Errorf("Write Failed")
}
