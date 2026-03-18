package debug

import (
	"fmt"
	"net/http"
	"strconv"
)

// sink - удерживает выделенную память
var sink [][]byte

// Register - регистрирует debug endpoint
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/debug/alloc", allocHandler)
}

// allocHandler - создает указанное количество мегабайт в памяти
func allocHandler(w http.ResponseWriter, r *http.Request) {
	mbStr := r.URL.Query().Get("mb")
	if mbStr == "" {
		mbStr = "10"
	}

	mb, err := strconv.Atoi(mbStr)
	if err != nil {
		http.Error(w, "invalid mb param", http.StatusBadRequest)
		return
	}

	for i := 0; i < mb; i++ {
		b := make([]byte, 1024*1024) // 1MB
		sink = append(sink, b)
	}

	if _, err := fmt.Fprintf(w, "allocated %d MB\n", mb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
