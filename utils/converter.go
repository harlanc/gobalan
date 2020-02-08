package utils

import "strconv"

//String2Int convert string to int
func String2Int(val string) int {

	i, err := strconv.Atoi(val)
	if err != nil {
		return -1
	}
	return i

}

//Int2String convert int to string
func Int2String(val int) string {
	return strconv.Itoa(val)
}

//Int322String convert int32 to string
func Int322String(val int32) string {
	return strconv.FormatInt(int64(val), 10)
}

//Int642String convert int64 to string
func Int642String(val int64) string {
	return strconv.FormatInt(val, 10)
}

//Float642String convert float64 to string
func Float642String(val float64) string {
	return strconv.FormatFloat(val, 'E', -1, 64)
}
