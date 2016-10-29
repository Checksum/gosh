# gosh
Simple HTTP file server for development with browser auto reload.

## Why yet another server?

Needed a single binary which I can use on remote servers without having to install Python or a million Node.js dependencies.

## Installation

Download the binary for your OS from [releases](https://github.com/Checksum/gosh/releases). Put it in some path where it's accessible.

## Usage

`gosh` serves the current directory on port `8000`. By default, it doesn't watch for changes. To watch for changes:

`gosh -watch -ext="html,css,js"`.  

`gosh -help` lists all available options.

## Browser support

If watching is enabled, the server intercepts all requests with `Accept` header `text/html` and injects a script which opens a `EventSource` connection to listen for changes. Since this requires `Server-sent events`, IE/Edge users are out of luck! Check if your [browser is supported](http://caniuse.com/#feat=eventsource)

