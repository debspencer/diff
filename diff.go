package diff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	DiffProgram     = "diff"
	DiffProgramArgs = "-u %f1 %f2"

	TempDir = ""

	writeData = fileWriteData
)

// DiffFile will run a context diff between file1 and file2
// If an error is returned diff failed, error strings is returned
// if error is nil and no data is returned, then no differences
func DiffFile(file1, file2 string) ([]byte, error) {
	return DiffFileCommand(file1, file2, DiffProgram, DiffProgramArgs)
}

// DiffFileCommand will run a context diff between file1 and file2 useing program and args
// If an error is returned diff failed, error strings is returned
// if error is nil and no data is returned, then no differences
func DiffFileCommand(file1, file2, program, args string) ([]byte, error) {
	// Diff will return
	// 0 for no diffs
	// 1 for diffs
	// 2 for error
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if strings.Contains(args, "%f") {
		args = strings.Replace(args, "%f1", file1, -1)
		args = strings.Replace(args, "%f2", file2, -1)
	} else {
		args += " " + file1 + " " + file2
	}
	args = strings.TrimSpace(args)

	cmd := exec.Command(program, strings.Fields(args)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	// There was standard error, something did not go right
	if stderr.Len() > 0 {
		if err == nil {
			err = fmt.Errorf("%s: failed", DiffProgram)
		}
		return stderr.Bytes(), err
	}

	// If there was no output, then there should be no error
	if stdout.Len() == 0 && err != nil {
		return []byte(err.Error()), fmt.Errorf("%s: failed", DiffProgram)
	}

	return stdout.Bytes(), nil
}

// Buffer Defines a buffer to be diffed
type Buffer struct {
	// Data to be diffed
	Data []byte

	// Filename is an option name of the to assocation with diff.  If not specified, it will default to "file1" or "file2"
	Filename string

	temp string // temp filename assigned to file for diffing
	name string // name of file
}

// DiffBuffer will run a context diff between buf1 and buf2.  file1 and file2 are file names for the buffer.  Returns context diff.
// Saves buffers to disk and the runs a context iff on them, replacing the file names
func DiffBuffer(file1 Buffer, file2 Buffer) ([]byte, error) {
	files := [2]Buffer{
		file1,
		file2,
	}

	for i := range files {
		f := &files[i]
		temp, err := createTmpFile(f.Data)
		if err != nil {
			return nil, err
		}
		f.temp = temp
		f.name = f.Filename
		if f.name == "" {
			f.name = "file" + strconv.Itoa(i+1)
		}
		//		defer os.Remove(temp)
	}

	diff, err := DiffFile(files[0].temp, files[1].temp)
	if err != nil {
		return nil, err
	}

	// if no data, no differnces
	if len(diff) > 0 {
		// replace the temp file names with the ones
		diff, err = fixFilenames(diff, files[0].name, files[1].name)
	}
	return diff, err
}

func createTmpFile(data []byte) (string, error) {
	tmp, err := ioutil.TempFile(TempDir, "diff")
	if err != nil {
		return "", err
	}
	file := tmp.Name()

	_, err = writeData(tmp, data)
	closeErr := tmp.Close()
	if err == nil {
		err = closeErr
	}

	if err != nil {
		os.Remove(file)
		return "", err
	}
	return file, nil
}

// fixFilenames will replace the temp file names with requested ones
func fixFilenames(diff []byte, filename1 string, filename2 string) ([]byte, error) {
	// We only need to replace the first three lines, `
	lines := bytes.SplitN(diff, []byte{'\n'}, 3)
	if len(lines) >= 3 && bytes.HasPrefix(lines[0], []byte("--- ")) && bytes.HasPrefix(lines[1], []byte("+++ ")) {
		// unified diff

		now := time.Now()
		ts := now.Format("2006-01-02 15:04:05.000000000 -0700")

		lines[0] = []byte(fmt.Sprintf("--- %s\t%s", filename1, ts))
		lines[1] = []byte(fmt.Sprintf("+++ %s\t%s", filename2, ts))
	}
	return bytes.Join(lines, []byte{'\n'}), nil
}

func fileWriteData(f *os.File, data []byte) (int, error) {
	return f.Write(data)
}
