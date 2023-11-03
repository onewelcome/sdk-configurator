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
go get -u github.com/spf13/cobra
go install github.com/jteeuwen/go-bindata/...@latest
```

Clone project:
```sh
go install github.com/onewelcome/sdk-configurator
```

Initialize module:
```sh
go mod init sdk-configurator
```

Install dependencies:
```sh
go mod tidy
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