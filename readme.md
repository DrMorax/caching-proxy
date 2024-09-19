# Caching Proxy

A proxy for caching requests with their status code, headers, and body.

## Usage

Important: You need Golang installed to run this app.

1. Pick a location
   `cd your/desired/directory`
2. Copy the source code
   `git clone https://github.com/DrMorax/caching-proxy`
3. Run the local server
   `go run ./cmd/web`
4. In you browser, go to
   `https://localhost:4000?url=<target url>`

once the response is successfully loaded, it will be cached in-memory and retrieved from it in the second attempt to access the same url with the same method
