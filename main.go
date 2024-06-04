package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"math/rand"
	"time"

	"github.com/ankitpaudel20/algoviz_go/sorting"
	"github.com/ebitenui/ebitenui/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	LightBlue     color.Color = color.RGBA{204, 245, 255, 255}
	Pink          color.Color = color.RGBA{255, 51, 85, 255}
	LightPurple   color.Color = color.RGBA{191, 162, 208, 255}
	DarkTeal      color.Color = color.RGBA{5, 83, 107, 255}
	old_time      int64       = time.Now().UnixMicro()
	frames        int64       = 0
	ticks         int64       = 0
	fps           float32     = 0
	tps           float32     = 0
	oldWindowSize image.Point = image.Pt(1000, 800)
)

// ///////////////////////////////////////////////
type CanvasState int

const (
	CanvasSortState CanvasState = iota
	CanvasGridState
)

// ///////////////////////////////////////////////
const (
	leftPad  int = 15
	rightPad int = 15
	topPad   int = 50
	botPad   int = 15
)

// ///////////////////////////////////////////////
type Game struct {
	currentState  CanvasState
	sortStateInfo *sorting.SortCanvas
	// gridStateInfo GridCanvas
	bgColor color.Color
	ui      *GameUI
}

func (g *Game) PlayPause() {
	switch g.currentState {
	case CanvasSortState:
		g.sortStateInfo.PlayPause()
	case CanvasGridState:
		log.Print("grid state is not implemented yet")
	}
}

// ///////////////////////////////////////////////
func getRandomArray(length int) []int {
	default_nums := make([]int, length)
	for i := range length {
		default_nums[i] = i
	}

	rng := rand.New(rand.NewSource(time.Now().Unix()))
	rng.Shuffle(len(default_nums), func(i, j int) {
		default_nums[i], default_nums[j] = default_nums[j], default_nums[i]
	})
	return default_nums
}

func get_default_starting_game_config() *Game {
	length := 100
	default_nums := getRandomArray(length)
	// default_nums := []int{8, 0, 4, 9, 2, 7, 5, 3, 6, 1}

	// default_nums := make([]int, length)
	// for i := range length {
	// 	default_nums[i] = i
	// }

	// log.Printf("%v", default_nums)
	canvasSize := image.Rect(leftPad, topPad, oldWindowSize.X-rightPad, oldWindowSize.Y-botPad)
	initial_sort_config := sorting.InitSortConfig{}
	initial_sort_config.Nums = default_nums
	initial_sort_config.Pos = canvasSize
	initial_sort_config.DefaultColor = LightPurple
	initial_sort_config.DefaultSort = sorting.SortTypeQuicksort
	// initial_sort_config.Nums = default_nums
	// initial_sort_config.Nums = default_nums
	// intial_sort_config.Nums = default_nums
	// sorting.InitSortConfig
	g := Game{
		currentState:  CanvasSortState,
		sortStateInfo: sorting.NewSortCanvas(&initial_sort_config),
		bgColor:       DarkTeal,
		ui:            getDefaultUI(),
	}
	return &g
}

func main() {
	ebiten.SetWindowSize(oldWindowSize.X, oldWindowSize.Y)
	ebiten.SetWindowTitle("Displaying image with a triangle overlay")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetVsyncEnabled(false)

	g := get_default_starting_game_config()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// ///////////////////////////////////////////////
func (g *Game) Update() error {
	ticks++
	input.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.PlayPause()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		switch g.currentState {
		case CanvasGridState:
			log.Print("Grid not yet implemented")
		case CanvasSortState:
			g.sortStateInfo.Start()
		}
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		switch g.currentState {
		case CanvasGridState:
			log.Print("Grid not yet implemented")
		case CanvasSortState:
			g.sortStateInfo.Nums = getRandomArray(len(g.sortStateInfo.Nums))
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.bgColor)
	g.ui.uiElements.Draw(screen)

	switch g.currentState {
	case CanvasGridState:
		log.Print("grid canvas not implemented yet")
	case CanvasSortState:
		g.sortStateInfo.Draw(screen)
	default:
		log.Fatal("this should not happen")
	}

	if time.Now().UnixMicro()-old_time > 1000*1000*0.2 {
		fps = 1000000 * float32(frames) / float32(time.Now().UnixMicro()-old_time)
		tps = 1000000 * float32(ticks) / float32(time.Now().UnixMicro()-old_time)
		old_time = time.Now().UnixMicro()
		frames = 0
		ticks = 0
	}
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("fps: %2.1f \t tps: %2.1f", fps, tps))
	frames++

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if oldWindowSize.Y != outsideHeight || oldWindowSize.X != outsideWidth {
		oldWindowSize.X, oldWindowSize.Y = outsideWidth, outsideHeight
		canvasSize := image.Rect(leftPad, topPad, outsideWidth-rightPad, outsideHeight-botPad)
		g.sortStateInfo.Position = canvasSize
		g.sortStateInfo.UpdateSizes()
	}
	return outsideWidth, outsideHeight
}
