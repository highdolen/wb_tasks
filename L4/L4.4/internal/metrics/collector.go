package metrics

import (
	"runtime"
	"runtime/debug"
	"time"

	"gcMetrics/pkg/meminfo"
)

// Collector собирает runtime-метрики памяти и GC
type Collector struct{}

// NewCollector - конструктор Collector
func NewCollector() *Collector {
	return &Collector{}
}

// Collect - читает runtime.MemStats
func (c *Collector) Collect() meminfo.Stats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	gcPercent := debug.SetGCPercent(-1)
	debug.SetGCPercent(gcPercent)

	return meminfo.Stats{
		AllocBytes:        ms.Alloc,
		TotalAllocBytes:   ms.TotalAlloc,
		SysBytes:          ms.Sys,
		LookupsTotal:      ms.Lookups,
		MallocsTotal:      ms.Mallocs,
		FreesTotal:        ms.Frees,
		HeapAllocBytes:    ms.HeapAlloc,
		HeapSysBytes:      ms.HeapSys,
		HeapIdleBytes:     ms.HeapIdle,
		HeapInuseBytes:    ms.HeapInuse,
		HeapReleasedBytes: ms.HeapReleased,
		HeapObjects:       ms.HeapObjects,
		StackInuseBytes:   ms.StackInuse,
		StackSysBytes:     ms.StackSys,
		MSpanInuseBytes:   ms.MSpanInuse,
		MSpanSysBytes:     ms.MSpanSys,
		MCacheInuseBytes:  ms.MCacheInuse,
		MCacheSysBytes:    ms.MCacheSys,
		BuckHashSysBytes:  ms.BuckHashSys,
		GCSysBytes:        ms.GCSys,
		OtherSysBytes:     ms.OtherSys,
		NextGCBytes:       ms.NextGC,
		LastGCUnix:        nsToUnixSeconds(ms.LastGC),
		PauseTotalNs:      ms.PauseTotalNs,
		NumGC:             ms.NumGC,
		NumForcedGC:       ms.NumForcedGC,
		GCCPUFraction:     ms.GCCPUFraction,
		Goroutines:        runtime.NumGoroutine(),
		CgoCalls:          runtime.NumCgoCall(),
		GCPercent:         gcPercent,
		TimestampUnix:     time.Now().Unix(),
	}
}

// nsToUnixSeconds - переводит время из наносекунд в Unix
func nsToUnixSeconds(ns uint64) int64 {
	if ns == 0 {
		return 0
	}
	return int64(ns / uint64(time.Second))
}
