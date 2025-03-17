package cgomemprof_test

import (
	"net/http"
	_ "net/http/pprof"
	"testing"

	_ "github.com/milvus-io/chyezh/cgomemprof"
)

func TestHTTP(t *testing.T) {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
