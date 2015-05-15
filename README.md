# Blob

Distributed file/media storage service.

## Build

Blob uses GPM + GVP for reproducible builds. Installation details here:

https://github.com/pote/gpm
https://github.com/pote/gvp

To set up your $GOPATH and fetch all of the dependencies:

	```shell
	source gvp
	gpm install
	```

The Blob source is also written using [mdweb](https://github.com/tokenshift/mdweb),
a Markdown-based literate programming tool. It is included as a dependency in
`Godeps`, but will still need to be installed manually:

	```shell
	go install github.com/tokenshift/mdweb/mdtangle
	go install github.com/tokenshift/mdweb/mdweave
	```

Then run `mdtangle` against all of the `.go.md` source files to generate the
`.go` source (`mdweave` will generate the associated documentation), which can
then be built/installed normally.

	```shell
	mdtangle *.go.md
	mdweave *.go.md

	go build
	go install
	```

## Use

TODO

### Options

TODO
