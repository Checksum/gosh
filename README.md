# gosh
Simple HTTP file server for development with browser auto reload

## Why yet another server?

Needed a single binary which I can use on remote servers without having to install Python or a million Node.js dependencies.

## Installation

Download the binary for your OS from [releases](https://github.com/Checksum/gosh/releases). Put it in some path where it's accessible.

## Usage

`gosh` serves the current directory on port `8000`. By default, it doesn't watch for changes. To watch for changes, use `gosh -watch -ext="html,css,js"`.  

`gosh -help` lists all available options.