#!/bin/bash

# Report install start
curl -sSL "https://cli.valar.dev/report?type=install_start"

arch=$(uname -m)
os=$(uname -s)

if [ "$os" == "Darwin" ] && [ "$arch" == "x86_64" ]; then
    sudo curl -sSL -o /usr/local/bin/valar https://github.com/valar/cli/releases/latest/download/valar_darwin_amd64
elif [ "$os" == "Darwin" ] && [ "$arch" == "arm64" ]; then
    sudo curl -sSL -o /usr/local/bin/valar https://github.com/valar/cli/releases/latest/download/valar_darwin_arm64
elif [ "$os" == "Linux" ] && [ "$arch" == "x86_64" ]; then
    sudo curl -sSL -o /usr/local/bin/valar https://github.com/valar/cli/releases/latest/download/valar_linux_amd64
elif [ "$os" == "Linux" ] && [ "$arch" == "arm64" ]; then
    sudo curl -sSL -o /usr/local/bin/valar https://github.com/valar/cli/releases/latest/download/valar_linux_arm64
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
