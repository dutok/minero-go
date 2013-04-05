Minero
======

Minero is an implementation of the Multiplayer Server for [Minecraft](http://minecraft.net) made in [Go](http://golang.org). It aims to fully support Minecraft 1.4.5 version.

It is licensed under the MIT open source license, please see the [LICENSE.txt](https://github.com/toqueteos/minero/blob/master/LICENSE.txt) file for more information.

Requirements
============

Right now it's aimed for: `go version go1.0.3`.

You can check your Go version typing `go version` on the terminal. If it outputs an error you don't have Go installed.

Go to [Go's install page](http://golang.org/doc/install) **Download the Go tools** section and follow the instructions.

Features
========

- NBT v19133 support.
    - nbtdebug command on [`bin/nbtdebug`](https://github.com/toqueteos/minero/blob/master/bin/nbtdebug).
- Basic [data types](http://wiki.vg/Data_Types) support (bool, byte, short, int, long, float, double and string). See [`types`](https://github.com/toqueteos/minero/blob/master/types), [`types/nbt`](https://github.com/toqueteos/minero/blob/master/types/nbt) and [`types/minecraft`](https://github.com/toqueteos/minero/blob/master/types/minecraft).
- Proxy with logging support available.
    - miproxy command on [`bin/miproxy`](https://github.com/toqueteos/minero/blob/master/bin/miproxy).

Instructions
============

Commands can be installed with `go install`.

- Copy & Paste fans: `go get github.com/toqueteos/minero/bin/<cmdName>`
- Working example: `go get github.com/toqueteos/minero/bin/nbtdebug`
