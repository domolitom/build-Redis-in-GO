# build-Redis-in-GO

Errors encountered during the build process:

- Make sure the redis server is running before running the test. (`conn, err := l.Accept()` does not go further)

## Things to fix: !!!

right now deadlock if we 'hget' a key that does not exist.
