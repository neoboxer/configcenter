package main

import (
	"github.com/neoboxer/configcenter/internal"
)

func main() {
	// router new
	router := internal.NewRouter()
	// router register http handlers
	// router run
	router.Run(":8080")
}
