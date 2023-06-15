package utils

import (
	"fmt"
	"time"
)

var stopAnimation chan bool

func StartLoadingAnimation(url string) {
	stopAnimation = make(chan bool)

	dots := []string{"   ", ".  ", ".. ", "..."}
	currentDot := 0

	for {
		select {
		case <-stopAnimation:
			return
		default:
			fmt.Printf("\rCrawling at %s %s", url, dots[currentDot])
			currentDot = (currentDot + 1) % len(dots)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func StopLoadingAnimation() {
	stopAnimation <- true
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("\r\nCrawling complete\n")
}
