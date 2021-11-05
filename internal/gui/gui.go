package gui

import (
	"context"
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
	"github.com/foxfurry/go_kitchen/internal/domain/repository"
	"github.com/foxfurry/go_kitchen/internal/service/cook"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
)

type ICUI interface {
	Create() bool
	Start(ctx context.Context, cancel context.CancelFunc)
	AddLogData(data string)
	AddCookData(cookNum int, cookData []cook.CookSlot)
}

type CookCUI struct {
	progress []*widgets.Gauge
	header *widgets.Paragraph
	items *widgets.List
}

type kitchenCUI struct {
	logs *widgets.List
	header *widgets.Paragraph
	tabs *widgets.TabPane
	orders *widgets.List
	cooks []*CookCUI
}

type mode int

const (
	CMDMode mode = iota
	CUIMode
)

var (
	AppMode = CMDMode
	lastGaugeY int
)

const (
	headerYStart = 0
	headerYEnd   = headerYStart + 1

	tabYStart = headerYEnd + 1
	tabYEnd   = tabYStart + 1

	cookHeaderYStart = tabYEnd + 1
	cookHeaderYEnd = cookHeaderYStart + 1

	logView   = "Current Logs"
	logYStart = tabYEnd + 3
	logX = 0.7

	gaugeYStart = logYStart
	gaugeYSize = 4
	gaugeX = logX

	orderYStart = logYStart
	orderXStart = logX + 0.01

	itemsYStart = logX
	itemsXStart = 0
	itemsX = gaugeX
)

func NewKitchenCUI() ICUI {
	return &kitchenCUI{
		logs: nil,
	}
}

func (k *kitchenCUI) Create() bool{
	if !viper.GetBool("enable_cui") {
		return false
	}

	cooks := repository.GetCooks()

	if err := ui.Init(); err != nil {
		return false
	}

	k.logs = buildLogWidget()
	k.header = buildHeaderWidget()
	k.tabs = buildTabPane(len(cooks))
	k.orders = buildOrderWidget()

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

func (k *kitchenCUI) AddLogData(data string) {
	k.logs.Rows = append(k.logs.Rows, data)
	k.logs.ScrollBottom()
	k.render()
}

func (k *kitchenCUI) AddCookData(cookNum int, cookData []cook.CookSlot){
	k.updateCook(cookNum, cookData)
}

func (k *kitchenCUI) updateCook(cookNum int, cookData []cook.CookSlot){
	for idx, _ := range k.cooks[cookNum].progress {
		k.cooks[cookNum].progress[idx].Percent = cookData[idx].Progress	// Update all gauges
	}
}

func (k *kitchenCUI) render(){
	ui.Clear()
	ui.Render(k.header, k.tabs)
	ui.Render(k.orders)
	switch k.tabs.ActiveTabIndex {
	case 0:
		ui.Render(k.logs)
	default:
		k.cooks[k.tabs.ActiveTabIndex - 1].Render()
	}
}

func (c *CookCUI) Render(){
	ui.Render(c.header)
	ui.Render(c.items)
	for idx, _ := range c.progress{
		ui.Render(c.progress[idx])
	}
}

func buildLogWidget() *widgets.List {
	maxX, maxY := ui.TerminalDimensions()
	view := widgets.NewList()
	view.Title = logView
	view.WrapText = true

	view.SetRect(0, logYStart, int(float64(maxX)*logX), maxY)
	return view
}

func buildOrderWidget() *widgets.List {
	maxX, maxY := ui.TerminalDimensions()
	view := widgets.NewList()
	view.Title = "Current Orders"
	view.WrapText = false

	view.SetRect(int(float64(maxX)*orderXStart), orderYStart, maxX, maxY)
	return view
}

func buildHeaderWidget() *widgets.Paragraph {
	maxX, _ := ui.TerminalDimensions()

	header := widgets.NewParagraph()
	header.Text = "Press q (C^c) to quit. Press h or l to switch tabs"
	header.Border = false

	header.SetRect(0, headerYStart, maxX, headerYEnd)

	return header
}

func buildTabPane(cooksNum int) *widgets.TabPane {
	maxX, _ := ui.TerminalDimensions()

	tab := widgets.NewTabPane()
	tab.Border=true

	tabsNames := []string{"logs"}

	for idx := 0; idx < cooksNum; idx++ {
		tabsNames = append(tabsNames, "cook #" + strconv.Itoa(idx + 1))
	}

	tab.TabNames = tabsNames
	tab.SetRect(0, tabYStart, maxX, tabYEnd)

	return tab
}

func buildCookCUI(cook entity.Cook) *CookCUI {
	return &CookCUI{
		progress: buildBarsGauge(cook.Proficiency),
		header:   buildBarHeader(cook.Name, cook.CatchPhrase),
		items:    buildBarItems(cook.Proficiency),
	}
}

func buildBarHeader(cookName, cookCatch string) *widgets.Paragraph {
	maxX, _ := ui.TerminalDimensions()

	header := widgets.NewParagraph()
	header.Text = cookName + " | " + cookCatch
	header.Border = false

	header.SetRect(0, cookHeaderYStart, maxX, cookHeaderYEnd)

	return header
}

func buildBarsGauge(proficiency int) []*widgets.Gauge {
	maxX, _ := ui.TerminalDimensions()

	var result []*widgets.Gauge

	for idx := 0; idx < proficiency; idx++{
		tmp := widgets.NewGauge()
		tmp.Title = "Item #" + strconv.Itoa(idx)
		tmp.Percent = rand.Int()%100

		tmp.BarColor = ui.ColorRed
		tmp.LabelStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
		tmpStart := gaugeYStart + gaugeYSize*idx
		tmp.SetRect(0, tmpStart, int(float64(maxX)*gaugeX), tmpStart+gaugeYSize)

		lastGaugeY = tmpStart+ gaugeYSize
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

	items.SetRect(int(float32(maxX)*itemsXStart), lastGaugeY, int(float64(maxX)*itemsX), maxY)

	return items
}


