survana [![Build Status](https://travis-ci.org/vpetrov/survana.png?branch=1.0)](https://travis-ci.org/vpetrov/survana)
=======

An HTML5 application for administering questionnaires on tablets and mobile devices. Developed by the Neuroinformatics Research Group at Harvard University.

Building
========

`make` or `make osx`

Prerequisites:

  * git, hg, bzr (for downloading Go libraries with `go get`)
  * make
  * go 1.2+

All binaries and apps will be placed in `bin`.

To build just the server, type `make`.

On OS X, you can type `make osx` instead. It builds `bin/server` and `bin/Survana`, which is a status bar application for managing the server.
