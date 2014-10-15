# Manifest

Tracks uploaded files and related metadata, such as size and MIME type.

	<<#-->>
	package main

	import (
		"database/sql"
		"fmt"
		"io"
		"os"
		"path/filepath"
		"time"

		"code.google.com/p/go-uuid/uuid"
		_ "github.com/mattn/go-sqlite3"
	)

The manifest maps file paths to file contents on disk. So that multple versions
of a file can be stored, the actual path on disk is a randomly generated UUID;
the manifest will always return the newest version of the file. File metadata
(the `Entry` structure) includes the unique ID of the file, the MIME type, the
file size, the time the file was added/created, and the filename/path.

	type Entry struct {
		ID string
		MimeType string
		Added time.Time
		Path string
		Size int64
	}

The Manifest itself handles all interaction with file metadata, including
supporting simultaneous access from many goroutines. A single global manifest
is initialized at program start. SQLite3 is used as the data store for the
manifest.

	type Manifest struct {
		conn *sql.DB
	}

	var manifest = initManifest("data.sqlite3")

	func initManifest(fname string) Manifest {
		conn, err := sql.Open("sqlite3", fname)
		if err != nil {
			panic(err)
		}

		if err = conn.Ping(); err != nil {
			panic(err)
		}

		prepareDB(conn)

		return Manifest {
			conn: conn,
		}
	}

Required tables are created when the application starts.

	func prepareDB(conn *sql.DB) {
		_, err := conn.Exec("CREATE TABLE IF NOT EXISTS Entries (" +
			"ID CHARACTER(36) PRIMARY KEY," +
			"MimeType NVARCHAR(255)," +
			"Added DATETIME NOT NULL," +
			"Path TEXT," +
			"Size BIGINT)")
		if err != nil {
			panic(err)
		}
	}

The manifest handles simple Get, Put and Delete actions for retrieving,
creating/updating and removing files, respectively. Deletions are represented
in the database by an entry with file size 0.

	func (m Manifest) Get(path string) (Entry, error) {
		row := m.conn.QueryRow(
			"SELECT ID, MimeType, Added, Size " +
			"FROM Entries WHERE Path = ? " +
			"ORDER BY Added DESC LIMIT 1",
			path)

		var id, mimeType string
		var added time.Time
		var size int64

		err := row.Scan(&id, &mimeType, &added, &size)
		if err != nil {
			return Entry{}, NotFound{}
		}

		if size == 0 {
			return Entry{}, NotFound{}
		}

		return Entry {
			ID: id,
			MimeType: mimeType,
			Added: added,
			Path: path,
			Size: size,
		}, nil
	}

	func (m Manifest) Put(path string, mimeType string, data io.Reader) error {
		id := uuid.New()

		f, err := os.Create(filepath.Join("data", id))
		if err != nil {
			return err
		}
		defer f.Close()

		size, err := io.Copy(f, data)
		if err != nil {
			return err
		}

		if size == 0 {
			os.Remove(filepath.Join("data", id))
			return BadRequest(fmt.Sprintf("Did not receive any data for file %s", path))
		}

		_, err = m.conn.Exec(
			"INSERT INTO Entries " +
			"(ID, MimeType, Added, Path, Size) " +
			"VALUES (?, ?, ?, ?, ?)",
			id, mimeType, time.Now(), path, size)
		if err != nil {
			return err
		}

		return nil
	}

	func (m Manifest) Delete(path string) error {
		_, err := m.conn.Exec(
			"INSERT INTO Entries " +
			"(ID, MimeType, Added, Path, Size) " +
			"VALUES (?, ?, ?, ?, ?)",
			uuid.New(), nil, time.Now(), path, 0)

		return err
	}

Entries provide access to the file data that they correspond to.

	func (e Entry) Open() (io.Reader, error) {
		f, err := os.Open(filepath.Join("data", e.ID))
		return f, err
	}
