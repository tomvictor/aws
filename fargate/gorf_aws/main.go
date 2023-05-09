package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting the Gorf")
	r := BootstrapRouter()
	err := r.Run(":80")
	if err != nil {
		log.Fatal("Unable to create the gin server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
