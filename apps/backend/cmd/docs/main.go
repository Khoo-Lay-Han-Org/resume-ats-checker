package main

import "resuming/controller"

func main() {
	router := api.RunOpenAPIDoc()
	if err := router.Start(":5781"); err != nil {
		panic(err)
	}
}
