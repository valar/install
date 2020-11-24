package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func increaseCounter(reptype string) {
	req, err := http.NewRequest("POST", "https://kv.valar.dev/valar/"+reptype+"?op=inc", nil)
	if err != nil {
		log.Println("request invalid:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("VALAR_TOKEN"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("report failed:", err)
		return
	}
	resp.Body.Close()
}

func main() {
	// Load script
	script, err := ioutil.ReadFile("install.sh")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		repType := r.URL.Query().Get("type")
		switch repType {
		case "install_finish":
			increaseCounter("total_installs")
		case "install_start":
			increaseCounter("started_installs")
		}
		w.Write([]byte(""))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(script)
		increaseCounter("total_downloads")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
