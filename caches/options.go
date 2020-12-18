package caches

type Options struct {
	MaxEntrySize int64
	MaxGcCount int
	GcDuration int64
	DumpFile string
	DumpDuration int64
}

func DefaultOptions() Options {
	return Options{
		MaxEntrySize: int64(4),
		MaxGcCount: 1000,
		GcDuration: 60,
		DumpFile: "cache-server.dump",
		DumpDuration: 30,
	}
}
