# UUID [![GoDoc](https://godoc.org/github.com/ryanfowler/uuid?status.svg)](https://godoc.org/github.com/ryanfowler/uuid) [![Go Report Card](https://goreportcard.com/badge/github.com/ryanfowler/uuid)](https://goreportcard.com/report/github.com/ryanfowler/uuid)

UUID provides functions for generating and formatting UUIDs according to RFC 4122.

## Sample Usage

#### Version 3

Don't use v3 - use [v5](#version-5) (unless for legacy reasons).

#### Version 4

To generate a new v4 UUID:
```go
package main

import (
	"fmt"

	"github.com/ryanfowler/uuid"
)

func main() {
	u, err := uuid.NewV4()
	if err != nil {
		// unable to read random bytes, bad!!!
		return
	}
	fmt.Println(u.FormatString())
}
```
will output something like: `9e754ef6-8dd9-4903-af43-7aea99bfb1fe`.

NewV4 will only return an error when it is unable to read random bytes from the OS.
Otherwise, the returned UUID will be made up of random bytes with the appropriate variant and version bits set.

#### Version 5

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
	fmt.Println(u.FormatString())
}
```
will output something like: `e83915a5-a2c6-573b-a5f7-8cf2badd0af5`.

UUIDs generated from NewV5 with the same namespace & name will be equal every time NewV5 is called.
In other words, `same input -> same output`.

#### Formatting

A UUID represents a 16 byte array (128 bits).
In order to use the UUID in a human-readable form, either `Format` or `FormatString` should be used.

Format will returns the UUID bytes in the format:
```
xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
where `x`'s are hexadecimal characters.

FormatString will return the string representation of the above UUID format.
The returned UUID will have a length of 36 bytes.

## License

The MIT License.

See the LICENSE file for more information.
