## ...---...
#### (SOS in Morse Code)

This package provides a Go error handling framework inspired by the following articles:
- [Failure is Your Domain](https://middlemost.com/failure-is-your-domain/)
- [Error Handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html)

Basic example.
```go
  1| package main
  2|
  3| import "github.com/bjaus/sos"
  4|
  5| func getByID(id int) (interface{}, error) {
  6|     return nil, fmt.Errorf("could not get record")
  7| }
  8|
  9| func main() {
 10|     id := 123
 11|
 12|     _, err := getByID(id)
 13|     if err != nil {
 14|         err = sos.New(sos.NOTFOUND).
 15|             WithError(err).
 16|             WithMessage("record not found: %d", id)
 17|     }
 18|
 19|     err = sos.Trace(err)
 20|     switch sos.Kind(err) {
 21|     case sos.NOTFOUND:
 22|         fmt.Println(err.Error())
 23|     default:
 24|         panic("...")
 25|     }
 26| }
```

```bash
$ go run .

could not get record
[not found] record not found: 123
    /path/to/file.go:14
    /path/to/file.go:19
```
