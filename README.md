# UUID [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/ryanfowler/uuid) [![Go Report Card](https://goreportcard.com/badge/github.com/ryanfowler/uuid)](https://goreportcard.com/report/github.com/ryanfowler/uuid)

UUID provides functions for generating and formatting UUIDs according to RFC 4122.

## Sample Usage

### Version 3

Use [v5](#version-5), if possible (unless for legacy reasons).

### Version 4

To generate a new v4 UUID:
```go
package main

import (
	"fmt"

	"github.com/ryanfowler/uuid"
)

func main() {
	u := uuid.Must(uuid.NewV4())
	fmt.Println(u.String())
}
```
will output something like: `9e754ef6-8dd9-4903-af43-7aea99bfb1fe`.

NewV4 will only return an error when it is unable to read random bytes from the OS.
Otherwise, the returned UUID will be made up of random bytes with the appropriate variant and version bits set.

### Version 5

To generate a new v5 UUID:
```go
package main

import (
	"fmt"

	"github.com/ryanfowler/uuid"
)

func main() {
	namespace := "9e754ef6-8dd9-5903-af43-7aea99bfb1fe"
	u := uuid.NewV5(namespace, []byte("unique bytes"))
	fmt.Println(u.String())
}
```
will output something like: `e83915a5-a2c6-573b-a5f7-8cf2badd0af5`.

UUIDs generated from NewV5 with the same namespace & name will be equal every time NewV5 is called.
In other words, `same input -> same output`.

### Version 7

To generate a new v7 UUID:

```go
package main

import (
	"fmt"

	"github.com/ryanfowler/uuid"
)

func main() {
	u := uuid.Must(uuid.NewV7(time.Now()))
	fmt.Println(u.String())
}
```

To get the timestamp out of a v7 UUID, you can use the following method:

```go
u := uuid.Must(uuid.NewV7(time.Now()))
timestamp, ok := u.Time()
if ok {
	fmt.Println(timestamp)
}
```

### Formatting

A UUID represents a 16 byte array (128 bits).
In order to use the UUID in a human-readable form, either `Format`, `Bytes`, or `String` should be used.

Format will return the UUID bytes as 36 byte array in the format:
```
xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
where `x`'s are hexadecimal characters.

To format the UUID into a 36 byte slice, use `Bytes`.

To format the UUID as a string, use `String`.

### Parsing

A UUID can be parsed from the following formats:

- A 16-byte raw UUID.
- A 32-byte hexadecimal UUID without dashes e.g. 9e754ef68dd94903af437aea99bfb1fe
- A 36-byte hexadecimal UUID with dashes e.g. 9e754ef6-8dd9-4903-af43-7aea99bfb1fe

Example:

```go
raw := "9e754ef6-8dd9-4903-af43-7aea99bfb1fe"
u := uuid.Must(uuid.ParseString(raw))
fmt.Println(u.String())
```

will ouput: `9e754ef6-8dd9-4903-af43-7aea99bfb1fe`.

## License

The MIT License.

See the LICENSE file for more information.
