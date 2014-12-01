## CASGO CAS Authentication Server

### What is CASGO?

Casgo is a simple to use, simple to deploy [Single Sign On](http://en.wikipedia.org/wiki/Single_sign-on) that uses the [CAS protocol](http://en.wikipedia.org/wiki/Central_Authentication_Service) developed by Shawn Bayern of Yale University.

### CAS Spec

Casgo implements the CAS Specification as defined with a few key changes:

- JSON is preferred over XML/plaintext responses

### Getting started

0. Install your database of choice (default is  [RethinkDB](http://rethinkdb.com))
1. Download the casgo binary for your operating system
2. Ensure port 9090 is open (and your database instance is at the right port, 28015 by default)
3. Run the binary


### Casgo configuration

|Variable                 |ENV              |default           |description                                |
|-------------------------|-----------------|------------------|-------------------------------------------|
|**Host**                 |HOST             |"0.0.0.0"         |The host on which to run casgo             |
|**Port**                 |PORT             |"8080"            |The port on which to run casgo             |
|**DBHost**               |DBHOST           |"localhost:28015" |The hostname of database instance          |
|**DBName**               |DBNAME           |"casgo"           |The database name for casgo to use         |
|**TemplatesDirectory**   |CASGO_TEMPLATES  |"templates/"      |The folder in which casgo templates reside |
|**CompanyName**          |CASGO_COMPNAME   |"companyABC"      |The database name for casgo to use         |


### Database Schema

So what does the database that powers casgo look like?

|Database |Table    |Description                                                   |
|---------|---------|--------------------------------------------------------------|
|casgo    |tickets  |The authentication tickets currently in use by the casgo      |
|casgo    |services |Services authorized to use casgo                              |
|casgo    |users    |User data stored by casgo (if not using external auth)        |


### Contributing 

0. Fork the repo
1. Install [Go](http://golang.org)
2. Install your database of choice (default is  [RethinkDB](http://rethinkdb.com))
3. Fix issues, make changes
4. Pull Request
5. Receive thanks from the community

