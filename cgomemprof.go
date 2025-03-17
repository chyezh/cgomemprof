package cgomemprof

// #cgo LDFLAGS: -L/home/chyezh/repository/jemalloc/jemalloc/lib/ -ljemalloc -lbacktrace
// #include "cgomemprof.h"
// #include "stdlib.h"
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

var mu sync.RWMutex

// EnableMemoryProfiling enables memory profiling.
func EnableMemoryProfiling() error {
	mu.Lock()
	defer mu.Unlock()

	if code := C.EnableMemoryProfiling(); code != 0 {
		return fmt.Errorf("EnableMemoryProfiling failed with code %d", code)
	}
	return nil
}

func DisableMemoryProfiling() error {
	mu.Lock()
	defer mu.Unlock()

	if code := C.DisableMemoryProfiling(); code != 0 {
		return fmt.Errorf("DisableMemoryProfiling failed with code %d", code)
	}
	return nil
}

func DumpMemoryProfileIntoFile(filename string) error {
	mu.RLock()
	defer mu.RUnlock()

	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	if code := C.DumpMemoryProfileIntoFile(cfilename); code != 0 {
		return fmt.Errorf("DumpMemoryProfileIntoFile failed with code %d", code)
	}
	return nil
}

func GetSymbol(addr uint64) string {
	result := C.GetSymbol(C.uintptr_t(addr))
	r := C.GoString(result)
	C.free(unsafe.Pointer(result))
	return r
}
