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