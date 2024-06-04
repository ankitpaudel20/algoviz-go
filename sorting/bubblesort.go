package sorting

import "context"

var BubblesortHandler SortHandler

func init() {
	BubblesortHandler = SortHandler{state: SortStateIdle, normalSort: &Bubblesort{}, multithread: true}
}

type Bubblesort struct{}

func (q *Bubblesort) Sort(ctx context.Context, arr []int, compare CompareCallback, swap SwapCallback) {
	BubblesortHandler.state = SortStateInProgress

	for i := len(arr) - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			if compare(ctx, arr, j+1, j) {
				swap(ctx, arr, j+1, j)
			}
		}
	}
	BubblesortHandler.state = SortStateIdle
}
