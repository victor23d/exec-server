package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type test struct {
	message string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "get.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		postMessage := r.PostForm["message"]
		message := os.Getenv("MESSAGE")
		if postMessage[0] == message {
			// if message != nil {

			fmt.Fprintln(w, "Authenticate succeed ! ")
			fmt.Fprintln(w, "Wait 10 seconds for build and deploy ...")

			result := execCommand()
			fmt.Fprintln(w, result)
			log.Println(time.Now().Format(time.RFC850))
		} else {

			log.Println(time.Now().Format(time.RFC850) + "     " + postMessage[0])
		}
		// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	default:
		fmt.Fprintf(w, "Site is closed.")
	}

}

func testHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t test

	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t.message)
}

func execCommand() string {
	cmd := exec.Command("/root/run.sh")
	// cmd := exec.Command("bash", "-c", "/root/run.sh")
	// cmd.Stdin = Ntrings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

func main() {
	http.HandleFunc("/", viewHandler)
	log.Fatal(http.ListenAndServe(":3222", nil))
}
