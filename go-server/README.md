# mta-proxy server

Simple HTTP web server written in Go to accomplish the following:

1. Proxy requests for the status directly to the [MTA's feed](http://web.mta.info/status/serviceStatus.txt) with CORS header attached
2. Provide some basic pre-parsed options for students who are not yet ready to parse themselves (with CORS headers as well).

## Roadmap

- [x] Parse to HTML table
- [ ] Add more options for testing without SSL
- [ ] Robustify error handling
- [ ] Parse to JSON (`Accept: application/json`)
