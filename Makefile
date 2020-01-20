create_dist:
	mkdir -p ./dist/

clean: create_dist
	rm -fr ./dist/*

build: create_dist
	gox -os="darwin linux windows" -arch="amd64 386" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

build_darwin: create_dist
	gox -os="darwin" -arch="amd64 386" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

build_linux: create_dist
	gox -os="linux" -arch="amd64 386" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

build_windows: create_dist
	gox -os="windows" -arch="amd64 386" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
