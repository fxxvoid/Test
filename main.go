package main

import (
	"asterism/caches"
	"asterism/servers"
	"flag"
	"log"
)

func main() {
	address := flag.String("address", ":5837", "The address used to listen, such as 127.0.0.1:5837")
	options := caches.DefaultOptions()
	flag.Int64Var(&options.MaxEntrySize, "maxEntrySize", options.MaxEntrySize, "The max memory size that entries can use. The unit is GB.")
	flag.IntVar(&options.MaxGcCount, "maxGcCount", options.MaxGcCount, "The max count of entries that gc will clean.")
	flag.Int64Var(&options.GcDuration, "gcDuration", options.GcDuration, "The duration between gc tasks, The unit is Minute.")
	flag.StringVar(&options.DumpFile, "dumpFile", options.DumpFile, "The file is used to dumpfile.")
	flag.Int64Var(&options.DumpDuration, "dumpDuration", options.DumpDuration, "The gap between two dump, The unit is Minute.")
	flag.Parse()

	cache := caches.NewCacheWith(options)
	cache.AutoGc()
	cache.AutoDump()
	log.Printf("Kafo is running on %s.", *address)
	err := servers.NewHTTPServer(cache).Run(*address)

	if err != nil {
		panic(err)
	}
}
