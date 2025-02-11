### Create database on Postgres

Connect to postgres as a superuser

```bash
psql -U postgres
```

 To create the database, use the following SQL command:
 
```bash
CREATE DATABASE jwt_auth_db;
```

After executing the command, you should see output like:

```bash
CREATE DATABASE
```

To exit psql, type: `\q`


### how to verify the creation of the users table and its columns within the jwt_auth_db database

connect to your jwt_auth_db database:
```bash
psql -U postgres -d jwt_auth_db
```

list all tables in the current database
```bash
jwt_auth_db=# \dt
         List of relations
 Schema | Name  | Type  |  Owner   
--------+-------+-------+----------
 public | users | table | postgres
(1 row)
```

 Describe the users Table to Check Columns:
```bash
jwt_auth_db=# \d users
                                        Table "public.users"
    Column     |           Type           | Collation | Nullable |              Default              
---------------+--------------------------+-----------+----------+-----------------------------------
 id            | bigint                   |           | not null | nextval('users_id_seq'::regclass)
 created_at    | timestamp with time zone |           |          | 
 updated_at    | timestamp with time zone |           |          | 
 deleted_at    | timestamp with time zone |           |          | 
 username      | text                     |           | not null | 
 password_hash | text                     |           | not null | 
 email         | text                     |           |          | 
 full_name     | text                     |           |          | 
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "idx_users_deleted_at" btree (deleted_at, deleted_at)
    "idx_users_email" UNIQUE, btree (email)
    "idx_users_username" UNIQUE, btree (username)
```

To exit psql, type: `\q`

### Register a user

```bash
curl -X POST -H "Content-Type: application/json" -d '{
    "username": "testuser",
    "password": "password123",
    "email": "testuser@example.com",
    "fullName": "Test User"
}' http://localhost:8080/api/register
```

You can verify it through `psql`

```bash
jwt_auth_db=# SELECT * FROM users;
 id |          created_at           |          updated_at           | deleted_at | username |                        password_hash                         |        email         | full_name 
----+-------------------------------+-------------------------------+------------+----------+--------------------------------------------------------------+----------------------+-----------
  1 | 2025-02-11 11:09:04.206541+00 | 2025-02-11 11:09:04.206541+00 |            | testuser | $2a$10$zSMDOzdDQm6MZb6Kl1qaouyhVwVzHpTSTd4tSbYYsuTAYXxg5Xth2 | testuser@example.com | Test User
(1 row)
```

### Login as a user

```bash
curl -X POST -H "Content-Type: application/json" -d '{
        "usernameOrEmail": "testuser",
        "password": "password123"
    }' http://localhost:8080/api/login
```

### Access a protected endpoint

Send a GET request to `/api/protected`, including the JWT token you obtained from login in the `Authorization` header as a Bearer token.  Replace `YOUR_JWT_TOKEN` with the actual JWT string you copied.

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/api/protected
```

