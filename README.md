# Multiconf

Load multiple configuration files on top of each other in a certain order, eg:

```
/etc/myapp.conf
~/.myapp.conf
~/.local/share/myapp/myapp.conf
```

## Installation

```
go get -u github.com/christopherobin/multiconf
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/christopherobin/multiconf"
)

// see the documentation for the list of variables available
var conf = multiconf.NewMulticonf(
	"myapp",              // app name
	multiconf.YamlParser, // a parser for my files
	// the actual configuration files
	"/etc/myapp.conf",
	"{{.Home}}/.myapp.conf",
	"{{.Conf}}/myapp.conf",
)
```

[API Documentation](http://godoc.org/github.com/christopherobin/multiconf)