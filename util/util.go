package util

import "strconv"

// SaveBit 保留小数位，四舍五入
func SaveBit(f float64, bit int) float64 {
	str := strconv.FormatFloat(f, 'f', bit, 64)
	float, _ := strconv.ParseFloat(str, 64)
	return float
}
