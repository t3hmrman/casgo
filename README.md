# CASGO CAS Authentication Server

## What is CASGO?

Casgo is a simple to use, simple to deploy [Single Sign On](http://en.wikipedia.org/wiki/Single_sign-on) that uses the [CAS protocol](http://en.wikipedia.org/wiki/Central_Authentication_Service) developed by Shawn Bayern of Yale University.

## CAS Spec

Casgo implements version 1.0 of the [CAS Specification](http://www.yale.edu/tp/cas/specification/CAS%202.0%20Protocol%20Specification%20v1.0.html) as defined with a few key changes:

- JSON is preferred over XML/plaintext responses
- The /validate endpoint behaves as specified in CAS 1.0 (success/failure and the username of the user)
- The /validate endpoint returns user attributes

## Getting started (deploying an instance of Casgo)

0. Install your database of choice (default is [RethinkDB](http://rethinkdb.com), version 2.0+)
1. Download the casgo binary for your operating system
2. Ensure port 443 is open (and your database instance is at the right port, 28015 by default)
3. Run the binary

## Contributing to Casgo (setting up your local development environment)

0. Install your database of choice (default is [RethinkDB](http://rethinkdb.com), version 2.0+)
1. Install [bower](http://bower.io)
2. `bower install`
3. Install [go.rice](https://github.com/GeertJohan/go.rice)
4. `go get github.com/t3hmrman/casgo`
5. `make all` (or `go install`/`go build`)
6. Ensure port 443 is open (and your database instance is at the right port, 28015 by default)
7. `casgo`

NOTE - You may have to add an exception for the included self-signed certificate

You may also find it helpful to install [Gin](https://github.com/codegangsta/gin)

## Running tests

0. (If using RethinkDB) Download the RethinkDB python driver
1. Run `load-test-fixtures.sh` to load a running RethinkDB instance with fixture data.
2. Install [ginkgo](https://github.com/onsi/ginkgo) and [agouti](https://github.com/sclevine/agouti)
3. `ginkgo -r` (from the main casgo directory)

*Note* As some tests rely on the database to be up, RethinkDB must be running.

`make test`

OR

`cd cas && ginkgo -r`

## Options

|Option       |Description                                            |
|-------------|-------------------------------------------------------|
|-config      | Specify a (JSON) configuration file for CasGo to use. |

## Configuration

### By File

Casgo can be configured by file if you specify the `-c/--config <filename>` flag. See **Options** section for a full list of CASGO's command line options.

### By ENV

|Variable (json)          |ENV                  |default                 |description                                        |
|-------------------------|---------------------|------------------------|---------------------------------------------------|
|**host**                 |CASGO_HOST           |"0.0.0.0"               |The host on which to run casgo                     |
|**port**                 |CASGO_PORT           |"8080"                  |The port on which to run casgo                     |
|**dbHost**               |CASGO_DBHOST         |"localhost:28015"       |The hostname of database instance                  |
|**dbName**               |CASGO_DBNAME         |"casgo"                 |The database name for casgo to use                 |
|**templatesDirectory**   |CASGO_TEMPLATES      |"templates/"            |The folder in which casgo templates reside         |
|**companyName**          |CASGO_COMPNAME       |"companyABC"            |The database name for casgo to use                 |
|**authMethod**           |CASGO_DEFAULT_AUTH   |"password"              |The default (user) authentication method for casgo |
|**logLevel**             |CASGO_LOG_LVL        |"WARN|DEBUG|INFO"       |The default log level for casgo                    |
|**tlsCertFile**          |CASGO_TLS_CERT       |"fixtures/ssl/cert.pem" |The TLS cert file that casgo will use              |
|**tlsKeyFile**           |CASGO_TLS_KEY        |"fixtures/ssl/eckey.pem"|The TLS key file that casgo will use               |


### Database Schema

So what does the database that powers casgo look like?

|Database |Table    |Description                                                   |
|---------|---------|--------------------------------------------------------------|
|casgo    |tickets  |The authentication tickets currently in use by the casgo      |
|casgo    |services |Services authorized to use casgo                              |
|casgo    |users    |User data stored by casgo (if not using external auth)        |
|casgo    |api_keys |Authentication API keys (enabling non-web app authentication) |


### Contributing

0. Fork the repo
1. Install [Go](http://golang.org)
2. Install your database of choice (default is  [RethinkDB](http://rethinkdb.com))
3. Fix issues, make changes
4. Ensure Makefile functions correctly (`make all`/`make resources`/etc)
5. Pull Request
6. Receive thanks from the community
