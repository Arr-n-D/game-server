package internal

import (
	"fmt"
	"time"
)

func InitLogger() {
	// get the current date in yyyy-mm-dd h:m:s in UTC
	date := time.Now()

	fmt.Println(date)
}
