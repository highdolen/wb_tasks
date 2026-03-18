package meminfo

//Stats - хранит snapshot runtime-метрик памяти.
type Stats struct {
	AllocBytes      uint64
	TotalAllocBytes uint64
	SysBytes        uint64
	LookupsTotal    uint64
	MallocsTotal    uint64
	FreesTotal      uint64

	HeapAllocBytes    uint64
	HeapSysBytes      uint64
	HeapIdleBytes     uint64
	HeapInuseBytes    uint64
	HeapReleasedBytes uint64
	HeapObjects       uint64

	StackInuseBytes uint64
	StackSysBytes   uint64

	MSpanInuseBytes  uint64
	MSpanSysBytes    uint64
	MCacheInuseBytes uint64
	MCacheSysBytes   uint64

	BuckHashSysBytes uint64
	GCSysBytes       uint64
	OtherSysBytes    uint64

	NextGCBytes   uint64
	LastGCUnix    int64
	PauseTotalNs  uint64
	NumGC         uint32
	NumForcedGC   uint32
	GCCPUFraction float64
	GCPercent     int

	Goroutines    int
	CgoCalls      int64
	TimestampUnix int64
}
