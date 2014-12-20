# Admin Interface

	<<#-->>
	package admin

	import (
		"bytes"
		"crypto/rand"
		"fmt"
		. "net/http"
		"os"
		"regexp"

		"github.com/tokenshift/blob/log"

		"github.com/boltdb/bolt"
		"github.com/gorilla/mux"
		"golang.org/x/crypto/sha3"
	)

The Admin Interface is a REST service intended to be exposed on a non-public
port that can be used by administrators of a **Node** cluster to configure
various settings and access internal data. It used HTTP Basic Auth to
authenticate the client and has a single admin user/password, configured using
environment variables.

	type AdminInterface interface {
		Serve(port int)
	}

	type server struct {
		mux.Router
		db *bolt.DB
		username string
		salt, hash []byte
	}

	func CreateAdminInterface(dbFile, username string, salt, hash []byte) (AdminInterface, error) {
		fi, err := os.Stat(dbFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if err == nil && fi.IsDir() {
			return nil, fmt.Errorf("%s is a directory", dbFile)
		}

		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			return nil, err
		}

		db.Update(func(tx *bolt.Tx) error {
			tx.CreateBucketIfNotExists([]byte(hashBucket))
			tx.CreateBucketIfNotExists([]byte(saltBucket))
			return nil
		})

		s := server{
			db: db,
			username: username,
			salt: salt,
			hash: hash,
		}

		s.HandleFunc("/apps", s.authorized(getAppsHandler)).Methods("GET")
		s.HandleFunc("/apps/{name}", s.authorized(putAppHandler)).Methods("PUT")

		return &s, nil
	}

Client app information is stored in a Bolt DB.

	const hashBucket = "Hashes"
	const saltBucket = "Salts"

Call `Serve` to begin serving the admin interface on the specified port.

	func (s *server) Serve(port int) {
		log.Info("Starting admin interface on port", port)
		ListenAndServe(fmt.Sprintf(":%d", port), s)
	}

Password hashing and verification is provided as a public function so that
consumers can create hashes matching the internal implementation. This is also
made available through the github.com/tokenshift/blob/admin/hash executable.

	func Hash(input, salt []byte) []byte {
		hash := sha3.New256()
		hash.Write(input)
		hash.Write(salt)
		return hash.Sum(nil)
	}

	const saltLen = 32

	func Salt() []byte {
		salt := make([]byte, saltLen)
		_, err := rand.Read(salt)
		if err != nil {
			panic(err)
		}
		return salt
	}

	func Verify(input, salt, hash []byte) bool {
		provided := Hash(input, salt)
		return bytes.Compare(provided, hash) == 0
	}

All endpoints require HTTP Basic authentication using the configured admin
username and password.

	type adminHandlerFunc func(*server, ResponseWriter, *Request)

	func (s *server) authorized(f adminHandlerFunc) HandlerFunc {
		return func(res ResponseWriter, req *Request) {
			username, password, ok := req.BasicAuth()
			if !ok {
				res.Header()["WWW-Authenticate"] = []string{"Basic"}
				res.WriteHeader(401)
				return
			}

			if username != s.username || !Verify([]byte(password), s.salt, s.hash) {
				res.Header()["WWW-Authenticate"] = []string{"Basic"}
				res.WriteHeader(401)
				res.Write([]byte("Invalid username or password.\n"))
				return
			}

			f(s, res, req)
		}
	}

Client applications can be set up with their own usernames and passwords at the
/apps endpoint.

	var rxValidAppName = regexp.MustCompile(`(?i)^[0-9a-z\-_]+$`)
	func validAppName(name string) bool {
		return rxValidAppName.MatchString(name)
	}

	func getAppsHandler(s *server, res ResponseWriter, req *Request) {
		s.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(saltBucket))
			c := b.Cursor()

			res.WriteHeader(200)
			// Only the app names (the keys) are returned.
			for app, _ := c.First(); app != nil; app, _ = c.Next() {
				res.Write(app)
				res.Write([]byte("\n"))
			}

			return nil
		})
	}

	func putAppHandler(s *server, res ResponseWriter, req *Request) {
		pass := req.URL.Query()["password"]
		if pass == nil {
			res.WriteHeader(400)
			res.Write([]byte("Password is required.\n"))
			return
		}
		if len(pass) > 1 {
			res.WriteHeader(400)
			res.Write([]byte("Password cannot be specified multiple times.\n"))
			return
		}

		appName := mux.Vars(req)["name"]
		if !validAppName(appName) {
			res.WriteHeader(400)
			res.Write([]byte("Invalid app name.\n"))
			return
		}

		salt := Salt()
		passHash := Hash([]byte(pass[0]), salt)

		s.db.Update(func(tx *bolt.Tx) error {
			salts := tx.Bucket([]byte(saltBucket))
			err := salts.Put([]byte(appName), salt)
			if err != nil {
				return err
			}

			hashes := tx.Bucket([]byte(hashBucket))
			return hashes.Put([]byte(appName), passHash)
		})
	}
