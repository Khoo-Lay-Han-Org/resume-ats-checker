package main

import "resuming/api"

func main() {
	router := api.RunOpenAPIDoc()
	if err := router.Run(":5781"); err != nil {
		panic(err)
	}
}
