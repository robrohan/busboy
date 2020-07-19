# Busboy

Busboy is a simple HTTP server used when testing things locally. Several newer internet things (tm) require a file to be loaded over http. This can be a pain when you just want to test something quickly - when doing a spike or something.

## Building

Busboy is written in golang, so if you have a golang environment setup, it's just a matter of building the application. There is nothing fancy with it.

    $ make 

or

    $ go build cmd/busboy/main.go

## Usage

Just running it...

    $ ./busyboy

will start a server on port :8580 and serve the current directory. You can do:

    $ ./busboy --help

For options.

I have a alias in my shell to just `bb`.

    alias bb=/home/rob/Projects/busboy/build/busboy
