package cui

import (
	"context"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

type mode int

const (
	CMDMode mode = iota
	CUIMmode
)

var (
	AppMode = CMDMode
)

const(
	logView = "log"
	logChannelSize = 20
	logX = 1
	logY = 0.4

	orderView = "order"
	orderX = 0.3
	orderY = 1 - logY

	cookView = "cook"
	cookX = 1 - orderX
	cookY = 1 - logY
)

var(
	LogListener = make(chan string, logChannelSize)
)

func buildLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(cookView, 0, 0, int(cookX*float32(maxX))-1, int(cookY*float32(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Cook status"
	}
	if v, err := g.SetView(orderView, int(cookX*float32(maxX)), 0, maxX-1, int(orderY*float32(maxY))-1);  err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Order Queue"
	}
	if v, err := g.SetView(logView, 0, int(cookY*float32(maxY)), maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Logs"
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func CUIInit(ctx context.Context, cancel context.CancelFunc) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {

	}
	defer g.Close()

	g.SetManagerFunc(buildLayout)

	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	go listen(g, ctx)

	AppMode = CUIMmode
	if err = g.MainLoop(); err != nil {
		if err == gocui.ErrQuit {
			cancel()
			AppMode = CMDMode
		}
		log.Panicln(err)
	}
}

func listen(g *gocui.Gui, ctx context.Context) {
	for{
		select {
		case <-ctx.Done():
			return
		case data := <-LogListener:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View(logView)
				if err != nil {
					return err
				}
				fmt.Fprintln(v, data)
				return nil
			})
		}
	}
}