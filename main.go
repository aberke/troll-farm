package main 

import (
	"os"
	"fmt"
	"log"
	"net/http"

    "github.com/aberke/troll-farm/trolls"
)

func main() {

	var port = os.Getenv("PORT")
	if len(port) == 0 {
		fmt.Println("$PORT not set -- defaulting to 5000", port)
		port = "5000"
	}
	fmt.Println("using port:", port)


	// troll websocket server
	trollServer := trolls.NewServer()
	go trollServer.Listen()
	
	http.HandleFunc("/static/", serveStatic)
	http.HandleFunc("/", serveHome)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
    }
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/")
}
func serveStatic(w http.ResponseWriter, r *http.Request) {
	var staticFileHandler = http.FileServer(http.Dir("./public/"))
	staticFileHandler.ServeHTTP(w, r)
}
