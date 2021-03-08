# Go-inject

## Flags

```sh
Usage:
  go-inject [flags]

Flags:
      --dir string       the working directory (default ".")
      --func string      the name of function that need to be injected
  -h, --help             help for go-inject
      --imports string   import paths used by injected go file
      --instr string     the golang snippet
      --method string    the name of method that need to be injected
      --pkg string       the name of package that need to be injected
```

## Demo

[1]

```go
package example

func A() {

}
```

```sh
~ make build
~ ./go-inject --dir . --func A --imports "f123 fmt" --pkg "github.com/gowins/go-inject/example" --instr 'f123.Println("Hello World")'
```

You will get to see

```go
package example

import (
	f123 "fmt"
)

func A() {
	f123.Println("Hello World")
}
```

[2]

```sh
package example

type B struct{}

func (B) Hi() {

}
```

```sh
~ make build
~ ./go-inject --dir . --method "B Hello" --imports "fmt,_ net/http/pprof" --pkg "github.com/gowins/go-inject/example" --instr 'fmt.Println("Hello World")'
```

You will get to see

```go
package example

import (
	"fmt"
	_ "net/http/pprof"
)

type B struct{}

func (*B) Hello() {
	fmt.Println("Hello World")
}
```

[3]

```sh
package example

type B struct{}

func (B) Hi() {

}
```

```sh
~ make build
~ ./go-inject --dir . --method "B Hi" --imports "f1234 fmt,_ net/http/pprof" --pkg "github.com/gowins/go-inject/example" --instr 'f1234.Println("Hello World")'
```

You will get to see

```go
package example

import (
	f1234 "fmt"
	_ "net/http/pprof"
)

type B struct{}

func (B) Hi() {
	f1234.Println("Hello World")
}
```
