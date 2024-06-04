package main

import (
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

// "golang.org/x/image/font"

type uiResources struct {
	background *image.NineSlice
	// fonts *font

	// separatorColor color.Color

	// text        *textResources
	// button      *buttonResources
	// label       *labelResources
	// checkbox    *checkboxResources
	// comboButton *comboButtonResources
	// list        *listResources
	// slider      *sliderResources
	// progressBar *progressBarResources
	// panel       *panelResources
	// tabBook     *tabBookResources
	// header      *headerResources
	// textInput   *textInputResources
	// textArea    *textAreaResources
	// toolTip     *toolTipResources
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.NRGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}

type GameUI struct {
	uiElements  ebitenui.UI
	uiResources uiResources
}

const (
	backgroundColor = "131a22"
)

func initResources() *uiResources {
	return &uiResources{
		background: image.NewNineSliceColor(hexToColor(backgroundColor)),
	}
}

func getDefaultUI() *GameUI {
	//This creates the root container for this UI.
	res := initResources()

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			// It is using a GridLayout with a single column
			widget.GridLayoutOpts.Columns(1),
			// It uses the Stretch parameter to define how the rows will be layed out.
			// - a fixed sized header
			// - a content row that stretches to fill all remaining space
			// - a fixed sized footer
			// widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			// Padding defines how much space to put around the outside of the grid.
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			// Spacing defines how much space to put between each column and row
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(res.background),
	)
	ui := &ebitenui.UI{
		Container: rootContainer,
	}

	gUI := GameUI{uiElements: *ui, uiResources: *res}
	return &gUI
}
