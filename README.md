blocksafari
===========

blocksafari is a web-based frontend to [btcd](https://github.com/conformal/btcd).

**blocksafari is proof-of-concept code only.**

## Installation

#### Build from source

- Install Go according to the installation instructions here:
  http://golang.org/doc/install

- Run the following command to obtain blocksafari, all dependencies, and install it:
  ```$ go get github.com/conformal/blocksafari```

- Enter the source directory
  ```cd $GOPATH/src/github.com/conformal/blocksafari```

- Copy sample-blocksafari.conf to blocksafari.conf and edit the options.

- Start blocksafari:
  ```blocksafari -C blocksafari.conf```

## Updating

#### Build from Source

- Run the following command to update blocksafari, all dependencies, and install it:
  ```$ go get -u -v github.com/conformal/blocksafari/...```

## License

blocksafari is licensed under the liberal ISC License.
