# 9. Web Helpers

Encoding/decoding are repetitive tasks that happen in almost every handler.

Move these behaviors to a dedicated package.

# Tasks:

- Add `internal/platform/web` package
  - Define `Encode` and `Decode` functions to centralize request/response behavior.
