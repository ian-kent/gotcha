Gotcha
======

A(nother) web framework for Go.

### Getting started

* Make sure your [Go environment](http://golang.org/doc/install) is configured.
* Download the latest Gotcha release
* Create a new application: ```gotcha new MyApp```
* Change to the application directory: ```cd MyApp```
* Install your application: ```gotcha install```
* Run your application: ```MyApp```
* Open your application in a browser: http://localhost:7050

### Why another web framework

This is probably best explained as what I'm trying to create:

* A platform portable web framework
* Easy to use, even easier to deploy
* Installation-free, suitable for cloud environments
* A friendly community willing to support each other

Every web framework I've used has failed on at least one of those.

This is my attempt at fixing it!

### Why Go?

It's portable. That makes the first point easy to achieve with
almost no additional work.

It's also got (mostly) nice syntax and awesome concurrency support!

### Principles

* Care more about syntax than implementation!
* If it's too big to fit in memory, it's **TOO BIG**
* If you want to read assets from disk, use a CDN
* If you have highly dynamic assets, use a cache/database

### Contributing

* Clone this repository: ```git clone https://github.com/ian-kent/gotcha```
* Run tests: ```make test```
* Install gotcha: ```make```

If you make any changes, run ```go fmt ./...``` before submitting a pull request.

### Licence

Copyright ©‎ 2014, Ian Kent (http://www.iankent.eu).

Released under MIT license, see [LICENSE](license) for details.
