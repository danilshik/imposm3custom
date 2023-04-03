package stats

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/danilshik/imposm3custom/log"
)

func StartHTTPPProf(bind string) {
	go func() {
		log.Println(http.ListenAndServe(bind, nil))
	}()
}
