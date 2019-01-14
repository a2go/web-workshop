# 9. Refactor Handlers

HTTP errors and Encoding/decoding are a very repetitive tasks that happens in
almost every handler. Move these behaviors to a dedicated package and create a
middleware to execute for every handler.

# Tasks:

- Add `internal/platform/web` package
  - Add a custom error type that knows about HTTP status codes.
  - Define a `Decode` function to centralize request decoding behavior.
  - Create a custom func type for a web handler that returns dynamic values and an error.
  - Add a function `Encode` that runs as middleware accepting the new handler type. Have it control error responses as well as normal response encoding.
- Update handler methods to match `web.HandlerFunc`
