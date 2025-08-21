# buffers-pool
[![Code Analysis](https://github.com/sv-tools/buffers-pool/actions/workflows/checks.yaml/badge.svg)](https://github.com/sv-tools/buffers-pool/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/buffers-pool.svg)](https://pkg.go.dev/github.com/sv-tools/buffers-pool)
[![codecov](https://codecov.io/gh/sv-tools/buffers-pool/branch/main/graph/badge.svg?token=0XVOTDR1CW)](https://codecov.io/gh/sv-tools/buffers-pool)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/buffers-pool?style=flat)](https://github.com/sv-tools/buffers-pool/releases)

[!IMPORTANT]
Archived becasue the package is stable and the new releases are not expected.
Feel free to copy the implementation directly to your code, instead of using it as a dependency.

The small library with an implementation of Buffer Pool.
The library was created to avoid repeating this code.

Here is a good article how to implement and properly use the Buffer Pools: https://www.captaincodeman.com/2017/06/02/golang-buffer-pool-gotcha

Check the [tests](pool_test.go) file for some examples.

## Usage

```go
package main

import (
	"fmt"
	"sync"
	"text/template"

	buffferspool "github.com/sv-tools/buffers-pool"
)

func render(tmpl string, data interface{}) (string, error) {
	tp, err := template.New("test").Parse(tmpl)
	if err != nil {
		return "", err
	}

	buf := buffferspool.Get()
	defer buffferspool.Put(buf)

	if err := tp.Execute(buf, data); err != nil {
		return "", err
	}

	// the usage of buf.String is safe 
	return buf.String(), nil
}

func main() {
    var tmpl string
    var data []interface{}
    // ... load template and data to variables ...

    var wg sync.WaitGroup
    res := make(chan string, len(data))
    for _, d := range data {
        wg.Add(1)
        go func(data interface{}) {
            defer wg.Done()
            val, err := render(tmpl, data)
            if err != nil {
                res <- err.Error()
                return
            }

            res <- val
        }(d)
    }

    wg.Wait()
    close(res)

    for val := range res {
    	fmt.Println(val)

    }
}
```

## License

MIT licensed. See the bundled [LICENSE](LICENSE) file for more details.
