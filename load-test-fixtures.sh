# Create and populate the users table
rethinkdb import -f fixtures/users.json --table casgo.users --pkey email
