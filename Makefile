default: prepareLibs

prepareLibs:
	go-bindata -pkg data -o data/bindata.go lib/

install: prepareLibs
	go install

clean:
	git checkout data/bindata.go
	rm -rf target

release: prepareLibs
	go get github.com/mitchellh/gox
	gox -os="darwin linux windows" -arch="386 amd64" -output="target/{{.OS}}/{{.Arch}}/{{.Dir}}"
	zip target/onegini-sdk-configurator-macos-64bit.zip -j target/darwin/amd64/onegini-sdk-configurator
	zip target/onegini-sdk-configurator-linux-32bit.zip -j target/linux/386/onegini-sdk-configurator
	zip target/onegini-sdk-configurator-linux-64bit.zip -j target/linux/amd64/onegini-sdk-configurator
	zip target/onegini-sdk-configurator-windows-32bit.zip -j target/windows/386/onegini-sdk-configurator.exe
	zip target/onegini-sdk-configurator-windows-64bit.zip -j target/windows/amd64/onegini-sdk-configurator.exe