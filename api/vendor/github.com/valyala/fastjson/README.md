[![Build Status](https://travis-ci.org/valyala/fastjson.svg)](https://travis-ci.org/valyala/fastjson)
[![GoDoc](https://godoc.org/github.com/valyala/fastjson?status.svg)](http://godoc.org/github.com/valyala/fastjson)
[![Go Report](https://goreportcard.com/badge/github.com/valyala/fastjson)](https://goreportcard.com/report/github.com/valyala/fastjson)
[![codecov](https://codecov.io/gh/valyala/fastjson/branch/master/graph/badge.svg)](https://codecov.io/gh/valyala/fastjson)

# fastjson - fast JSON parser and validator for Go


## Features

  * Fast. As usual, up to 15x faster than the standard [encoding/json](https://golang.org/pkg/encoding/json/).
    See [benchmarks](#benchmarks).
  * Parses arbitrary JSON without schema, reflection, struct magic and code generation
    contrary to [easyjson](https://github.com/mailru/easyjson).
  * Provides simple [API](http://godoc.org/github.com/valyala/fastjson).
  * Outperforms [jsonparser](https://github.com/buger/jsonparser) and [gjson](https://github.com/tidwall/gjson)
    when accessing multiple unrelated fields, since `fastjson` parses the input JSON only once.
  * Validates the parsed JSON unlike [jsonparser](https://github.com/buger/jsonparser)
    and [gjson](https://github.com/tidwall/gjson).
  * May quickly extract a part of the original JSON with `Value.Get(...).MarshalTo`.
  * May parse array containing values with distinct types (aka non-homogenous types).
    For instance, `fastjson` easily parses the following JSON array `[123, "foo", [456], {"k": "v"}, null]`.
  * `fastjson` preserves the original order of object items when calling
    [Object.Visit](https://godoc.org/github.com/valyala/fastjson#Object.Visit).


## Known limitations

  * Requies extra care to work with - references to certain objects recursively
    returned by [Parser](https://godoc.org/github.com/valyala/fastjson#Parser)
    must be released before the next call to [Parse](https://godoc.org/github.com/valyala/fastjson#Parser.Parse).
    Otherwise the program may work improperly.
    Adhere recommendations from [docs](https://godoc.org/github.com/valyala/fastjson).
  * Cannot parse JSON from `io.Reader`. There is [Scanner](https://godoc.org/github.com/valyala/fastjson#Scanner)
    for parsing stream of JSON values from a string.


## Usage

One-liner accessing a single field:
```go
	s := []byte(`{"foo": [123, "bar"]}`)
	fmt.Printf("foo.0=%d\n", fastjson.GetInt(s, "foo", "0"))

	// Output:
	// foo.0=123
```

Accessing multiple fields with error handling:
```go
        var p fastjson.Parser
        v, err := p.Parse(`{
                "str": "bar",
                "int": 123,
                "float": 1.23,
                "bool": true,
                "arr": [1, "foo", {}]
        }`)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("foo=%s\n", v.GetStringBytes("str"))
        fmt.Printf("int=%d\n", v.GetInt("int"))
        fmt.Printf("float=%f\n", v.GetFloat64("float"))
        fmt.Printf("bool=%v\n", v.GetBool("bool"))
        fmt.Printf("arr.1=%s\n", v.GetStringBytes("arr", "1"))

        // Output:
        // foo=bar
        // int=123
        // float=1.230000
        // bool=true
        // arr.1=foo
```

See also [examples](https://godoc.org/github.com/valyala/fastjson#pkg-examples).


## Security

  * `fastjson` shouldn't crash or panic when parsing input strings specially crafted
    by an attacker. It must return error on invalid input JSON.
  * `fastjson` requires up to `sizeof(Value) * len(inputJSON)` bytes of memory
    for parsing `inputJSON` string. Limit the maximum size of the `inputJSON`
    before parsing it in order to limit the maximum memory usage.


## Performance optimization tips

  * Re-use [Parser](https://godoc.org/github.com/valyala/fastjson#Parser) and [Scanner](https://godoc.org/github.com/valyala/fastjson#Scanner)
    for parsing many JSONs. This reduces memory allocations overhead.
    [ParserPool](https://godoc.org/github.com/valyala/fastjson#ParserPool) may be useful in this case.
  * Prefer calling `Value.Get*` on the value returned from [Parser](https://godoc.org/github.com/valyala/fastjson#Parser)
    instead of calling `Get*` one-liners when multiple fields
    must be obtained from JSON, since each `Get*` one-liner re-parses
    the input JSON again.
  * Prefer calling once [Value.Get](https://godoc.org/github.com/valyala/fastjson#Value.Get)
    for common prefix paths and then calling `Value.Get*` on the returned value
    for distinct suffix paths.
  * Prefer iterating over array returned from [Value.GetArray](https://godoc.org/github.com/valyala/fastjson#Object.Visit)
    with a range loop instead of calling `Value.Get*` for each array item.


## Benchmarks

Go 1.11 has been used for benchmarking.

Legend:

  * `small` - parse [small.json](testdata/small.json) (190 bytes).
  * `medium` - parse [medium.json](testdata/medium.json) (2.3KB).
  * `large` - parse [large.json](testdata/large.json) (28KB).
  * `canada` - parse [canada.json](testdata/canada.json) (2.2MB).
  * `citm` - parse [citm_catalog.json](testdata/citm_catalog.json) (1.7MB).
  * `twitter` - parse [twitter.json](testdata/twitter.json) (617KB).

  * `stdjson-map` - parse into a `map[string]interface{}` using `encoding/json`.
  * `stdjson-struct` - parse into a struct containing
    a subset of fields of the parsed JSON, using `encoding/json`.
  * `stdjson-empty-struct` - parse into an empty struct using `encoding/json`.
    This is the fastest possible solution for `encoding/json`, may be used
    for json validation. See also benchmark results for json validation.
  * `fastjson` - parse using `fastjson` without fields access.
  * `fastjson-get` - parse using `fastjson` with fields access similar to `stdjson-struct`.

```
$ GOMAXPROCS=1 go test github.com/valyala/fastjson -bench='Parse$'
goos: linux
goarch: amd64
pkg: github.com/valyala/fastjson
BenchmarkParse/small/stdjson-map         	  200000	      7249 ns/op	  26.21 MB/s	     960 B/op	      51 allocs/op
BenchmarkParse/small/stdjson-struct      	  500000	      3731 ns/op	  50.92 MB/s	     224 B/op	       4 allocs/op
BenchmarkParse/small/stdjson-empty-struct         	  500000	      2566 ns/op	  74.04 MB/s	     168 B/op	       2 allocs/op
BenchmarkParse/small/fastjson                     	 3000000	       411 ns/op	 461.69 MB/s	       0 B/op	       0 allocs/op
BenchmarkParse/small/fastjson-get                 	 2000000	       674 ns/op	 281.73 MB/s	       0 B/op	       0 allocs/op
BenchmarkParse/medium/stdjson-map                 	   30000	     41120 ns/op	  56.64 MB/s	   10195 B/op	     208 allocs/op
BenchmarkParse/medium/stdjson-struct              	   30000	     47793 ns/op	  48.73 MB/s	    9174 B/op	     258 allocs/op
BenchmarkParse/medium/stdjson-empty-struct        	  100000	     20962 ns/op	 111.10 MB/s	     280 B/op	       5 allocs/op
BenchmarkParse/medium/fastjson                    	  500000	      3499 ns/op	 665.50 MB/s	       0 B/op	       0 allocs/op
BenchmarkParse/medium/fastjson-get                	  300000	      3679 ns/op	 632.94 MB/s	       0 B/op	       0 allocs/op
BenchmarkParse/large/stdjson-map                  	    2000	    613724 ns/op	  45.82 MB/s	  210761 B/op	    2785 allocs/op
BenchmarkParse/large/stdjson-struct               	    5000	    284621 ns/op	  98.79 MB/s	   15616 B/op	     353 allocs/op
BenchmarkParse/large/stdjson-empty-struct         	    5000	    250047 ns/op	 112.45 MB/s	     280 B/op	       5 allocs/op
BenchmarkParse/large/fastjson                     	   30000	     48737 ns/op	 576.93 MB/s	       9 B/op	       0 allocs/op
BenchmarkParse/large/fastjson-get                 	   30000	     49540 ns/op	 567.58 MB/s	       9 B/op	       0 allocs/op
BenchmarkParse/canada/stdjson-map                 	      20	  67939960 ns/op	  33.13 MB/s	12260502 B/op	  392539 allocs/op
BenchmarkParse/canada/stdjson-struct              	      20	  67322448 ns/op	  33.44 MB/s	12260123 B/op	  392534 allocs/op
BenchmarkParse/canada/stdjson-empty-struct        	     100	  17227143 ns/op	 130.67 MB/s	     280 B/op	       5 allocs/op
BenchmarkParse/canada/fastjson                    	     200	   6082885 ns/op	 370.06 MB/s	  472007 B/op	     571 allocs/op
BenchmarkParse/canada/fastjson-get                	     200	   6038088 ns/op	 372.81 MB/s	  472007 B/op	     571 allocs/op
BenchmarkParse/citm/stdjson-map                   	      50	  27171438 ns/op	  63.57 MB/s	 5213923 B/op	   95401 allocs/op
BenchmarkParse/citm/stdjson-struct                	     100	  14270811 ns/op	 121.03 MB/s	    1989 B/op	      75 allocs/op
BenchmarkParse/citm/stdjson-empty-struct          	     100	  14335648 ns/op	 120.48 MB/s	     280 B/op	       5 allocs/op
BenchmarkParse/citm/fastjson                      	    1000	   2026047 ns/op	 852.50 MB/s	   17633 B/op	      30 allocs/op
BenchmarkParse/citm/fastjson-get                  	    1000	   2025209 ns/op	 852.85 MB/s	   17633 B/op	      30 allocs/op
BenchmarkParse/twitter/stdjson-map                	     100	  11230390 ns/op	  56.23 MB/s	 2187443 B/op	   31264 allocs/op
BenchmarkParse/twitter/stdjson-struct             	     300	   5411600 ns/op	 116.70 MB/s	     408 B/op	       6 allocs/op
BenchmarkParse/twitter/stdjson-empty-struct       	     300	   5347742 ns/op	 118.09 MB/s	     408 B/op	       6 allocs/op
BenchmarkParse/twitter/fastjson                   	    2000	    834384 ns/op	 756.86 MB/s	    2536 B/op	       2 allocs/op
BenchmarkParse/twitter/fastjson-get               	    2000	    827723 ns/op	 762.95 MB/s	    2536 B/op	       2 allocs/op
```

Benchmark results for json validation:

```
$ GOMAXPROCS=1 go test github.com/valyala/fastjson -bench='Validate$'
goos: linux
goarch: amd64
pkg: github.com/valyala/fastjson
BenchmarkValidate/small/stdjson 	 2000000	       846 ns/op	 224.51 MB/s	      72 B/op	       2 allocs/op
BenchmarkValidate/small/fastjson         	 5000000	       352 ns/op	 539.57 MB/s	       0 B/op	       0 allocs/op
BenchmarkValidate/medium/stdjson         	  200000	      9783 ns/op	 238.05 MB/s	     184 B/op	       5 allocs/op
BenchmarkValidate/medium/fastjson        	  500000	      3297 ns/op	 706.34 MB/s	       0 B/op	       0 allocs/op
BenchmarkValidate/large/stdjson          	   10000	    119264 ns/op	 235.76 MB/s	     184 B/op	       5 allocs/op
BenchmarkValidate/large/fastjson         	   30000	     48079 ns/op	 584.82 MB/s	       0 B/op	       0 allocs/op
BenchmarkValidate/canada/stdjson         	     200	   8229071 ns/op	 273.55 MB/s	     184 B/op	       5 allocs/op
BenchmarkValidate/canada/fastjson        	     500	   3662160 ns/op	 614.68 MB/s	       0 B/op	       0 allocs/op
BenchmarkValidate/citm/stdjson           	     200	   6760134 ns/op	 255.50 MB/s	     184 B/op	       5 allocs/op
BenchmarkValidate/citm/fastjson          	    1000	   1581792 ns/op	1091.93 MB/s	       0 B/op	       0 allocs/op
BenchmarkValidate/twitter/stdjson        	     500	   2596556 ns/op	 243.21 MB/s	     312 B/op	       6 allocs/op
BenchmarkValidate/twitter/fastjson       	    2000	    889861 ns/op	 709.68 MB/s	       0 B/op	       0 allocs/op
```

## FAQ

  * Q: _There are a ton of other high-perf packages for JSON parsing in Go. Why creating yet another package?_
    A: Because other packages require either rigid JSON schema via struct magic
       and code generation or perform poorly when multiple unrelated fields
       must be obtained from the parsed JSON.
       Additionally, `fastjson` provides nicer [API](http://godoc.org/github.com/valyala/fastjson).

  * Q: _What is the main purpose for `fastjson`?_
    A: High-perf JSON parsing for [RTB](https://www.iab.com/wp-content/uploads/2015/05/OpenRTB_API_Specification_Version_2_3_1.pdf)
       and other [JSON-RPC](https://en.wikipedia.org/wiki/JSON-RPC) services.

  * Q: _Why fastjson doesn't provide fast marshaling (serialization)?_
    A: Because other solutions exist. I'd recommend [quicktemplate](https://github.com/valyala/quicktemplate#use-cases)
       for high-performance JSON marshaling :)

  * Q: _`fastjson` crashes my program!_
    A: There is high probability of improper use.
       * Make sure you don't hold references to objects recursively returned by `Parser` / `Scanner`
         beyond the next `Parser.Parse` / `Scanner.Next` call
         if such restriction is mentioned in [docs](https://github.com/valyala/fastjson/issues/new).
       * Make sure you don't access `fastjson` objects from concurrently running goroutines
         if such restriction is mentioned in [docs](https://github.com/valyala/fastjson/issues/new).
       * Build and run your program with [-race](https://golang.org/doc/articles/race_detector.html) flag.
         Make sure the race detector detects zero races.
       * If your program continue crashing after fixing issues mentioned above, [file a bug](https://github.com/valyala/fastjson/issues/new).
