The following steps must be executed to create a release

1. Update the release notes
2. Update the version constant in cmd/version.go
3. Update package bindata with `go-bindata -pkg data -o data/bindata.go lib`
4. Commit & push to GitHub
5. Create a new tag and push
6. Update the release notes in the release
7. Build the configurator binaries for OS X, Linux and Windows
8. Add the binaries to the GitHub release
