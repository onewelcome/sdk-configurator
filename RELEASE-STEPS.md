The following steps must be executed to create a release

1. Update the release notes
2. Update the version constant in version/version.go
4. Commit & push to GitHub
5. Create a new tag and push
6. Update the release notes in the release
7. Build the configurator binaries for OS X, Linux and Windows (execute: `make release`)
8. Add the binaries to the GitHub release
