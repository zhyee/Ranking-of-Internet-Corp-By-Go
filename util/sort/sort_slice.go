package sort

type Comparable interface {
	GetValue() int64
	// 和另一个对象比较，返回值<0 表明当前对象小于another
	// 返回值=0，当前对象等于another
	// 返回值>0，当前对象大于another
	Compare(another Comparable) int
}

func FindInsertIndex(array []Comparable, startIdx, endIdx int, target Comparable) int {

	if startIdx > endIdx {
		return 0
	}

	midIdx := (startIdx + endIdx) / 2

	if target.Compare(array[midIdx]) > 0 {
		if midIdx >= endIdx {
			return endIdx + 1
		}
		return FindInsertIndex(array, midIdx + 1, endIdx, target)
	} else if target.Compare(array[midIdx]) < 0 {
		if midIdx <= startIdx {
			return startIdx
		}
		return FindInsertIndex(array, startIdx, midIdx - 1, target)
	} else {
		return midIdx+1
	}

}

func Insert(array []Comparable, target Comparable) []Comparable {

	if array == nil {
		array = make([]Comparable, 0)
	}

	idx := FindInsertIndex(array, 0, len(array) - 1, target)

	array = append(array, target)

	pos := len(array) - 2

	for pos >= idx {
		array[pos+1] = array[pos]
		pos--
	}

	array[idx] = target

	return array
}
