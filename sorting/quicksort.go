package sorting

import (
	"context"
	"sync"
)

var QuicksortHandler SortHandler

func init() {
	QuicksortHandler = SortHandler{state: SortStateIdle, multithreadedSort: &QuicksortMulti{}, multithread: true}
}

type QuicksortMulti struct{}

func unsub(ctx context.Context) {
	stateCommChannel := ctx.Value(SortContextKey("CommChan")).(chan SortState)
	unsubscribe := ctx.Value(SortContextKey("Unsubscribe")).(func(chan SortState))
	unsubscribe(stateCommChannel)
}

func (q *QuicksortMulti) Sort(ctx context.Context, arr []int, compare CompareCallback, swap SwapCallback) {

	if len(arr) == len(QuicksortHandler.arr) {
		QuicksortHandler.state = SortStateInProgress
	}

	pivot := 0
	if len(arr) < 3 {
		if len(arr) > 1 && !compare(ctx, arr, 0, 1) {
			swap(ctx, arr, 1, 0)
		}
		return
	}
	lp := 1
	rp := len(arr) - 1

	for lp <= rp {
		if compare(ctx, arr, lp, pivot) {
			lp++
			continue
		} else if compare(ctx, arr, rp, pivot) {
			swap(ctx, arr, lp, rp)
			lp++
		}
		rp--
	}
	swap(ctx, arr, pivot, lp-1)

	// Need to do this here too because'
	// channel of this context will block the publish and the channel will not be used below this
	unsub(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		subscribe := ctx.Value(SortContextKey("Subscribe")).(func() chan SortState)
		ctx1 := context.WithValue(ctx, SortContextKey("CommChan"), subscribe())
		defer unsub(ctx1)
		q.Sort(ctx1, arr[:lp-1], compare, swap)
		wg.Done()
	}()

	go func() {
		subscribe := ctx.Value(SortContextKey("Subscribe")).(func() chan SortState)
		ctx2 := context.WithValue(ctx, SortContextKey("CommChan"), subscribe())
		ctx2 = context.WithValue(ctx2, SortContextKey("offset"), ctx.Value(SortContextKey("offset")).(int)+lp)
		defer unsub(ctx2)
		q.Sort(ctx2, arr[lp:], compare, swap)
		wg.Done()
	}()

	wg.Wait()

	if len(arr) == len(QuicksortHandler.arr) {
		QuicksortHandler.state = SortStateIdle
	}
}

type QuicksortSingle struct{}

func (q *QuicksortSingle) Sort(ctx context.Context, arr []int, compare CompareCallback, swap SwapCallback) {
	if len(arr) == len(QuicksortHandler.arr) {
		QuicksortHandler.state = SortStateInProgress
	}

	pivot := 0
	if len(arr) < 3 {
		if len(arr) > 1 && !compare(ctx, arr, 0, 1) {
			swap(ctx, arr, 1, 0)
		}
		return
	}
	lp := 1
	rp := len(arr) - 1

	for lp <= rp {
		if compare(ctx, arr, lp, pivot) {
			lp++
			continue
		} else if compare(ctx, arr, rp, pivot) {
			swap(ctx, arr, lp, rp)
			lp++
		}
		rp--
	}
	swap(ctx, arr, pivot, lp-1)

	q.Sort(ctx, arr[:lp-1], compare, swap)
	ctx = context.WithValue(ctx, SortContextKey("offset"), ctx.Value(SortContextKey("offset")).(int)+lp)
	q.Sort(ctx, arr[lp:], compare, swap)

	if len(arr) == len(QuicksortHandler.arr) {
		QuicksortHandler.state = SortStateIdle
		stateCommChannel := ctx.Value(SortContextKey("CommChan")).(chan SortState)
		unsubscribe := ctx.Value(SortContextKey("Unsubscribe")).(func(chan SortState))
		unsubscribe(stateCommChannel)
	}
}
