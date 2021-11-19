xyr
====
> `xyr` is a very lightweight, simple, and powerful data ETL platform that helps you to query available data sources using `SQL`.

Example (Local Filesystem)
===========================
> here we define a new table called `users` which will load all json files in that directory (recursive) with any of the following json formats: (object/object[] per-file, newline delimited json objects/object[], or event no delimiter json objects/object[] like what kinesis firehose json output format).

> Let's image we have a directory of json files called `/tmp/data/users` and here is an example of a json file there:
```json
{"id":10,"email":"u10@example.com"}{"id":20,"email":"u20@example.com"}{"id": 3,"email":"u3@example.com"}{"id": 4,"email":"u4@example.com"}
```

> Then we can define its schema as following
```hcl
# where xyr should store its internal database
data_dir = "./tmp/db/"

# this file is `./config.xyr.hcl`
table "users" {
    // the driver we want
    driver = "jsondir"

    // the data source directory
    source = "/tmp/data/users"

    // xyr will try to create a table into its internal storage, so it needs
    // to know at least what are the required columns names of your data.
    // i.e: {"id": 1, "email": "user@example.com", "age": 20}
    // but we only need "id" and "email", so we defined both in the below columns array
    // and not that the ordering is the same as our example.
    columns = ["id", "email"]

    // what do you want to load
    // in case of jsondir, we can specify a regex pattern to filter the files 
    // using the filename
    // but if we're using an SQL driver we can provide an sql statement that reads the data
    // from the source SQL based database.
    // i.e: "SELECT * FROM SOME_TABLE"
    filter = ".*"
}
```

> Now its the time to load it
```bash
$ xyr table:import users
```

> Now let's query it
```bash
$ xyr exec "SELECT * FROM users"
```

> All tables you define could be joined in the same query easily, let's imagine that we have the following defination
```hcl
# where xyr should store its internal database
data_dir = "./tmp/db/"

table "users" {
    driver = "s3jsondir"
    source = "s3://ACCESS_KEY:SECRET_KEY@/BUCKET_NAME?region=&ssl=false&path=true&perpage=1000"

    # which prefix we want to select
    filter = "xyr/users/"

    columns = ["id", "email"]
}

table "user_vists" {
    driver = "postgres"
    source = "postgresql://username:password@server:port/dbname?option1=value1"
    columns = ["user_id", "vists"]
    filter = "SELECT user_id, count(vists) FROM USERS GROUP BY user_id"
}
```

> Now let's join them
```bash
$ xyr exec "SELECT * FROM users LEFT JOIN user_vists ON users_vists.user_id = users.id"
```

Installation
============
> use this [docker package](https://github.com/alash3al/xyr/pkgs/container/xyr)

Supported Drivers
=================
| Driver | Source Connection String |
---------| ------------------------ |
| `jsondir`     | `/PATH/TO/JSON/DATA/DIR`|
| `s3jsondir`   | `s3://[access_key_url_encoded]:[secret_key_url_encoded]@[endpoint_url]/bucket_name?region=&ssl=false&path=true&perpage=1000`|
| `mysql`       | `usrname:password@tcp(server:port)/dbname?option1=value1&...`|
| `postgres`    | `postgresql://username:password@server:port/dbname?option1=value1`|
| `sqlite3`     | `/path/to/db.sqlite?option1=value1`|
| `sqlserver`   | `sqlserver://username:password@host/instance?param1=value&param2=value` |
|               | `sqlserver://username:password@host:port?param1=value&param2=value`|
|               | `sqlserver://sa@localhost/SQLExpress?database=master&connection+timeout=30`|
| `hana`        | `hdb://user:password@host:port` |
| `clickhouse`  | `tcp://host1:9000?username=user&password=qwerty&database=clicks&read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000` |
| `oracle`      | `oracle://user:pass@server1/service?server=server2&server=server3` |

Use Cases
=========
- Simple Presto Alternative.
- Simple AWS Athena Alternative.
- Convert your JSON documents into a SQL DB.
- Query your CSV files easily and join them with other data.

How does it work?
==================
> internaly `xyr` utilizes `SQLite` as an embeded sql datastore (it may be changed in future and we can add multiple data stores), when you define a table in `XYRCONFIG` file then run `$ xyr table:import` you will be able to import all defined tables as well querying them via `$ xyr exec "SELECT * FROM TABLE_NAME_HERE"` which outputs json result by default.

Plan
====
- [x] Building the initial core.
- [x] Add the basic `import` command for importing the tables into `xyr`.
- [x] Add the `exec` command to execute SQL query.
- [x] Add well known SQL drivers
    - [x] mysql
    - [x] postgres
    - [x] sqlite3
    - [x] clickhouse
    - [x] oracle
    - [x] hana
    - [x] sqlserver
- [x] Add an S3 driver
- [ ] Adding/Improving documentations
- [ ] Expose another API beside the `CLI` to enable external Apps to query `xyr`.
    - [ ] JSON Endpoint?
    - [ ] Mysql Protocol?
    - [ ] Redis Protocol?
- [ ] Improving the code base (iteration 1).
- [ ] Add another backend instead of sqlite3 as internal datastore?
