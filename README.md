rollingcsv
==========

A rolling CSV library for Go

##Usage

Almost identical to standard usage of the encoding/csv package:

```go
writer := rollingcsv.New("filename", ".", MAX_BYTES, MAX_LINES) 
var record []string
...
writer.Write(record)
writer.Close()
```
