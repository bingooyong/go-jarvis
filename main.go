package main

import (
	"flag"
	"github.com/lvyong1985/go-jarvis/config"
	"github.com/lvyong1985/go-jarvis/routers"
	"strconv"
	"fmt"
	"github.com/lvyong1985/go-jarvis/g"
	"os"
)

var c = flag.String("c", "./etc/go-jarvis-conf.yaml", "config file path")
var version = flag.Bool("v", false, "show version")

func main() {
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	cfg := config.Instance()
	cfg.Load(*c)

	router := routers.Router()
	router.Run(":" + strconv.Itoa(cfg.Server.Port))
}
