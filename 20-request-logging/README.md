# 20. Request Logging

- Add a middleware to log something for each request.
- Requires creating a struct with some request values and passing it down through context.
- Ensure web.Respond updates the value.
- Extract error handling to a middleware so it can finish before the request logger.
