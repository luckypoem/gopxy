package main

import "gopxy"

func main() {
	svr := gopxy.New(nil)
	err := svr.Start()
	if err != nil {
		panic(nil)
	}
}