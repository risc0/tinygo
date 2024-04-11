package main

import "time"

var timer = time.NewTimer(time.Millisecond)

func main() {
	// Test ticker.
	ticker := time.NewTicker(time.Nanosecond * 5)
	println("waiting on ticker")
	go func() {
		time.Sleep(time.Nanosecond * 5)
		println(" - after 150ms")
		time.Sleep(time.Nanosecond * 10)
		println(" - after 200ms")
		time.Sleep(time.Nanosecond * 20)
		println(" - after 300ms")
	}()
	<-ticker.C
	println("waited on ticker at 500ms")
	<-ticker.C
	println("waited on ticker at 1000ms")
	ticker.Stop()
	time.Sleep(time.Nanosecond * 7)
	select {
	case <-ticker.C:
		println("fail: ticker should have stopped!")
	default:
		println("ticker was stopped (didn't send anything after 750ms)")
	}

	timer := time.NewTimer(time.Nanosecond * 5)
	println("waiting on timer")
	go func() {
		time.Sleep(time.Nanosecond * 2)
		println(" - after 200ms")
		time.Sleep(time.Nanosecond * 4)
		println(" - after 400ms")
	}()
	<-timer.C
	println("waited on timer at 500ms")
	time.Sleep(time.Nanosecond * 5)

	reset := timer.Reset(time.Nanosecond * 5)
	println("timer reset:", reset)
	println("waiting on timer")
	go func() {
		time.Sleep(time.Nanosecond * 2)
		println(" - after 200ms")
		time.Sleep(time.Nanosecond * 4)
		println(" - after 400ms")
	}()
	<-timer.C
	println("waited on timer at 500ms")
	time.Sleep(time.Nanosecond * 5)
}
