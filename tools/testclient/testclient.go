package main

import (
	"context"
	"time"
	"flag"
	"fmt"
	"strconv"
	
	"github.com/wangthomas/gobloomfield/client"
)

var (
	create   bool
	drop     bool
	add      bool
	has      bool
	hostname string
)

func init() {
	flag.BoolVar(&create, "create", false, "create a filter <filter name>")
	flag.BoolVar(&drop, "drop", false, "drop a filter <filter name>")
	flag.BoolVar(&add, "add", false, "add a key to a filter <filter name> <key ...>")
	flag.BoolVar(&has, "has", false, "check a key in a filter <filter name> <key ...>")
	flag.StringVar(&hostname, "hostname", "localhost:8679", "Specifc hostname:port")
}


func main() {

	flag.Parse()

	client, err := client.NewBloomClient(hostname, time.Millisecond*5000)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	switch {
		case create:
			filter := flag.Arg(0)
			err = client.Create(ctx, filter)

		case drop:
			filter := flag.Arg(0)
			err = client.Drop(ctx, filter)

		case add:
			filter := flag.Arg(0)

			hash1, _ := strconv.ParseUint(flag.Args()[1], 10, 64)
			hash2, _ := strconv.ParseUint(flag.Args()[2], 10, 64)

			var hasKey bool
			hasKey, err = client.Add(ctx, filter, []uint64{hash1, hash2})
			fmt.Println(hasKey)

		case has:
			filter := flag.Arg(0)
			hash1, _ := strconv.ParseUint(flag.Args()[1], 10, 64)
			hash2, _ := strconv.ParseUint(flag.Args()[2], 10, 64)

			var hasKey bool
			hasKey, err = client.Has(ctx, filter, []uint64{hash1, hash2})
			fmt.Println(hasKey)
	}
	if err != nil {
		fmt.Println(err)
	}
}