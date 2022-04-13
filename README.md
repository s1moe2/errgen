# Sample - error codes

This sample repo serves two purposes. To demonstrate (in a very simplified manner) a setup to easily reuse error codes
for API responses AND a generator  that looks for declarations of type `AppError`
and generates markdown content in the description field of an OpenAPI spec file.

For visualization purposes, open `openapi.yaml` and remove the lines between the comments:
```
<!-- ERROR_GENERATOR_START -->
<!-- ERROR_GENERATOR_END -->
```

Run `go generate`.

Profit.