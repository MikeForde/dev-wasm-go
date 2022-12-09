# Devcontainer WASM-Go
Simple devcontainer for Go development

# Usage

## Github Codespaces
Just click the button:

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new?hide_repo_select=true&ref=main&repo=575629782)



## Visual Studio Code
Note this assumes that you have the VS code support for remote containers and `docker` installed 
on your machine.

```sh
git clone https://github.com/dev-wasm/dev-wasm-go
cd dev-wasm-go
code ./
```

Visual studio should prompt you to see if you want to relaunch the workspace in a container, you do.

# Building and Running

```sh
tinygo build -wasm-abi=generic -target=wasi -o main.wasm main.go
wasmtime main.wasm --dir .
```
