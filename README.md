# Blob

Distributed file/media storage service. **Blob** provides a REST service
wrapping files stored on the local filesystem, along with a database of file
metadata, as an alternative to storing files (like images) in your database.

**In Development**

* HTTP Basic Authentication, files namespaced by client application.
* Sibling nodes for distributed file storage.
* Certificate-based authentication of sibling nodes.

## Build

Blob uses GPM + GVP for reproducible builds. Installation details here:

https://github.com/pote/gpm
https://github.com/pote/gvp

To set up your $GOPATH and fetch all of the dependencies:

```bash
$ source gvp
$ gpm install
```

The Blob source is also written using [mdweb](https://github.com/tokenshift/mdweb),
a Markdown-based literate programming tool. It is included as a dependency in
`Godeps`, but will still need to be installed manually:

```bash
$ go install github.com/tokenshift/mdweb/mdtangle
$ go install github.com/tokenshift/mdweb/mdweave
```

Then run `mdtangle` against all of the `.go.md` source files to generate the
`.go` source (`mdweave` will generate the associated documentation), which can
then be built/installed normally.

```bash
$ mdtangle *.go.md
$ mdweave *.go.md

$ go build
$ go install
```

## Use

All configuration is performed using environment variables (listed below).
After setting the required variables, simply run the `blob` executable to start
a single node.

### File Service

The file service responds to HTTP requests as follows, where the provided path
is used as the filename. No additional validation of the path is performed, so
client applications can implement their own folder structure/namespacing of
files as needed. (This is safe, because Blob stores the files at paths based on
an ID generated internally, with the filename stored as part of the metadata.)

* **GET**  
  Retrieves an existing file.
  * _200_  
    File exists, returned in body of response.
  * _404_  
    File not found.
* **HEAD**  
  Retrieves file metadata.
  * _200_  
  File exists; only the file metadata will be returned (in response headers).
  * _404_  
  File not found.
* **PUT**  
  Uploads a new file or updates an existing one.
  * _201_  
    Uploaded a new file.
  * _200_  
    Updated an existing file.
* **DELETE**  
  Deletes an existing file.
  * _200_  
    File was deleted.
  * _404_  
    File not found.

TODO: Content-Type, ETag, caching

### Admin Service

TODO

### Options

* **$BLOB_FILE_SERVICE_PORT**  
  The port that the main file service will run on.
* **$BLOB_ADMIN_SERVICE_PORT**  
  The port that the admin/config interface will run on.
* **$BLOB_ADMIN_SERVICE_USERNAME**  
  Username for accessing the admin service (using HTTP basic auth).
* **$BLOB_ADMIN_SERVICE_PASSWORD**  
  Password for accessing the admin service (using HTTP basic auth).
* **$BLOB_FILE_STORE_DB**  
  Filename for the file metadata db. Will be created if it does not already
  exist.
* **$BLOB_FILE_STORE_DIR**  
  Folder where uploaded files will be stored.
