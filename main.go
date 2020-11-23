package main

import (
	"log"
	"net/http"
	"os"
)

var script = []byte(`#!/bin/bash

# Report install start
curl -sSL "https://cli.valar.dev/report?type=install_start"

arch=$(uname -m)
os=$(uname -s)

if [ "$os" == "Darwin" ] && [ "$arch" == "x86_64" ]; then
    sudo curl -sSL -o /usr/local/bin/valar https://github.com/valar/cli/releases/latest/download/valar_darwin_amd64
elif [ "$os" == "Linux" ] && [ "$arch" == "x86_64" ]; then
    sudo curl -sSL -o /usr/local/bin/valar https://github.com/valar/cli/releases/latest/download/valar_linux_amd64
elif [ "$os" == "Linux" ] && [ "${arch:0:3}" == "arm" ]; then
    sudo curl -sSL -o /usr/local/bin/valar https://github.com/valar/cli/releases/latest/download/valar_linux_arm
else
    echo "Unsupported OS/ARCH $arch/$os"
fi

sudo chmod +x /usr/local/bin/valar

if [ ! -f "$HOME/.valar/valarcfg" ]; then
    echo "Configuring Valar ..."
    mkdir -p $HOME/.valar
    printf "Token: "
    read -s API_TOKEN
    cat - > $HOME/.valar/valarcfg <<EOF
token: $API_TOKEN
endpoint: https://api.valar.dev/v1
EOF
fi

# Report install finish
curl -sSL "https://cli.valar.dev/report?type=install_finish"

echo "Valar is now installed on your machine. Enjoy :)"
`)

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
