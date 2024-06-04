package sorting

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/ankitpaudel20/algoviz_go/pubsub"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/math/f32"
)

var (
	LightBlue   color.Color = color.RGBA{204, 245, 255, 255}
	Pink        color.Color = color.RGBA{255, 51, 85, 255}
	LightPurple color.Color = color.RGBA{191, 162, 208, 255}
	DarkTeal    color.Color = color.RGBA{5, 83, 107, 255}
)

type Sortable interface {
	Sort(ctx context.Context, arr []int, compare CompareCallback, swap SwapCallback)
}
type CompareCallback func(ctx context.Context, arr []int, n1, n2 int) bool
type SwapCallback func(ctx context.Context, arr []int, n1, n2 int)

// //////////////////////////////////

type SortType int

const (
	SortTypeQuicksort SortType = iota
	SortTypeBubbleSort
	SortTypeMergeSort
	SortTypeInsertionSort
)

// //////////////////////////////////

type SortState int

const (
	SortStateIdle SortState = iota
	SortStateInProgress
	SortStatePaused
)

func (ss SortState) String() string {
	return [...]string{"SortStateIdle", "SortStateInProgress", "SortStatePaused"}[ss]
}

type SortHandler struct {
	arr               []int
	state             SortState
	normalSort        Sortable
	multithreadedSort Sortable
	multithread       bool
}

func (handler *SortHandler) Sort(ctx context.Context, arr []int, cmp CompareCallback, swap SwapCallback) {
	if handler.multithread && handler.multithreadedSort != nil {
		handler.multithreadedSort.Sort(ctx, arr, cmp, swap)
	} else {
		if handler.multithread {
			log.Print("Multithreaded version of Sort not available, using normal sort algo")
		}
		handler.normalSort.Sort(ctx, arr, cmp, swap)
	}
}

// //////////////////////////////////////
type SortCanvas struct {
	Nums                  []int
	activeSort            SortType
	sortStateCommunicator pubsub.PubSub[SortType, SortState]
	SortTypeStateMap      map[SortType]*SortHandler
	colors                []color.Color
	default_color         color.Color
	Position              image.Rectangle
	rect_size             f32.Vec2 //because it has property of a vector
	rect_padding          float32
	start_pos             float32
	sleep_time            time.Duration
}

type InitSortConfig struct {
	Nums         []int
	Pos          image.Rectangle
	DefaultSort  SortType
	DefaultColor color.Color
}

func NewSortCanvas(sortConfig *InitSortConfig) *SortCanvas {
	canv := SortCanvas{}
	canv.SortTypeStateMap = map[SortType]*SortHandler{
		SortTypeQuicksort:     &QuicksortHandler,
		SortTypeBubbleSort:    &BubblesortHandler,
		SortTypeMergeSort:     &MergeSortHandler,
		SortTypeInsertionSort: &InnsertionSortHandler,
	}
	canv.sortStateCommunicator = *pubsub.NewPubSub[SortType, SortState]()
	canv.Nums = sortConfig.Nums
	canv.Position = sortConfig.Pos
	canv.default_color = sortConfig.DefaultColor
	canv.colors = make([]color.Color, len(canv.Nums))
	canv.activeSort = sortConfig.DefaultSort
	for i := range canv.Nums {
		canv.colors[i] = canv.default_color
	}
	canv.sleep_time = 30 * 1000 * time.Microsecond
	canv.UpdateSizes()
	return &canv
}

func (canv *SortCanvas) getCurrSortHandler() *SortHandler {
	return canv.SortTypeStateMap[canv.activeSort]
}

func change_col(delay time.Duration, place *color.Color, col color.Color) {
	time.Sleep(10 * delay)
	*place = col
}

func (canv *SortCanvas) handlePlayPause(ctx context.Context) {
	stateCommChannel := ctx.Value(SortContextKey("CommChan")).(chan SortState)

	if canv.getCurrSortHandler().state == SortStatePaused {
		// log.Print("current state is paused, so waiting to get play signal")
	L:
		for {
			state := <-stateCommChannel
			// log.Printf("current state is paused, got a signal %v", state)
			switch state {
			case SortStateInProgress:
				// log.Print("unpausing the run")
				canv.getCurrSortHandler().state = state
				break L
			default:
				log.Print("can only unpause when paused.")
			}
		}
	}

	select {
	case new_state := <-stateCommChannel:
		// log.Printf("current state: %v |||| received state: %v", canv.getCurrSortHandler().state, new_state)
		if new_state == SortStatePaused && canv.getCurrSortHandler().state == SortStateInProgress {
			canv.getCurrSortHandler().state = SortStatePaused
		} else if new_state == SortStateInProgress && canv.getCurrSortHandler().state == SortStatePaused {
			// log.Print("unpausing the run")
			canv.getCurrSortHandler().state = new_state
		}
	default:
	}

}

func (canv *SortCanvas) compare(ctx context.Context, arr []int, n1, n2 int) bool {
	canv.handlePlayPause(ctx)

	offset := ctx.Value(SortContextKey("offset")).(int)
	canv.colors[n1+offset] = LightBlue
	canv.colors[n2+offset] = LightBlue
	time.Sleep(canv.sleep_time / 2)
	val := arr[n1] < arr[n2]
	time.Sleep(canv.sleep_time / 2)
	go change_col(0*canv.sleep_time, &canv.colors[n1+offset], canv.default_color)
	go change_col(0*canv.sleep_time, &canv.colors[n2+offset], canv.default_color)

	canv.handlePlayPause(ctx)
	return val
}

func (canv *SortCanvas) swap(ctx context.Context, arr []int, n1, n2 int) {

	canv.handlePlayPause(ctx)
	offset := ctx.Value(SortContextKey("offset")).(int)

	canv.colors[n1+offset] = Pink
	canv.colors[n2+offset] = Pink
	time.Sleep(canv.sleep_time / 2)
	arr[n1], arr[n2] = arr[n2], arr[n1]
	time.Sleep(canv.sleep_time / 2)
	go change_col(canv.sleep_time, &canv.colors[n1+offset], canv.default_color)
	go change_col(canv.sleep_time, &canv.colors[n2+offset], canv.default_color)

	canv.handlePlayPause(ctx)
}

func (canv *SortCanvas) PlayPause() {
	switch canv.getCurrSortHandler().state {
	case SortStatePaused:
		log.Printf("publishing start signal to the state communicator  %v", canv.getCurrSortHandler().state)
		canv.sortStateCommunicator.Publish(canv.activeSort, SortStateInProgress)
	case SortStateInProgress:
		log.Printf("publishing pause signal to the state communicator  %v", canv.getCurrSortHandler().state)
		canv.sortStateCommunicator.Publish(canv.activeSort, SortStatePaused)
	}
}

// ////////////////////////////////////
func (canv *SortCanvas) UpdateSizes() {
	pos := canv.Position
	rect_width := float32(pos.Dx()) / float32(len(canv.Nums))
	padding := rect_width / 4
	rect_width = rect_width - padding
	max_rect_height := pos.Dy()
	canv.rect_size = f32.Vec2{float32(rect_width), float32(max_rect_height)}
	canv.rect_padding = padding
	canv.start_pos = float32(pos.Min.X) + padding + float32(canv.rect_size[0])/2

	log.Printf("padding: %f rect_size: %f total: %f", padding, rect_width, padding*float32(len(canv.Nums))+rect_width*float32(len(canv.Nums)))
}

func (canv *SortCanvas) Draw(screen *ebiten.Image) {
	pos := canv.start_pos
	max_num := len(canv.Nums)
	for i, num := range canv.Nums {
		rect := image.Rect(
			int((pos - float32(canv.rect_size[0]/2))),
			canv.Position.Min.Y+(canv.Position.Dy()-int(canv.rect_size[1])*num/max_num),
			int((pos + float32(canv.rect_size[0]/2))),
			canv.Position.Max.Y,
		)

		vector.DrawFilledRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), canv.colors[i], false)
		pos += canv.rect_size[0] + canv.rect_padding
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("current sort state: %v", canv.getCurrSortHandler().state))

	// vector.StrokeRect(screen, float32(canv.Position.Min.X), float32(canv.Position.Min.Y), float32(canv.Position.Dx()), float32(canv.Position.Dy()), 1, color.White, true)
}

type SortContextKey string

func (canv *SortCanvas) Start() {
	if canv.getCurrSortHandler().state == SortStateIdle {
		canv.getCurrSortHandler().arr = canv.Nums
		ctx := context.WithValue(context.Background(), SortContextKey("offset"), 0)

		subscribe := func() chan SortState {
			if comm_chan, err := canv.sortStateCommunicator.Subscribe(canv.activeSort, 1); err == nil {
				return comm_chan
			} else {
				log.Panicf("can't create channel for a subscriber. | error: %v", err)
				return nil
			}
		}

		unsubscribe := func(commChan chan SortState) {
			canv.sortStateCommunicator.Unsubscribe(canv.activeSort, commChan)
		}

		ctx = context.WithValue(ctx, SortContextKey("Subscribe"), subscribe)
		ctx = context.WithValue(ctx, SortContextKey("Unsubscribe"), unsubscribe)

		go func() {
			comm_chan := subscribe()
			ctx = context.WithValue(ctx, SortContextKey("CommChan"), comm_chan)
			defer unsubscribe(comm_chan)
			canv.getCurrSortHandler().Sort(ctx, canv.getCurrSortHandler().arr, canv.compare, canv.swap)
		}()

	} else if canv.getCurrSortHandler().state == SortStatePaused {
		canv.PlayPause()
	}
}
