package sorting

import "context"

var MergeSortHandler SortHandler

func init() {
	MergeSortHandler = SortHandler{state: SortStateIdle, multithread: true}
}

type MergeSort struct{}

func (m *MergeSort) Merge(ctx context.Context, arr []int, pivot int, compare CompareCallback, swap SwapCallback, offset int) {
	p1, p2 := 0, pivot
	// final_arr := make([]int, len(arr1))
	for p1 < pivot || pivot < len(arr) {
		if arr[p2] < arr[p1] {
			swap(ctx, arr, p1, p2)

		}
	}
}
func (m *MergeSort) Sort(ctx context.Context, arr []int, compare CompareCallback, swap SwapCallback) {
	if len(arr) == len(MergeSortHandler.arr) {
		MergeSortHandler.state = SortStateInProgress
	}
	if len(arr) < 1 {
		return
	}

	// pivot := len(arr) / 2
	// m.Sort(ctx, arr[:pivot], compare, swap, offset)
	// m.Sort(ctx, arr[pivot:], compare, swap, offset+pivot)

}
