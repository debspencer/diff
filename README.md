# Diff - wrapper for `diff`

Package `diff` provides a wrapper around the `diff` program.

`import "github.com/debspencer/diff"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func DiffBuffer(file1 Buffer, file2 Buffer) ([]byte, error)](#DiffBuffer)
* [func DiffFile(file1, file2 string) ([]byte, error)](#DiffFile)
* [func DiffFileCommand(file1, file2, program, args string) ([]byte, error)](#DiffFileCommand)
* [type Buffer](#Buffer)


#### <a name="pkg-files">Package files</a>
[diff.go](/src/github.com/debspencer/diff/diff.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    DiffProgram     = "diff"
    DiffProgramArgs = "-u %f1 %f2"

    TempDir = ""
)
```


## <a name="DiffBuffer">func</a> [DiffBuffer](/src/target/diff.go?s=2203:2262#L84)
``` go
func DiffBuffer(file1 Buffer, file2 Buffer) ([]byte, error)
```
DiffBuffer will run a context diff between buf1 and buf2.  file1 and file2 are file names for the buffer.  Returns context diff.
Saves buffers to disk and the runs a context iff on them, replacing the file names



## <a name="DiffFile">func</a> [DiffFile](/src/target/diff.go?s=402:452#L26)
``` go
func DiffFile(file1, file2 string) ([]byte, error)
```
DiffFile will run a context diff between file1 and file2
If an error is returned diff failed, error strings is returned
if error is nil and no data is returned, then no differences



## <a name="DiffFileCommand">func</a> [DiffFileCommand](/src/target/diff.go?s=747:819#L33)
``` go
func DiffFileCommand(file1, file2, program, args string) ([]byte, error)
```
DiffFileCommand will run a context diff between file1 and file2 useing program and args
If an error is returned diff failed, error strings is returned
if error is nil and no data is returned, then no differences




## <a name="Buffer">type</a> [Buffer](/src/target/diff.go?s=1699:1983#L71)
``` go
type Buffer struct {
    // Data to be diffed
    Data []byte

    // Filename is an option name of the to assocation with diff.  If not specified, it will default to "file1" or "file2"
    Filename string
    // contains filtered or unexported fields
}

```
Buffer Defines a buffer to be diffed

