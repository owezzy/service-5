// This program performs administrative tasks for the garage sale service.
package main

import (
	"github.com/owezzy/service-5/app/services/sales-api/v1/cmd"
	"github.com/owezzy/service-5/app/services/sales-api/v1/cmd/all"
	"github.com/owezzy/service-5/app/services/sales-api/v1/cmd/crud"
	"github.com/owezzy/service-5/app/services/sales-api/v1/cmd/reporting"
	"os"
)

var build = "develop"
var routes = "all" // go build -ldflags "-X main.routes=crud"

func main() {

	// The idea here is that we can build different versions of the binary
	// with different sets of exposed web APIs. By default we build a single
	// an instance with all the web APIs.
	//
	// Here is the scenario. It would be nice to build two binaries, one for the
	// transactional APIs (CRUD) and one for the reporting APIs. This would allow
	// the system to run two instances of the database. One instance tuned for the
	// transactional database calls and the other tuned for the reporting calls.
	// Tuning meaning indexing and memory requirements. The two databases can be
	// kept in sync with replication.

	switch routes {
	case "all":
		if err := cmd.Main(build, all.Routes()); err != nil {
			os.Exit(1)
		}

	case "crud":
		if err := cmd.Main(build, crud.Routes()); err != nil {
			os.Exit(1)
		}

	case "reporting":
		if err := cmd.Main(build, reporting.Routes()); err != nil {
			os.Exit(1)
		}
	}
}
