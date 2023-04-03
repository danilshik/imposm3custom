package writer

import (
	"runtime"
	"sync"

	"github.com/danilshik/imposm3custom/cache"
	"github.com/danilshik/imposm3custom/database"
	"github.com/danilshik/imposm3custom/expire"
	"github.com/danilshik/imposm3custom/geom/limit"
	"github.com/danilshik/imposm3custom/proj"
	"github.com/danilshik/imposm3custom/stats"
	osm "github.com/omniscale/go-osm"
)

type ErrorLevel interface {
	Level() int
}

type looper interface {
	loop()
}

type OsmElemWriter struct {
	osmCache   *cache.OSMCache
	diffCache  *cache.DiffCache
	progress   *stats.Statistics
	inserter   database.Inserter
	wg         *sync.WaitGroup
	limiter    *limit.Limiter
	writer     looper
	srid       int
	expireor   expire.Expireor
	concurrent bool
}

func (writer *OsmElemWriter) SetLimiter(limiter *limit.Limiter) {
	writer.limiter = limiter
}

func (writer *OsmElemWriter) EnableConcurrent() {
	writer.concurrent = true
}

func (writer *OsmElemWriter) Start() {
	concurrency := 1
	if writer.concurrent {
		concurrency = runtime.NumCPU()
	}
	for i := 0; i < concurrency; i++ {
		writer.wg.Add(1)
		go writer.writer.loop()
	}
}

func (writer *OsmElemWriter) SetExpireor(exp expire.Expireor) {
	writer.expireor = exp
}

func (writer *OsmElemWriter) Wait() {
	writer.wg.Wait()
}

func (writer *OsmElemWriter) NodesToSrid(nodes []osm.Node) {
	if writer.srid == 4326 {
		return
	}
	if writer.srid != 3857 {
		panic("invalid srid. only 4326 and 3857 are supported")
	}

	for i, nd := range nodes {
		nodes[i].Long, nodes[i].Lat = proj.WgsToMerc(nd.Long, nd.Lat)
	}
}

func (writer *OsmElemWriter) NodeToSrid(node *osm.Node) {
	if writer.srid == 4326 {
		return
	}
	if writer.srid != 3857 {
		panic("invalid srid. only 4326 and 3857 are supported")
	}
	node.Long, node.Lat = proj.WgsToMerc(node.Long, node.Lat)
}
