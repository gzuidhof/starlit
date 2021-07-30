# starlit (cli)


## nbtest
Starlit NBTest runs notebooks end-to-end, if no errors are thrown the notebook passes. It spawns a local webserver that serves the notebook tests, and will talk to a Chrome Debug Protocol-compatible browser (Chrome, Edge, Safari, but not Firefox) to run the tests in parallel.


### Usage examples
```bash
# Print usage (there are some options like concurrency and test timeout)
startlit nbtest --help

# Test notebooks under a path
starlit nbtest some/path

# Test a single notebook
starlit nbtest some/path/my-notebook.nb

# Serve mode, in this mode the webserver will stay alive. When a test fails a clickable link will show for easy debugging.
starlit nbtest --serve some/path

# Non-headless mode, you will be able to see the actual browser running the tests.
starlit nbtest --headless=false some/path

# Specify where pyodide artifacts should be loaded from (defaults to JSDelivr CDN)
starlit nbtest --pyodide_artifacts folder/with/artifacts some/path
starlit nbtest --pyodide_artifacts https://cdn.jsdelivr.net/pyodide/v0.17.0/full/ some/path

# Specify where starboard artifacts should be loaded from (defaults to the starboard files baked into the starlit binary)
starlit nbtest --starboard_artifacts some/path/to/a/folder
starlit nbtest --starboard_artifacts http://localhost:9001 some/path  # reminder: cors required!
starlit nbtest --starboard_artifacts https://cdn.jsdelivr.net/starboard-notebook/v0.12.3/dist some/path
```

### Skipping tests
Notebooks that a skip flag in their metadata are skipped.
```yaml
nbtest:
  skip: true
```

## Development
To download a new runtime, update `build_static.sh`.

To run the serve command with the latest static assets and templates without having to `go generate`, use:

```bash
cd starlit/web/app
npm run start

# In another terminal
go build &&  ./starlit serve . --static_folder web/app/dist/static --templates_folder web/templates
```

Consider it live-reload for everything that is not defined in Go files :).


```
go build && ./starlit.exe nbtest ../../starboard-notebook/src/debugNotebooks --templates_folder web/templates
```

### Windows development
When using the embedded filesystem templates don't get found on Windows platforms, unless they were built on a UNIX-based platform (i.e. downloaded from the Github releases page).
For local development pass the `--templates_folder web/templates` flag to work around this issue.

## Releases

Releases are minted on CI, you can create one locally by running
```
goreleaser --snapshot --skip-publish --rm-dist
```
