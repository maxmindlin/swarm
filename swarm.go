package main

import (
	"fmt"
	"os"

	"github.com/maxmindlin/swarm/workers"
)

func main() {
	start := os.Args[1]
	stories := workers.Crawl(start)
	fmt.Println(stories)
}
