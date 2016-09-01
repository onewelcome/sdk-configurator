default: prepareLibs

prepareLibs:
	go-bindata -pkg data -o data/bindata.go lib/

install: prepareLibs
	go install