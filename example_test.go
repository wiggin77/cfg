package cfg_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/wiggin77/cfg"
)

func sampleMap() map[string]string {
	return map[string]string{
		"fps":        "30",
		"retryDelay": "1 minute",
		"logRotate":  "1 day",
		"ratio":      "1.85",
	}
}

type listener struct {
	// empty
}

func (l *listener) ConfigChanged(cfg *cfg.Config, src cfg.SourceMonitored) {
	fmt.Println("Config changed!")
}

func Example() {
	// create a Config instance
	cfg := &cfg.Config{}
	// shutdown will stop monitoring the sources for changes
	defer cfg.Shutdown()

	// for this sample use a source backed by a simple map
	m := sampleMap()
	src := cfg.NewSrcMapFromMap(m)

	// add the source to the end of the searched sources
	cfg.AppendSource(src)

	// add a source to the beginning of the searched sources,
	// providing defaults for missing properties.
	cfg.PrependSource(cfg.NewSrcMapFromMap(map[string]string{"maxRetries": "10"}))

	// listen for changes (why not use a func type here intead of interface? Because we
	// need to be able to remove listeners and cannot do that with funcs).
	cfg.AddChangedListener(&listener{})

	// change a property every 1 seconds for 5 seconds.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	done := time.After(5 * time.Second)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		case <-ticker.C:
			m["fps"] = strconv.Itoa(rnd.Intn(30))
		case <-done:
			return
		}
	}
}
