set -e

# Download the npm packages we want to bundle
go run scripts/download_runtime/main.go starboard-notebook 0.12.0 web/app/static/vendor
go run scripts/download_runtime/main.go starboard-wrap 0.4.0 web/app/static/vendor

# Build web assets
pushd ./web/app 
  npm ci && npm run build
popd

# Basically clear everything that is not web/static/fs.go
rm -rf web/static/fonts 
rm -rf web/static/js 
rm -rf web/static/images 
rm -rf web/static/styles
rm -rf web/static/vendor

cp -r web/app/dist/static web/

# Generate a go file which helps us resolve the npm packages to the bundled version.
go run scripts/index_vendored_packages/main.go web/static/vendor web/vendored_libraries.go
