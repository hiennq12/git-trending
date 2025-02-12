package utils

import (
	"fmt"
	"time"
)

func GetToday() string {
	// Lấy thời gian hiện tại
	now := time.Now()

	// Định dạng ngày tháng năm
	date := now.Format("02/01/2006") // dd/mm/yyyy

	// Lấy thứ trong tuần (tiếng Anh)
	weekday := now.Weekday()

	// In ra kết quả
	return fmt.Sprintf("Today is: %s, date %s\n", weekday, date)
}
