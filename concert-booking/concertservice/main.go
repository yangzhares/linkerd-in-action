package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var f string

func init() {
	flag.StringVar(&f, "config", "config.json", "concert service config file")
	flag.StringVar(&f, "c", "config.json", "concert service config file")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	var help = `
   -c= (--config=)       Concert service config file(default: config.json)
`

	fmt.Fprint(os.Stderr, help)
}

func main() {
	conf, err := InitConfig(f)
	if err != nil {
		log.Println(err)
		return
	}

	log.Fatal(ServeAPI(conf))
}
