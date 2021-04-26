package http_util


func CalcUtf8TextWidth(text string) int {
	unicodeArr := []rune(text)
	num := 0
	for i := 0; i < len(unicodeArr); i++ {
		// 一个中文字符显示宽度约等于7/4个英文字符，一个英文字符当作3的宽度，则一个中文字符
		// 的宽度即为5，最终返回的值/3即为英文字符数
		if unicodeArr[i] > 127 {
			num += 5
		} else {
			num += 3
		}
	}
	return num
}

/**
计算文本占据多少个tab的宽度
 */
func CalcTabNum(text string) int {
	num := CalcUtf8TextWidth(text)
	// 4个英文字符占据一个tab
	return  num / 12
}

