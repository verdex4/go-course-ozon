package main

import (
	"sync"

	cart "github.com/verdex4/go-course-ozon/test-concurrency/cmd/cart-service"
	product "github.com/verdex4/go-course-ozon/test-concurrency/cmd/product-service"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		product.StartProductService()
	}()
	go func() {
		defer wg.Done()
		cart.StartCartService()
	}()
	wg.Wait()
}
