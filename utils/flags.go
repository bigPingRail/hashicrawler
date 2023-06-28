package utils

import "flag"

var (
	Port    = flag.Int("p", 8080, "listen port")
	Caching = flag.Bool("c", false, "enable local caching")
)
