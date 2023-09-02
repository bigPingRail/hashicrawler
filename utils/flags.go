package utils

import "flag"

var (
	Port    = flag.Int("p", 8080, "Listen port")
	Caching = flag.Bool("c", false, "Enable local caching")
	Auth    = flag.Bool("a", false, "Enable basic auth.\n By default login is 'admin' password 'testpass',\n you can override this values by setting up 'P_user' and 'P_password' environment variables")
)
