## CasGo Database Schema

A breakdown of Casgo's RethinkDB (No-SQL) schema. This schema only applies to servers using the default RethinkDB based database adapter.
If creating your own database adapter, please ensure to provide a similar document so others can build on your work. 

### Ticket

Tickets that will be used by CAS to validate logins

**Primary Key** - id (generated)

|field      |type    |description                                      |
|-----------|--------|-------------------------------------------------|
|serviceId  |string  |ID of the service this ticket belongs to         |
|userEmail  |string  |Email (id) of the user that was authenticated    |

#### Example
    {
       "serviceId": "this-is-not-a-real-id",
       "userEmail": "test@test.com"
    }

### Service

Registered services (applications) that may authenticate through the CasGO instance

**Primary Key** - id (generated)

|field      |type    |description                                      |
|-----------|--------|-------------------------------------------------|
|name       |string  |Name of the service (displayable)                |
|url        |string  |Redirect URL used upon successful user auth      |
|adminEmail |string  |Administrator contact email                      |

#### Example
    {
       "name": "test_service",
       "url": "localhost:9090/validateCASLogin",
       "adminEmail": "admin@test.com"
    }

### User

Users who are registered to use the application

**Primary Key** - email

Passwords are hashed with [bcrypt](http://en.wikipedia.org/wiki/Bcrypt)

|field      |type    |description                                      |
|-----------|--------|-------------------------------------------------|
|email      |string  |Email address of the user                        |
|password   |string  |Password of the user                             |
|isAdmin    |boolean |Whether user is admin                            |
|services   |list    |List of user's services eventually-consistent    |

#### Example
    {
       "email": "test@test.com",
       "password": "NczbWiVimnqUegfmoQYOqjCLNYXjFGJooHwbUezKXyYqFXHzCZAgZwMRAsmXKFfM",
       "isAdmin": false,
       "services": [
           {name: "casgo test service", url...},
           {name: "casgo second test service", url...},
           ...
           ]
    }
