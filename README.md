xyr [WIP]
=========
> `xyr` is a very lightweight, simple and powerful data ETL platform that helps you to query available data sources using `SQL`.

Supported Drivers
=================
- [x] `jsondir`: for extracting, transforming and loading json documents (objects/[]objects) from local filesystem directory (recursive) into `xyr`.
- [x] `csvdir`: for extracting, transforming and loading csv documents from local filesystem into `xyr`.
- [ ] `s3+json`: for extracting, transforming and loading json documents from `S3`.
- [ ] `s3+csv`: for extracting, transforming and loading csv documents from `S3`.
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
> internaly `xyr` utilizes `SQLite` as an embeded sql datastore (it may be changed in future and we can add multiple data stores), when you define a table in `XYRCONFIG` file then run `$ xyr import` you will be able to import all defined tables as well querying them via `$ xyr exec "SELECT * FROM TABLE_NAME_HERE"`.

Plan
====
- [x] Building the initial core.
- [x] Add the basic `import` command for importing the tables into `xyr`.
- [ ] Add the `exec` command to execute SQL query.
- [ ] Expose another API beside the `CLI` to enable external Apps to query `xyr`.
    - [ ] JSON Endpoint?
    - [ ] Mysql Protocol?
    - [ ] Redis Protocol?
- [ ] Improving the code base (iteration 1).
