package cgomemprof

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
)

const URLPathPrefix = "/debug/jemallocprof/pprof/"

func init() {
	http.HandleFunc(URLPathPrefix+"heap", Heap)
	http.HandleFunc(URLPathPrefix+"symbol", Symbol)
	http.HandleFunc(URLPathPrefix+"cmdline", pprof.Cmdline)
	http.HandleFunc(URLPathPrefix+"active", Active)
}

func Symbol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// We have to read the whole POST body before
	// writing any output. Buffer the output here.
	var buf bytes.Buffer

	if r.Method == "POST" {
		b := bufio.NewReader(r.Body)
		for {
			word, err := b.ReadSlice('+')
			if err == nil {
				word = word[0 : len(word)-1] // trim +
			}
			pc, _ := strconv.ParseUint(string(word), 0, 64)
			if pc != 0 {
				symbol := GetSymbol(pc)
				fmt.Fprintf(&buf, "%s\n", symbol)
			}

			// Wait until here to check for err; the last
			// symbol will have an err because it doesn't end in +.
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(&buf, "reading request: %v\n", err)
				}
				break
			}
		}
	} else {
		// We don't know how many symbols we have, but we
		// do have symbol information. Pprof only cares whether
		// this number is 0 (no symbols available) or > 0.
		fmt.Fprintf(&buf, "num_symbols: 1\n")
	}

	w.Write(buf.Bytes())
}

func Heap(w http.ResponseWriter, r *http.Request) {
	tmpFile, err := os.CreateTemp("", "memprofile-*.dump")
	if err != nil {
		http.Error(w, "could not create temp file to dump", http.StatusInternalServerError)
		return
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	if err := DumpMemoryProfileIntoFile(tmpFile.Name()); err != nil {
		http.Error(w, "could not dump memory profile", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, tmpFile.Name())
}

func Active(w http.ResponseWriter, r *http.Request) {
	enable := r.URL.Query().Get("enable")
	enableNum, err := strconv.ParseInt(enable, 10, 64)
	if err != nil {
		http.Error(w, "invalid enable value", http.StatusBadRequest)
		return
	}
	if enableNum != 0 {
		EnableMemoryProfiling()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("jemalloc memprof enabled"))
	} else {
		DisableMemoryProfiling()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("jemalloc memprof disabled"))
	}
}
