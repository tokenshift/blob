# Errors

Error types used throughout the program.

	<<#-->>

	package main

400 - BAD REQUEST

	type BadRequest string

	func (err BadRequest) Error() string {
		return string(err)
	}

404 - NOT FOUND

	type NotFound struct{}

	func (err NotFound) Error() string {
		return "Not Found"
	}
