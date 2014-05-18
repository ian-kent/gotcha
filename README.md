Gotcha
======

A(nother) web framework for Go.

### Getting started

* Make sure your [Go environment](http://golang.org/doc/install) is configured.
* Download the latest [Gotcha release](https://github.com/ian-kent/gotcha/releases)
* Create a new application: ```gotcha new MyApp```
* Change to the application directory: ```cd MyApp```
* Install your application: ```gotcha install```
* Run your application: ```MyApp```
* Open your application in a browser: [http://localhost:7050](http://localhost:7050)

### Features

* Easy deployment
  * Produces a single binary file with no dependencies
  * Easy to cross-compile using [gox](https://github.com/mitchellh/gox)
* Simple request routing
  * Supports regexes with named capture groups
  * GET/POST/PUT/PATCH/DELETE/OPTIONS
  * Static content handler
* Simple action composition
* Per-connection data stash
* Streaming/chunked responses
* Cached template rendering with [html/template](http://golang.org/pkg/html/template)
* Embedded assets using [jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata)
* Not tied to any design pattern (e.g. MVC, MVP etc)

### Why another web framework

Deploying web applications is difficult. Almost every language has
some kind of dependency which needs installation.

Go makes it easy. And Gotcha makes it better:

* A platform portable web framework
* Easy to use, even easier to deploy
* Installation-free, suitable for cloud environments

### Principles of Gotcha

* Care more about syntax than implementation!
* If it's too big to fit in memory, it's **TOO BIG**
* If you want to read assets from disk, use a CDN
* If you have highly dynamic assets, use a cache/database
* Configure from your environment, not from file
* Re-use core Go libraries wherever possible
* Try to avoid using external dependencies

### Contributing

#### Feature requests

* [Open a new issue](https://github.com/ian-kent/gotcha/issues/new)
* Explain your use-case(s)
* Explain why it should be part of the framework
* Have a go at implementing it :)

#### Pull requests

* Clone this repository: ```git clone https://github.com/ian-kent/gotcha```
* Run tests: ```make test```
* Install gotcha: ```make```
* Stick to the principles!

Before submitting a pull request:

  * Run ```go fmt ./...```
  * Make sure tests pass: ```make test```

### Roadmap

Without putting too much thought into it:

* Add nicer routing syntax which translates to a regex
* Write some tests

### Licence

Copyright ©‎ 2014, Ian Kent (http://www.iankent.eu).

Released under MIT license, see [LICENSE](LICENSE.md) for details.
