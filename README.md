## CASGO CAS Authentication Server

### What is CASGO?

Casgo is a simple to use, simple to deploy [Single Sign On](http://en.wikipedia.org/wiki/Single_sign-on) that uses the [CAS protocol](http://en.wikipedia.org/wiki/Central_Authentication_Service) developed by Shawn Bayern of Yale University.


### Casgo configuration


|Variable               |ENV              |default           |description                                |
|-----------------------|-----------------|------------------|-------------------------------------------|
|*Host*                 |HOST             |"0.0.0.0"         |The host on which to run casgo             |
|*Port*                 |PORT             |"8080"            |The port on which to run casgo             |
|*DBHost*               |DBHOST           |"localhost:28015" |The hostname of database instance          |
|*DBName*               |DBNAME           |"casgo"           |The database name for casgo to use         |
|*TemplatesDirectory*   |CASGO_TEMPLATES  |"templates/"      |The folder in which casgo templates reside |
|*CompanyName*          |CASGO_COMPNAME   |"companyABC"      |The database name for casgo to use         |
