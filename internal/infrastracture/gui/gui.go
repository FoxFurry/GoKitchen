package gui

import (
	"context"
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
	"github.com/foxfurry/go_kitchen/internal/domain/repository"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"math/rand"
	"strconv"
)

type ICUI interface {
	Create() bool
	Start(ctx context.Context, cancel context.CancelFunc)
	AddLog(data string)
}

type CookCUI struct {
	progress []*widgets.Gauge
	header *widgets.Paragraph
	catchPhrase *widgets.Paragraph
	items *widgets.List
}

type kitchenCUI struct {
	logs *widgets.List
	header *widgets.Paragraph
	tabs *widgets.TabPane
	cooks []*CookCUI
}

type mode int

const (
	CMDMode mode = iota
	CUIMode
)

var (
	AppMode = CMDMode
)

const (
	headerYStart = 0
	headerYEnd   = headerYStart + 1

	tabYStart = headerYEnd + 1
	tabYEnd   = tabYStart + 1

	cookHeaderYStart = tabYEnd + 1
	cookHeaderYEnd = cookHeaderYStart + 1

	cookCatchPhraseYStart = cookHeaderYEnd + 1
	cookCatchPhraseYEnd = cookCatchPhraseYStart + 1

	logView   = "log"
	logYStart = tabYEnd + 2

	gaugeYStart = logYStart
	gaugeYSize = 4
	gaugeYGap = 1
	gaugeX = 0.7

	itemsYStart = logYStart
	itemsXStart = gaugeX + 0.01
)

func NewKitchenCUI() ICUI {
	return &kitchenCUI{
		logs: nil,
	}
}

func (k *kitchenCUI) Create() bool{
	cooks := repository.GetCooks()

	if err := ui.Init(); err != nil {
		return false
	}

	k.logs = buildLogWidget()
	k.header = buildHeaderWidget()
	k.tabs = buildTabPane(len(cooks))

	for _, val := range cooks {
		k.cooks = append(k.cooks, buildCookCUI(val))
	}

	return true
}

func (k *kitchenCUI) Start(ctx context.Context, cancel context.CancelFunc) {
	defer ui.Close()
	k.render()

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				AppMode = CMDMode
				cancel()
				return
			case "h":
				k.tabs.FocusLeft()
				k.render()
			case "l":
				k.tabs.FocusRight()
				k.render()
			}
		case <-ctx.Done():
			return
		}
	}
}

func (k *kitchenCUI) AddLog(data string) {
	k.logs.Rows = append(k.logs.Rows, data)
	k.logs.ScrollBottom()
	k.render()
}

func (k *kitchenCUI) render(){
	ui.Clear()
	ui.Render(k.header, k.tabs)

	switch k.tabs.ActiveTabIndex {
	case 0:
		ui.Render(k.logs)
	default:
		k.cooks[k.tabs.ActiveTabIndex - 1].Render()
	}
}

func (c *CookCUI) Render(){
	ui.Render(c.header)
	ui.Render(c.catchPhrase)
	ui.Render(c.items)
	for idx, _ := range c.progress{
		ui.Render(c.progress[idx])
	}
}

func buildLogWidget() *widgets.List {
	maxX, maxY := ui.TerminalDimensions()
	view := widgets.NewList()
	view.Title = logView
	view.WrapText = false

	view.SetRect(0, logYStart, maxX-1, maxY-1)
	return view
}

func buildHeaderWidget() *widgets.Paragraph {
	maxX, _ := ui.TerminalDimensions()

	header := widgets.NewParagraph()
	header.Text = "Press q (C^c) to quit. Press h or l to switch tabs"
	header.Border = false

	header.SetRect(0, headerYStart, maxX-1, headerYEnd)

	return header
}

func buildTabPane(cooksNum int) *widgets.TabPane {
	maxX, _ := ui.TerminalDimensions()

	tab := widgets.NewTabPane()
	tab.Border=true

	tabsNames := []string{"Logs"}

	for idx := 0; idx < cooksNum; idx++ {
		tabsNames = append(tabsNames, "cook #" + strconv.Itoa(idx + 1))
	}

	tab.TabNames = tabsNames
	tab.SetRect(0, tabYStart, maxX, tabYEnd)

	return tab
}

func buildCookCUI(cook entity.Cook) *CookCUI {
	return &CookCUI{
		progress:    buildBarsGauge(cook.Proficiency),
		header:      buildBarHeader(cook.Name),
		catchPhrase: buildBarCatch(cook.CatchPhrase),
		items:       buildBarItems(cook.Proficiency),
	}
}

func buildBarHeader(cookName string) *widgets.Paragraph {
	maxX, _ := ui.TerminalDimensions()

	header := widgets.NewParagraph()
	header.Text = cookName
	header.Border = false

	header.SetRect(0, cookHeaderYStart, maxX-1, cookHeaderYEnd)

	return header
}

func buildBarCatch(catchPhrase string)  *widgets.Paragraph{
	maxX, _ := ui.TerminalDimensions()

	header := widgets.NewParagraph()
	header.Text = catchPhrase
	header.Border = false

	header.SetRect(0, cookCatchPhraseYStart, maxX-1, cookCatchPhraseYEnd)

	return header
}

func buildBarsGauge(proficiency int) []*widgets.Gauge {
	maxX, _ := ui.TerminalDimensions()

	var result []*widgets.Gauge

	for idx := 0; idx < proficiency; idx++{
		tmp := widgets.NewGauge()
		tmp.Title = "Item #" + strconv.Itoa(idx)
		tmp.Percent = rand.Int()%100

		tmp.BarColor = ui.ColorBlack
		tmp.LabelStyle = ui.NewStyle(ui.ColorWhite)
		tmpStart := gaugeYStart + gaugeYSize*idx + gaugeYGap
		tmp.SetRect(0, tmpStart, int(float64(maxX)*gaugeX), tmpStart+gaugeYSize)

		result = append(result, tmp)
	}

	return result
}

func buildBarItems(proficiency int) *widgets.List{
	maxX, maxY := ui.TerminalDimensions()

	items := widgets.NewList()
	items.Title = "Current items: "
	var itemsRows []string
	for idx := 0; idx < proficiency; idx++ {
		itemsRows = append(itemsRows, "#" + strconv.Itoa(idx) + ": n/a")
	}

	items.Rows = itemsRows

	items.SetRect(int(float32(maxX)*itemsXStart), itemsYStart, maxX-1, maxY-1)

	return items
}


