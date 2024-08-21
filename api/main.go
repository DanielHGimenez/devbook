package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

// func init() {
// 	chave := make([]byte, 64)
// 	rand.Read(chave)
// 	fmt.Println(base64.StdEncoding.EncodeToString(chave))
// }

// func init() {
// 	var x, _ = security.Hash("password123")
// 	fmt.Println(string(x))
// }

func main() {
	config.Load()
	r := router.Build()

	log.Println("Listening on server port", config.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), r))
}
