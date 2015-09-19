package multiconf_test

import (
	"fmt"
	"github.com/christopherobin/multiconf"
)

func ExampleNewMulticonf() {
	conf := multiconf.NewMulticonf(
		"myapp",              // app name
		multiconf.YamlParser, // a parser for my files
		// the actual configuration files
		"/etc/myapp.conf",
		"{{.Home}}/.myapp.conf",
		"{{.Conf}}/myapp.conf",
	)
	err := conf.Load()
	if err != nil {
		panic(err)
	}

	// access your configuration here
	if key, ok := conf.Root["key"]; ok {
		fmt.Println(key)
	}
}
