package metrics

import (
	"fmt"
	"net/http"
	"strings"
)

// PrometheusHandler формирует HTTP-ответ с метриками в формате Prometheus.
type PrometheusHandler struct {
	collector *Collector
}

// NewPrometheusHandler - создает handler, который отдает runtime-метрики
func NewPrometheusHandler(collector *Collector) http.Handler {
	return &PrometheusHandler{
		collector: collector,
	}
}

// ServeHTTP - обрабатывает запрос к metrics
func (h *PrometheusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stats := h.collector.Collect()

	var sb strings.Builder

	writeGauge(&sb, "go_memory_alloc_bytes", "Currently allocated memory in bytes", float64(stats.AllocBytes))
	writeCounter(&sb, "go_memory_total_alloc_bytes", "Total allocated memory in bytes since start", float64(stats.TotalAllocBytes))
	writeGauge(&sb, "go_memory_sys_bytes", "Total memory obtained from the OS", float64(stats.SysBytes))

	writeCounter(&sb, "go_gc_mallocs_total", "Total number of mallocs", float64(stats.MallocsTotal))
	writeCounter(&sb, "go_gc_frees_total", "Total number of frees", float64(stats.FreesTotal))
	writeCounter(&sb, "go_gc_lookups_total", "Total number of pointer lookups", float64(stats.LookupsTotal))

	writeGauge(&sb, "go_memory_heap_alloc_bytes", "Heap allocated bytes", float64(stats.HeapAllocBytes))
	writeGauge(&sb, "go_memory_heap_sys_bytes", "Heap system bytes", float64(stats.HeapSysBytes))
	writeGauge(&sb, "go_memory_heap_idle_bytes", "Heap idle bytes", float64(stats.HeapIdleBytes))
	writeGauge(&sb, "go_memory_heap_inuse_bytes", "Heap in-use bytes", float64(stats.HeapInuseBytes))
	writeGauge(&sb, "go_memory_heap_released_bytes", "Heap released bytes", float64(stats.HeapReleasedBytes))
	writeGauge(&sb, "go_memory_heap_objects", "Number of heap objects", float64(stats.HeapObjects))

	writeGauge(&sb, "go_memory_stack_inuse_bytes", "Stack in-use bytes", float64(stats.StackInuseBytes))
	writeGauge(&sb, "go_memory_stack_sys_bytes", "Stack system bytes", float64(stats.StackSysBytes))

	writeGauge(&sb, "go_memory_mspan_inuse_bytes", "MSpan in-use bytes", float64(stats.MSpanInuseBytes))
	writeGauge(&sb, "go_memory_mspan_sys_bytes", "MSpan system bytes", float64(stats.MSpanSysBytes))
	writeGauge(&sb, "go_memory_mcache_inuse_bytes", "MCache in-use bytes", float64(stats.MCacheInuseBytes))
	writeGauge(&sb, "go_memory_mcache_sys_bytes", "MCache system bytes", float64(stats.MCacheSysBytes))

	writeGauge(&sb, "go_memory_buck_hash_sys_bytes", "Profiling bucket hash table bytes", float64(stats.BuckHashSysBytes))
	writeGauge(&sb, "go_memory_gc_sys_bytes", "GC metadata system bytes", float64(stats.GCSysBytes))
	writeGauge(&sb, "go_memory_other_sys_bytes", "Other system bytes", float64(stats.OtherSysBytes))

	writeGauge(&sb, "go_gc_next_gc_bytes", "Target heap size for the next GC cycle", float64(stats.NextGCBytes))
	writeCounter(&sb, "go_gc_cycles_total", "Total number of completed GC cycles", float64(stats.NumGC))
	writeCounter(&sb, "go_gc_forced_cycles_total", "Total number of forced GC cycles", float64(stats.NumForcedGC))
	writeCounter(&sb, "go_gc_pause_total_ns", "Total GC pause time in nanoseconds", float64(stats.PauseTotalNs))
	writeGauge(&sb, "go_gc_last_time_unix", "Unix timestamp of the last GC", float64(stats.LastGCUnix))
	writeGauge(&sb, "go_gc_cpu_fraction", "Fraction of CPU time used by GC", stats.GCCPUFraction)
	writeGauge(&sb, "go_gc_percent", "Current GCPercent value", float64(stats.GCPercent))

	writeGauge(&sb, "go_runtime_goroutines", "Current number of goroutines", float64(stats.Goroutines))
	writeCounter(&sb, "go_runtime_cgo_calls_total", "Total number of cgo calls", float64(stats.CgoCalls))
	writeGauge(&sb, "go_runtime_exporter_timestamp_unix", "Unix timestamp when metrics were collected", float64(stats.TimestampUnix))

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(sb.String()))
}

// writeGauge - записывает метрику типа gauge
func writeGauge(sb *strings.Builder, name, help string, value float64) {
	fmt.Fprintf(sb, "# HELP %s %s\n", name, help)
	fmt.Fprintf(sb, "# TYPE %s gauge\n", name)
	fmt.Fprintf(sb, "%s %v\n", name, value)
}

// writeCounter - записывает метрику типа counter
func writeCounter(sb *strings.Builder, name, help string, value float64) {
	fmt.Fprintf(sb, "# HELP %s %s\n", name, help)
	fmt.Fprintf(sb, "# TYPE %s counter\n", name)
	fmt.Fprintf(sb, "%s %v\n", name, value)
}
