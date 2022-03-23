# tableman-golang

Table based text generation language interpreter.

## TODO

- Web server + api
  - (DONE) guid-based user context
  - (Done) session-based root context and roll history. Don't think shared ctx is needed, but history is done.
  - (DONE) Load packs from file.
  - (DONE) `/pack` for a list of pack names, load a pack to the session.
  - `/history` get roll history.
  - `/tables` get tables in the pack.
  - (DONE) `/eval` evaluate expression in a pack.
- Basic Web UI
- Execution stack limit?
- roll history max horizon.
- User Docs
- Go Docs/lint
  - compilation and runtime explanation.
- stack trace errors
  - runtime
  - compile
  - parse
- Add float support?
