package condition_variable

import (
	"fmt"
	"sync"
	"testing"
)

func TestConditionVariable(t *testing.T) {
	data := NewData()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		data.WaitForReady()
	}()

	fmt.Println("Preparing data...")
	data.SetReady()

	wg.Wait()
}
