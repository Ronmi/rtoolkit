# ratelimit

ratelimit implements [Token-Bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm to help you limit transfer rate. It is rewrote from and inspired by [juju/ratelimit](https://github.com/juju/ratelimit).

## Synopsis

```go
resp, err := http.Get("http://example.com/")
if err != nil {
	// error process
}
defer resp.Body.Close()

bucket := ratelimit.NewFromRate(
	100*1024, // limit transfer rate to 100kb/s
	4*1024, // allocates 4k tokens each time
	48*1024, // at most 48k tokens in the bocket
)
r := ratelimit.NewReader(resp.Body, bucket)
io.Copy(dst, r)
```

See example folder for more.

## Total transfer rate

Bucket is thread-safe. You can share same Bucket betweens readers/writers to limit the total transfer rate.

## License

LGPL v3 or later.
