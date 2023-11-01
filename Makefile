default: prepareLibs

prepareLibs:
	bin/go-bindata -pkg data -o data/bindata.go lib/

install: prepareLibs
	go install

clean:
	git checkout data/bindata.go
	rm -rf target

release: prepareLibs
	go install github.com/mitchellh/gox@latest
	bin/gox -os="darwin linux windows" -arch="amd64 arm64" -output="target/{{.OS}}/{{.Arch}}/{{.Dir}}"
	zip target/sdk-configurator-macos-64bit.zip -j target/darwin/amd64/sdk-configurator
	zip target/sdk-configurator-linux-64bit.zip -j target/linux/amd64/sdk-configurator
	zip target/sdk-configurator-windows-64bit.zip -j target/windows/amd64/sdk-configurator.exe
