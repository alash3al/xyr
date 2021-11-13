xyr [WIP]
=========
> `xyr` is a very lightweight, simple and powerful data ETL platform that helps you to query available data sources using `SQL`.

Example
=======
> here we define a new table called `users` which will load all json files in that directory (recursive) with any of the following json formats: (object/object[] per-file, newline delimited json objects/object[], or event no delimiter json objects/object[] like what kinesis firehose json output format).

```hcl
# this file is `./config.xyr.hcl`
table "users" {
    // [D]ata [S]ource [N]ame as: "[drivername]://driver_specs"
    dsn = "jsondir://./tmp/data/users"

    // which columns we want to load
    // this expect that our json objects contains these properties.
    columns = ["id", "email"]

    // which filenames we want to be only loaded into xyr using regular expressions
    loader = ".*"
}
```

```bash
$ xyr table:import users
```

Installation
============
> use this [docker package](https://github.com/alash3al/xyr/pkgs/container/xyr)

Supported Drivers
=================
- [x] `jsondir`: for extracting, transforming and loading json documents (objects/[]objects) from local filesystem directory (recursive) into `xyr`.
- [ ] `jsons3`: for extracting, transforming and loading json documents from `S3` (something like AWS Athena).
- [ ] `postgresql`: for extracting, transforming and loading postgres results.
- [ ] `clickhouse`: for extracting, transforming and loading clickhouse results.
- [ ] `redis`: for extracting, transforming and loading redis datastructures.

Use Cases
=========
- Simple Presto Alternative.
- Simple AWS Athena Alternative.
- Convert your json documents into a SQL db.
- Query your CSV files easily and join them with other data.

How it works?
=============
> internaly `xyr` utilizes `SQLite` as an embeded sql datastore (it may be changed in future and we can add multiple data stores), when you define a table in `XYRCONFIG` file then run `$ xyr table:import` you will be able to import all defined tables as well querying them via `$ xyr exec "SELECT * FROM TABLE_NAME_HERE"` which outputs json result by default.

Plan
====
- [x] Building the initial core.
- [x] Add the basic `import` command for importing the tables into `xyr`.
- [x] Add the `exec` command to execute SQL query.
- [ ] Expose another API beside the `CLI` to enable external Apps to query `xyr`.
    - [ ] JSON Endpoint?
    - [ ] Mysql Protocol?
    - [ ] Redis Protocol?
- [ ] Improving the code base (iteration 1).
