# Test protoc plugin

The test plugin adds to proto message new method which returns json string for types `bool`, `int32` and `string`
The option can be added to proto message like:
```protobuf
...
import "github.com/YReshetko/protoc-gen-marshaler/proto/marshaler.proto";
...
message TestMessage {
  option (marshaler.enable) = true;
  string StrField = 1;
  bool BoolField = 2;
  int32 IntField = 3;
}
```

in the result you will get `<proto_file_name>.marshaler.pb.go` with next method:
```go

func (m *TestMessage) CustomMarshal() string {
	out := "{"
	out += "\"StrField\": \"" + m.StrField + "\""
	out += ","
	out += "\"BoolField\":" + strconv.FormatBool(m.BoolField)
	out += ","
	out += "\"IntField\":" + strconv.FormatInt(int64(m.IntField), 10)
	out += "}"
	return out
}
```

### Install

```go get -u github.com/YReshetko/protoc-gen-marshaler```

### Usage

```protoc -I $(GOPATH)/src:. --go_out=paths=source_relative:. --marshaler_out=paths=source_relative:. proto/test.proto```