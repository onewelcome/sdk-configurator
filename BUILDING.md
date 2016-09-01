## Building from source

>**Warning:** Only build the SDK configurator from source if you must modify it. The master version may not be compatible with the Onegini SDK version that you 
are using.

Make sure you have go installed:
```sh
brew install go
```

Read more about setting up go in the [official docs](https://golang.org/doc/install)

Install dependencies:
```sh
go get github.com/spf13/cobra
go get -u github.com/jteeuwen/go-bindata/...
```

Clone project:
```sh
go get github.com/Onegini/onegini-sdk-configurator
```

Build project with:
```sh
make
```

Install the go binary with:
```sh
make install
```

Or run without export a binary using:
```sh
go run main.go
```

Update binary assets using
```sh
go-bindata -pkg data -o data/bindata.go lib/
```