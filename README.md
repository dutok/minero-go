IMPORTANT NOTE
===============

I've seen some activity related to minero so I'm adding this little note to README.

Project is *NOT ACTIVE* anymore until further notice.

- **Why?** I have less time now to tinker with this project also Mojang is taking like forever to remake Minecraft properly from scratch so there are piles and piles of shit in the server code I must reimplement again and I refuse to do that.
- **Will that change?** If they remake the game properly I'll continue with minero.

If you want to contact me you can try any of these:
1. Create a new issue here
2. Drop me an email on: `toqueteos at gmail dot com`.

If you are still interested in Go + Minecraft check this out: http://github.com/NetherrackDev

Minero
======

Minero is an implementation of the Multiplayer Server for [Minecraft](http://minecraft.net) made in [Go](http://golang.org). It aims to fully support Minecraft 1.5.1 version.

It is licensed under the MIT open source license, please read the [LICENSE.txt](https://github.com/minero/minero/blob/master/LICENSE.txt) file for more information.

Requirements
============

Just Go, also Git (encouraged) if you want to use `go get`.

More specifically aimed for: `go version go1.0.3`. 

You can check your Go version typing `go version` on the terminal. If it outputs an error you don't have Go installed.

Go to [Go's install page](http://golang.org/doc/install) **Download the Go tools** section and follow the instructions.

**NOTE:** It should work with newer versions of Go, go1.1.1 right now.

Features
========

- Basic [data types](http://wiki.vg/Data_Types) support (bool, byte, short, int, long, float, double and string). See [`types`](https://github.com/minero/minero/blob/master/types), [`types/nbt`](https://github.com/minero/minero/blob/master/types/nbt) and [`types/minecraft`](https://github.com/minero/minero/blob/master/types/minecraft).
- NBT v19133 support.
- Proxy with logging support available.
- Server list ping client & server (ping other servers, fake a server).

**NOTE:** Right now the server allows client to log in, move and do most actions but any changes on map aren't saved.

Tools
=====
- Minero server: [`bin/minero`](https://github.com/minero/minero/blob/master/bin/minero)

        go get github.com/minero/minero/bin/minero

- NBT pretty printer: [`bin/minbtd`](https://github.com/minero/minero/blob/master/bin/minbtd)

        go get github.com/minero/minero/bin/minbtd

- Server proxy with logging support: [`bin/miproxy`](https://github.com/minero/minero/blob/master/bin/miproxy)

        go get github.com/minero/minero/bin/miproxy

- Server list ping client & server: [`bin/mipingd`](https://github.com/minero/minero/blob/master/bin/mipingd)

        go get github.com/minero/minero/bin/mipingd

Notes
=====

Everything can be go-get'd.
