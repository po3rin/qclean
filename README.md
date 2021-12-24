# qclean

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

