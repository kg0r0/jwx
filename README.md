# github.com/lestrrat-go/jwx ![](https://github.com/lestrrat-go/jwx/workflows/CI/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/lestrrat-go/jwx.svg)](https://pkg.go.dev/github.com/lestrrat-go/jwx) [![codecov.io](http://codecov.io/github/lestrrat-go/jwx/coverage.svg?branch=main)](http://codecov.io/github/lestrrat-go/jwx?branch=main)

Command line tool [jwx](./cmd/jwx) and libraries implementing various JWx technologies. Please click on the package names in the table below to find the synopsis/description for each package.

| Package name                                              | Notes                                           |
|-----------------------------------------------------------|-------------------------------------------------|
| [jwt](https://github.com/lestrrat-go/jwx/tree/main/jwt) | [RFC 7519](https://tools.ietf.org/html/rfc7519) |
| [jwk](https://github.com/lestrrat-go/jwx/tree/main/jwk) | [RFC 7517](https://tools.ietf.org/html/rfc7517) + [RFC 7638](https://tools.ietf.org/html/rfc7638) |
| [jwa](https://github.com/lestrrat-go/jwx/tree/main/jwa) | [RFC 7518](https://tools.ietf.org/html/rfc7518) |
| [jws](https://github.com/lestrrat-go/jwx/tree/main/jws) | [RFC 7515](https://tools.ietf.org/html/rfc7515) + [RFC 7797](https://tools.ietf.org/html/rfc7797) |
| [jwe](https://github.com/lestrrat-go/jwx/tree/main/jwe) | [RFC 7516](https://tools.ietf.org/html/rfc7516) |

# Index

* [Documentation on pkg.go.dev](https://pkg.go.dev/github.com/lestrrat-go/jwx)
  * HTML version of what you can see using `go doc` command
* [How-to style documentation](./docs)
  * Frequently asked questions.
  * How to JWx That? Documentation by example.
* Overview of this package
  * Read on for more gory details.

# Description

## History

My goal was to write a server that heavily uses JWK and JWT. At first glance
the libraries that already exist seemed sufficient, but soon I realized that

1. To completely implement the protocols, I needed the entire JWT, JWK, JWS, JWE (and JWA, by necessity).
2. Most of the libraries that existed only deal with a subset of the various JWx specifications that were necessary to implement their specific needs

For example, a certain library looks like it had most of JWS, JWE, JWK covered, but then it lacked the ability to include private claims in its JWT responses. Another library had support of all the private claims, but completely lacked in its flexibility to generate various different response formats.

Because I was writing the server side (and the client side for testing), I needed the *entire* JOSE toolset to properly implement my server, **and** they needed to be *flexible* enough to fulfill the entire spec that I was writing.

So here's `github.com/lestrrat-go/jwx`. This library is extensible, customizable, and hopefully well organized to the point that it is easy for you to slice and dice it.

## Why would I use this library?

There are several other major Go modules that handle JWT and related data formats,
so why should you use this library?

From a purely functional perspective, the only major difference is this:
Whereas most other projects only deal with what they seem necessary to handle
JWTs, this module handles the entire spectrum of JWS, JWE, JWK, and JWT.

That is, if you need to not only parse JWTs, but also to control JWKs, or
if you need to handle payloads that are NOT JWTs, you should probably consider
using this module.

Next, from an implementation perspective, this module differs significantly
from others in that it tries very hard to expose only the APIs, and not the
internal data. For example, individual JWT claims are not accessible through
struct field lookups. You need to use one of the getter methods.

This is because this library takes the stance that the end user is fully capable
and even willing to shoot themselves on the foot when presented with a lax
API. By making sure that users do not have access to open structs, we can protect
users from doing silly things like creating _incomplete_ structs, or access the
structs concurrently without any protection. This structure also allows
us to put extra smarts in the structs, such as doing the right thing when
you want to parse / write custom fields (this module does not require the user
to specify alternate structs to parse objects with custom fields)

In the end I think it comes down to your usage pattern, and priorities.
Some general guidelines that come to mind are:

* If you want a single library to handle everything JWx, such as handling [auto-refreshing JWKs](https://github.com/lestrrat-go/jwx/blob/main/docs/04-jwk.md#auto-refreshing-remote-keys), use this module.
* If you want to honor all possible custom fields transparently, use this module.
* If you want a standardized clean API, use this module.

Otherwise, feel free to choose something else.

# How to Use

* [API documentation](https://pkg.go.dev/github.com/lestrrat-go/jwx)
* [How-to style documentation](./docs)
* [Runnable Examples](./examples)
* Test files.

(Depending on what you want to do, you will need navigate between multiple packages within this package)

# Global Settings

## Allowing single element in 'aud' field

When you marshal `"github.com/lestrrat-go/jwx/jwt".Token` into JSON, by default the `aud` field is serialized as an array of strings. This field may take either a single string or array form, but apparently there are parsers that do not understand the array form.

The examples below shoud both be valid, but apparently there are systems that do not understand the former ([AWS Cognito has been reported to be one such system](https://github.com/lestrrat-go/jwx/issues/368)).

```
{
  "aud": ["foo"],
  ...
}
```

```
{
  "aud": "foo",
  ...
}
```

To workaround these problematic parsers, you may use the `jwt.Settings()` function with the `jwt.WithFlattenAudience(true)` option.

```go
func init() {
  jwt.Settings(jwt.WithFlattenAudience(true))
}
```

The above call will force all calls to marshal JWT tokens to flatten the `aud` field when it can. This has global effect.

## Enabling ES256K

Some algorithms are intentionally left out because they are not as common in the wild, and you may want to avoid compiling this extra information in.
To enable these, you must explicitly provide a build tag.

| Algorithm        | Build Tag  |
|:-----------------|:-----------|
| secp256k1/ES256K | jwx_es256k |

If you do not provide these tags, the program will still compile, but it will return an error during runtime saying that these algorithms are not supported.

## Switching to a faster JSON library

By default we use the standard library's `encoding/json` for all of our JSON needs.
However, if performance for parsing/serializing JSON is really important to you, you might want to enable [github.com/goccy/go-json](https://github.com/goccy/go-json) by enabling the `jwx_goccy` tag.

```shell
% go build -tags jwx_goccy ...
```

[github.com/goccy/go-json](https://github.com/goccy/go-json) is *disabled* by default because it uses some really advanced black magic, and I really do not feel like debugging it **IF** it breaks. Please note that that's a big "if".
As of github.com/goccy/go-json@v0.3.3 I haven't see any problems, and I would say that it is mostly stable.

However, it is a dependency that you can go without, and I won't be of much help if it breaks -- therefore it is not the default.
If you know what you are doing, I highly recommend enabling this module -- all you need to do is to enable this tag.
Disable the tag if you feel like it's not worth the hassle.

And when you *do* enable [github.com/goccy/go-json](https://github.com/goccy/go-json) and you encounter some mysterious error, I also trust that you know to file an issue to [github.com/goccy/go-json](https://github.com/goccy/go-json) and **NOT** to this library.

## Using json.Number

If you want to parse numbers in the incoming JSON objects as json.Number
instead of floats, you can use the following call to globally affect the behavior of JSON parsing.

```go
func init() {
  jwx.DecoderSettings(jwx.WithUseNumber(true))
}
```

Do be aware that this has *global* effect. All code that calls in to `encoding/json`
within `jwx` *will* use your settings.

## Decode private fields to objects

Packages within `github.com/lestrrat-go/jwx` parses known fields into pre-defined types,
but for everything else (usually called private fields/headers/claims) are decoded into
wharever `"encoding/json".Unmarshal` deems appropriate.

For example, JSON objects are converted to `map[string]interface{}`, JSON arrays into
`[]interface{}`, and so on.

Sometimes you know beforehand that it makes sense for certain fields to be decoded into
proper objects instead of generic maps or arrays. When you encounter this, you can use
the `RegisterCustomField()` method in each of `jwe`, `jwk`, `jws`, and `jwt` packages.

```go
func init() {
  jwt.RegisterCustomField(`x-foo-bar`, mypkg.FooBar{})
}
```

This tells the decoder that when it encounters a JWT token with the field named
`"x-foo-bar"`, it should be decoded to an instance of `mypkg.FooBar`. Then you can
access this value by using `Get()`

```go
v, _ := token.Get(`x-foo-bar`)
foobar := v.(mypkg.FooBar)
```

Do be aware that this has *global* effect. In the above example, all JWT tokens containing
the `"x-foo-bar"` key will decode in the same way. If you need this behavior from
`jwe`, `jwk`, or `jws` packages, you need to do the same thing for each package.

# Command Line Tool

Since v1.1.1 we have a command line tool `jwx` (*). With `jwx` you can create JWKs (from PEM files, even), sign and verify JWS message, encrypt and decrypt JWE messages, etc.

(*) Okay, it existed since a long time ago, but it was never useful.

## Installation

```
go install github.com/lestrrat-go/jwx/cmd/jwx
```

# Caveats

## Backwards Compatibility Notice

### Users of github.com/lestrrat/go-jwx

Uh, why are you using such an ancient version? You know that repository is archived for a reason, yeah? Please use the new version.

### Pre-1.0.0 users

The API has been reworked quite substantially between pre- and post 1.0.0 releases. Please check out the [Changes](./Changes) file (or the [diff](https://github.com/lestrrat-go/jwx/compare/v0.9.2...v1.0.0), if you are into that sort of thing)

### v1.0.x users

The API has gone under some changes for v1.1.0. If you are upgrading, you might want to read the relevant parts in the [Changes](./Changes) file.

# Contributions

## Issues

For bug reports and feature requests, please try to follow the issue templates as much as possible.
For either bug reports or feature requests, failing tests are even better.

## Pull Requests

Please make sure to include tests that excercise the changes you made.

If you are editing auto-generated files (those files with the `_gen.go` prefix, please make sure that you do the following:

1. Edit the generator, not the generated files (e.g. internal/cmd/genreadfile/main.go)
2. Run `make generate` (or `go generate`) to generate the new code
3. Commit _both_ the generator _and_ the generated files

## Discussions / Usage

Please try [discussions](https://github.com/lestrrat-go/jwx/discussions) first.

# Credits

* Work on this library was generously sponsored by HDE Inc (https://www.hde.co.jp)
* Lots of code, especially JWE was taken from go-jose library (https://github.com/square/go-jose)
* Lots of individual contributors have helped this project over the years. Thank each and everyone of you very much.

