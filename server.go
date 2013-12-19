package main 

import (
	"os"
	"fmt"
	"net/http"
)

func main() {
	var port = os.Getenv("PORT")

	if len(port) == 0 {
		fmt.Println("$PORT not set -- defaulting to 5000", port)
		port = "5000"
	}
	fmt.Println("using port:", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, GO")
}
