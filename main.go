package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	appModeFlag := flag.String("mode", appModeClient, "start in client or server mode")
	userHandleFlag := flag.String("user", "", "Set user handle through cli")
	verboseFlag := flag.Bool("v",false,"turn on verbose mode on")
	flag.Parse()
	if *verboseFlag {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}else{
		log.SetOutput(ioutil.Discard)
	}
	log.Println("test")
	switch strings.ToLower(*appModeFlag) {
	case appModeServer:
		startServerMode()
	case appModeClient:
		startClientMode(*userHandleFlag)
	case appModeTest:
		test()
	default:
		panic("Illegal running mode used")
	}
}

func test()  {

}