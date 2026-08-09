package serverconfig

import (
	"os"
	"sync"
)

//-----------------------------------------------------------------------------
// Get the http-port from the environmet
//-----------------------------------------------------------------------------

var (
	httpPort     string
	httpPortOnce sync.Once
)

func GetHttpPort() string {

	httpPortOnce.Do(func() {

		if value := os.Getenv(httpPortKey); value != "" {

			httpPort = value
		} else {

			httpPort = httpPortDefault
		}
	})

	return httpPort
}

//-----------------------------------------------------------------------------
// Get the https-port from the environmet
//-----------------------------------------------------------------------------

var (
	httpsPort     string
	httpsPortOnce sync.Once
)

func GetHttpsPort() string {

	httpsPortOnce.Do(func() {

		if value := os.Getenv(httpsPortKey); value != "" {

			httpsPort = value
		} else {

			httpsPort = httpsPortDefault
		}
	})

	return httpsPort
}
