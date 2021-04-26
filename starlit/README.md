# starlit (cli)

To download a new runtime, specify its version and the target folder in the `go:generate` in main.go.

## Development
To run the serve command with the latest static assets and templates without having to `go generate`, use:

```bash
cd starlit/web/app
npm run start

# In another terminal
go build &&  ./starlit serve . --static_folder web/app/dist/static --templates_folder web/templates
```

Consider it live-reload for everything that is not defined in Go files :).

## Releases

Releases are minted on CI, you can create one locally by running
```
goreleaser --snapshot --skip-publish --rm-dist
```
