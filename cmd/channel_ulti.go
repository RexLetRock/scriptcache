package cmd

import (
	"fmt"
	"time"
)

func showResult() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				timeNow = time.Now().Unix()
				total := 0
				for i, v := range ACount {
					if v != 0 {
						total += v
						ACount[i] = 0
					}
				}
				stotal += total
				fmt.Printf("Threadnum %v - Msg/s %v - Msg %v \n", NCpu, commaize(int(float64(total)/float64(timeNow-timeStart))), commaize(stotal))
				timeStart = timeNow
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func commaize(n int) string {
	s1, s2 := fmt.Sprintf("%d", n), ""
	for i, j := len(s1)-1, 0; i >= 0; i, j = i-1, j+1 {
		if j%3 == 0 && j != 0 {
			s2 = "," + s2
		}
		s2 = string(s1[i]) + s2
	}
	return s2
}
