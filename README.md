Gotcha
======

A(nother) web framework for Go.

### Getting started

Set your GOPATH and add it to PATH:
```
    export GOPATH=/path/to/src/root
    export PATH=$PATH:$GOPATH/bin
```

Install Gotcha (unless you've downloaded a binary release):
```
    cd /path/to/src/root/github.com/ian-kent/Gotcha
    make
```

Somewhere else inside your GOPATH, create an application:
```
    cd /path/to/src/root/github.com/your-name
    gotcha new MyApp
    cd MyApp
```

Install and run your application:
```
    gotcha install
    MyApp
```

Open your application in a browser: http://localhost:7050

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
