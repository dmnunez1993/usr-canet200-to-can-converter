package usrcanettocan

import "sync"

func LoopForever() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
