# Example

You can try to run this example like this.

```
go build ..
./atgen gen --templateDir=./template
```

This generates `v1_main_test.go` .

Run generated test code by this command.

```
go test
```


## Files of this example


- [`main.go`](main.go)
  - An example HTTP API Server. This is from https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo and I made a few changes.
- [`template/main_test.yaml`](template/main_test.yaml)
  - This defines requests/responses to test.
- [`template/template_test.go`](template/template_test.go)
  - This is a template of test code. See comments in this file how to write template code.
