package entity

type Comparable interface {
	GetValue() int64
	// 和另一个对象比较，返回值<0 表明当前对象小于another
	// 返回值=0，当前对象等于another
	// 返回值>0，当前对象大于another
	Compare(another Comparable) int
}
