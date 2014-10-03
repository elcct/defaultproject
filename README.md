Default Project [![Build Status](https://drone.io/github.com/elcct/defaultproject/status.png)](https://drone.io/github.com/elcct/defaultproject/latest)
===============

Provides essentials that most web applications need - MVC pattern and user authorisation that can be easily extended.

It consists of 3 core components:

- Goji - A web microframework for Golang - http://goji.io/
- Gorilla web toolkit sessions - cookie (and filesystem) sessions - http://www.gorillatoolkit.org/pkg/sessions
- mgo - MongoDB driver for the Go language - http://labix.org/mgo

# Dependencies

Default Project requires `Go`, `MongoDB` and few other tools installed.

Instructions below have been tested on `Ubuntu 14.04`.

## Installation

If you don't have `Go` installed, follow installation instructions described here: http://golang.org/doc/install

Then install remaining dependecies:

```
sudo apt-get install git mercurial subversion bzr
```

MongoDB:

```
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv 7F0CEB10
sudo echo 'deb http://downloads-distro.mongodb.org/repo/debian-sysvinit dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list
sudo apt-get update
sudo apt-get install mongodb-org
```



No go to your GOPATH location and run:

```
go get github.com/elcct/defaultproject
```

And then:

```
go install github.com/elcct/defaultproject
```

In your GOPATH directory you can create `config.json` file:

```
{
	"secret": "secret",
	"public_path": "./src/github.com/elcct/defaultproject/public",
	"template_path": "./src/github.com/elcct/defaultproject/views",	
	"database": {
		"hosts": "localhost",
		"database": "defaultproject"
	}
}
```

Finally, you can run:

```
./bin/defaultproject
```

That should output something like:

```
2014/06/19 15:31:15.386961 Starting Goji on [::]:8000
```

And it means you can now direct your browser to `localhost:8000`

# Project structure

`/controllers`

All your controllers that serve defined routes.

`/helpers`

Helper functions.

`/models`

You database models.

`/public`

It has all your static files mapped to `/assets/*` path except `robots.txt` and `favicon.ico` that map to `/`.

`/system`

Core functions and structs.

`/views`

Your views using standard `Go` template system.

`server.go`

This file starts your web application and also contains routes definition.

# Make it your own

I assume you have followed installation instructions and you have `defaultproject` installed in your `GOPATH` location.

Let's say I want to create `Amazing Website`. I create new `GitHub` repository `https://github.com/elcct/amazingwebsite` (of course replace that with your own repository).

Now I have to prepare `defaultproject`. First thing is that I have to delete its `.git` directory.

I issue:

```
rm -rf src/github.com/elcct/defaultproject/.git
```

Then I want to replace all references from `github.com/elcct/defaultproject` to `github.com/elcct/amazingwebsite`:

```
grep -rl 'github.com/elcct/defaultproject' ./ | xargs sed -i 's/github.com\/elcct\/defaultproject/github.com\/elcct\/amazingwebsite/g'
```

Now I have to move all `defaultproject` files to the new location:

```
mv src/github.com/elcct/defaultproject/ src/github.com/elcct/amazingwebsite
```

And push it to my new repository at `GitHub`:

```
cd src/github.com/elcct/amazingwebsite
git init
git add --all .
git commit -m "Amazing Website First Commit"
git remote add origin https://github.com/elcct/amazingwebsite.git
git push -u origin master
```

You can now go back to your `GOPATH` and check if everything is ok:

```
go install github.com/elcct/amazingwebsite
```

And that's it. 

# Continuous Development

For Continuous Development I recommend using `Fresh` - https://github.com/pilu/fresh

You can install `Fresh` by issuing:

```
go get github.com/pilu/fresh
```

Then create a config file `runner.conf` in your `GOPATH`:

```
root:              ./src/github.com/elcct/amazingwebsite
tmp_path:          ./tmp
build_name:        runner-build
build_log:         runner-build-errors.log
valid_ext:         .go, .tpl, .tmpl, .html
build_delay:       600
colors:            1
log_color_main:    cyan
log_color_build:   yellow
log_color_runner:  green
log_color_watcher: magenta
log_color_app:
```

Note: Remember to replace `./src/github.com/elcct/amazingwebsite` with your own location

Now if you run:

```
./bin/fresh -c runner.conf
```

Project should automatically rebuild itself when a change occurs.
