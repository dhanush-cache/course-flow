package adapters

import (
	"fmt"

	"github.com/vbauerster/mpb/v8/decor"
)

func Line(fn decor.DecorFunc, wcc ...decor.WC) decor.Decorator {
	return line{initWC(wcc...), fn}
}

type line struct {
	decor.WC
	fn decor.DecorFunc
}

func (d line) Decor(s decor.Statistics) (string, int) {
	result := d.fn(s)
	format, _ := d.Format(result)
	return fmt.Sprintf("%s\n", format), 0
}

func initWC(wcc ...decor.WC) decor.WC {
	var wc decor.WC
	for _, nwc := range wcc {
		wc = nwc
	}
	return wc.Init()
}
