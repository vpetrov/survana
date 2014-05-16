survana [![Build Status](https://travis-ci.org/vpetrov/survana.png?branch=1.0)](https://travis-ci.org/vpetrov/survana)
=======

An HTML5 application for administering questionnaires on tablets and mobile devices. Developed by the Neuroinformatics Research Group at Harvard University.

Download
========

To download pre-built binaries, go to https://github.com/vpetrov/survana/releases

Building From Source
====================

`git submodule init`
`git submodule update`
`make` or `make osx`

Prerequisites:

  * git, hg, bzr (for downloading Go libraries with `go get`)
  * make
  * go 1.2.1+

All binaries and apps will be placed in `bin`.

To build just the server, type `make`.

On OS X, you can type `make osx` instead. It builds `bin/server` and `bin/Survana`, which is a status bar application for managing the server.
