Introduction to Gotcha
======================

## Creating your first application

Creating your first application is easy. Once you've installed
Gotcha, just run ```gotcha new AppName``` from the command line.

You'll get a directory structure similar to this:

    AppName/
      - assets/
        - css/
          - default.css
        - images/
          - logo-ish.png
        - templates/
          - error.html
          - index.html
          - notfound.html
      - Makefile
      - README.md
      - main.go

To build your application, run ```make```.

This will build a binary application named ```AppName```, with
the assets in the ```assets``` directory embedded.

You can run your application using ```AppName```, which will start
on port ```7050``` by default.

### Command line parameters

| Parameter | Description
| --------- | -----------
| help      | Display application usage information
| listen    | The interface and port to bind on, e.g. ```0.0.0.0:7050```

For example:

    AppName -listen=:1234

### Embedded assets

Any content inside the ```asset``` directory is automatically embedded
in your application when you build it with ```make```.

Assets are embedded using ```github.com/jteeuwen/go-bindata```, and
can be retrieved at run-time using the built-in asset loader:

    b, err := Asset("assets/css/default.css")

This returns either a byte array, or an error if the file isn't found.

### Application

The application object provides your interface to the web application:

	import(
		gotcha "github.com/ian-kent/gotcha/app"
	)
	var app = gotcha.Create(Asset)

### Routing

Routing is configured using the router:

	r := app.Router

A route is a pattern and an associated handler:

	r.Get("/", welcome)

The pattern is treated as a regex, anchored at the start and end:

    ^/$

#### Placeholders

You can use placeholders in your route to capture path information:

	r.Get("/users/(?P<user>.*)", user_info)

Placeholder values are placed in the stash, which is available through
the ```http.Session``` object passed to the handler function.

    val := session.Stash["user"]

#### Handlers

The handler is any function matching ```func(*http.Session)```.

    func example(session *http.Session) {
		session.Stash["Title"] = "Welcome to Gotcha"
		session.Render("index.html")
	}

#### Data modeling and form validation

You can create a data model to represent your form data:

    type MyForm struct {
    	Name string
    }

You can assign validation rules to fields in your model:

    type MyForm struct {
    	Name string `minlength:1; maxlength:200`
    }

The form helper can populate and validate this for you:

    model := &MyForm{}
	session.Stash["form"] = form.New(session, model).Populate().Validate()
	session.Render("index.html")

You can then access form properties from templates:

    {{ .form.Model.Name }}
    {{ .form.HasErrors }}
    {{ .form.Errors["Name"] }}
    {{ .form.Errors["Name"]["minlength"] }}
    {{ .form.Values["Name"] }}

#### Serving static content

You can setup routes which only serve static content without a handler:

	r.Get("/file.txt)", r.Static("assets/file.txt"))

```r.Static``` implements a default handler which uses the asset loader
to return the named file.

You can also use a route placeholder to serve multiple files or a directory:

	r.Get("/images/(?P<file>.*)", r.Static("assets/images/{{file}}"))

### Session stash

A stash is available throughout the life of the session.

You can use it to pass information between handlers and lifecycle events, 
or to pass data to a template for rendering.

    session.Stash["foo"] = "bar"

Named placeholders in a route are automatically placed in the stash.

### Templates

Templates are implemented using ```html/template```.

To render a template, call the session Render function:

    session.Render("welcome.html")

The stash is automatically passed to the template during rendering,
so variables in the stash are available to the template:

    <h1>{{ .foo }}</h1>

## Examples

Examples of Gotcha applications can be found in the assets directory:

| Example                     | Description
| --------------------------- | -----------
| [demo_app](assets/demo_app) | A demo application which uses most of Gotcha's features
| [new_app](assets/new_app)   | The template for a ```gotcha new``` application
| [tiny_app](assets/tiny_app) | A pointless attempt at golfing a gotcha application

### Licence

Copyright ©‎ 2014, Ian Kent (http://www.iankent.eu).

Released under MIT license, see [LICENSE](LICENSE.md) for details.
