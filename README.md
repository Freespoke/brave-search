# Brave Search API for Go

A Go client for the [Brave Search API](https://api.search.brave.com/).

## Install

```sh
$ go get dev.freespoke.com/brave-search
```

## Usage

```go
client, err := New(key)
if err != nil {
    log.Fatal(err)
}

res, err := client.WebSearch(context.Background(), "freespoke")
if err != nil {
    log.Fatal(err)
}
```

## Documentation

* [Official Documentation](https://api.search.brave.com/app/documentation/get-started) (account required)

## Status

The client works, but is poorly tested. Some features and types are not available.
The documentation is unclear, and in certain places it seems that the API was
rushed out to take advantage of the opportunity to capture businesses fleeing
from Bing. That's not to say it's a bad API, just that there's a bit of weirdness
that I try to paper over, and that the documentation seems to miss certain things,
and that certain types don't seem to be returned as I'd expect.

I'll try to keep this up-to-date and adopt new stuff as it becomes available.

## License

MIT (See LICENSE.md).
