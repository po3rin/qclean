# qclean

<img src="https://img.shields.io/badge/go-v1.17-blue.svg"/> [![GoDoc](https://godoc.org/github.com/po3rin/qclean?status.svg)](https://godoc.org/github.com/po3rin/qclean) ![Go Test](https://github.com/po3rin/qclean/workflows/Go%20Test/badge.svg) 

qclean lets you to clean up search query in japanese. This is mainly used to remove wasted space.

## Quick Start

```go
package main

var cleaner *qclean.Cleaner

// It takes time to read the dictionary, so it is recommended to initialize Cleaner in the init function.
func init() {
	cleaner, _ = qclean.NewCleaner()
}

func main() {
	result, _ := cleaner.Clean("苔 癬 治っ た")
	fmt.Println(result) // 苔癬 治った
}
```

