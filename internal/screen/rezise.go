package screen

import (
	"time"

	"github.com/EnotInc/Bard/config"
	"golang.org/x/term"
)

func TermSizeMonitor() {
	go captureResize()
	go listenResize()
}

func captureResize() {
	cfg := config.GetConfig()
	duration := time.Duration(cfg.ResizeTime)
	ticker := time.NewTicker(duration * time.Microsecond)
	defer ticker.Stop()

	var last_w, last_h = global.w, global.h

	for range ticker.C {
		w, h, err := term.GetSize(global.fdOut)
		if err != nil {
			continue
		}

		if last_w != w || last_h != h {
			last_w = w
			last_h = h

			global.w = w
			global.h = h
			global.redraw <- true
		}
	}
}

func listenResize() {
	var last_w, last_h = global.w, global.h
	for {
		changed := <-global.redraw
		if changed {
			ofset := len(global.tiles)

			diff_w := (global.w - last_w) / ofset
			diff_h := global.h - last_h

			last_w = global.w
			last_h = global.h

			for _, t := range global.tiles {
				t.w += diff_w
				t.h += diff_h
				t.hash = make(map[int]uint32)
				t.object.Resize(t.w, t.w)
			}

			DrawAll()
		}
	}
}
