package core

import "time"

type queueProcesses []*goLoop

func (q *queueProcesses) Add() {}

type goLoop struct {
	duration time.Duration
	ticker   *time.Ticker
	function func() bool // aun no he implementado goloops ????
}

func _SetStateProcesses(newState bool) {
	isHung = newState
	if newState {
		println("---- (WIN) ---> Alert: Hunging, stopping processes")
		for _, ll := range processes {
			ll.ticker.Stop()
		}
	} else {
		println("---- (WIN) ---> Alert: Runnig restarting Processes")
		for _, ll := range processes {
			ll.ticker.Reset(ll.duration)
		}
	}
}
