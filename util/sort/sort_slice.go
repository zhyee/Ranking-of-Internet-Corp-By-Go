package sort

import "Ranking-of-Internet-Corp-By-Go/entity"

func FindInsertIndex(array []entity.Comparable, startIdx, endIdx int, target entity.Comparable) int {

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

func Insert(array []entity.Comparable, target entity.Comparable) []entity.Comparable {

	if array == nil {
		array = make([]entity.Comparable, 0)
	}

	idx := FindInsertIndex(array, 0, len(array) - 1, target)

	array = append(array, target)
	// 插入的位置不是在slice末尾, 则把插入位置及之后的元素往后挪一位
	if idx < len(array) - 1 {
		copy(array[idx+1:], array[idx:])
		array[idx] = target
	}

	return array
}
