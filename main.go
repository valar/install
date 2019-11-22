package main

import (
	"log"
	"net/http"
)

var script = []byte(`#!/bin/bash

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

if [ ! -f "$HOME/.valar/valar.cloud.yml" ]; then
    echo "Configuring Valar Cloud ..."
    mkdir -p $HOME/.valar
    printf "Token: "
    read -s API_TOKEN
    cat - > $HOME/.valar/valar.cloud.yml <<EOF
token: $API_TOKEN
endpoint: https://api.valar.dev
EOF
fi

echo "Valar is now installed on your machine. Enjoy :)"
`)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(script)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
