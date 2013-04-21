Minero
======

Minero is an implementation of the Multiplayer Server for [Minecraft](http://minecraft.net) made in [Go](http://golang.org). It aims to fully support Minecraft 1.5.1 version.

It is licensed under the MIT open source license, please read the [LICENSE.txt](https://github.com/toqueteos/minero/blob/master/LICENSE.txt) file for more information.

Requirements
============

Just Go, also Git (encouraged) if you want to use `go get`.

More specifically aimed for: `go version go1.0.3`.

You can check your Go version typing `go version` on the terminal. If it outputs an error you don't have Go installed.

Go to [Go's install page](http://golang.org/doc/install) **Download the Go tools** section and follow the instructions.

Features
========

- Basic [data types](http://wiki.vg/Data_Types) support (bool, byte, short, int, long, float, double and string). See [`types`](https://github.com/toqueteos/minero/blob/master/types), [`types/nbt`](https://github.com/toqueteos/minero/blob/master/types/nbt) and [`types/minecraft`](https://github.com/toqueteos/minero/blob/master/types/minecraft).
- NBT v19133 support.
- Proxy with logging support available.
- Server list ping client & server (ping other servers, fake a server).

Tools
=====
- Minero server: [`bin/minero`](https://github.com/toqueteos/minero/blob/master/bin/minero)

        go get github.com/toqueteos/minero/bin/minero

- NBT pretty printer: [`bin/minbtd`](https://github.com/toqueteos/minero/blob/master/bin/minbtd)

        go get github.com/toqueteos/minero/bin/minbtd

- Server proxy with logging support: [`bin/miproxy`](https://github.com/toqueteos/minero/blob/master/bin/miproxy)

        go get github.com/toqueteos/minero/bin/miproxy

- Server list ping client & server: [`bin/mipingd`](https://github.com/toqueteos/minero/blob/master/bin/mipingd)

        go get github.com/toqueteos/minero/bin/mipingd

Notes
=====

Everything can be go-get'd.
