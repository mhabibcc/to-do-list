//   Product Api: 
//    version: 0.1 
//    title: To Do List API
//   Schemes: http, https 
//   Host: localhost:3000
//   BasePath: /api/
//      Consumes: 
//      - application/json
//   Produces: 
//   - application/json
//   swagger:meta
package main

import (
	"log"
	"to-do-list/internal/config"
)

const repoName = "to-do-list"

func main() {
	Config, err := config.New(repoName)
	if err != nil {
		panic(err)
	}
	log.Fatal(startApp(Config))
}
