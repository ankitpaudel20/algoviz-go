package sorting

var InnsertionSortHandler SortHandler

func init() {
	InnsertionSortHandler = SortHandler{state: SortStateIdle}
}
