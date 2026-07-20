package main

import "resuming/api"

func main() {
	router := api.RunOpenAPIDoc()
	if err := router.Start(":5781"); err != nil {
		panic(err)
	}
}
