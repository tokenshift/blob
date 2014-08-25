# Blob - Request Handler

The request handler listens on the configured port (80 by default) for file
requests and dispatches them to workers, after checking them against the
manifest and ensuring the requestor is authorized.
