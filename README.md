xyr [WIP]
=========
> `xyr` is a very lightweight, simple and powerful data ETL platform that helps you to query available data sources using `SQL`.

Supported Drivers
=================
- [x] `local+json`: for extracting, transforming and loading json documents from local filesystem into `xyr`.
- [x] `local+csv`: for extracting, transforming and loading csv documents from local filesystem into `xyr`.
- [ ] `s3+json`: for extracting, transforming and loading json documents from `S3`.
- [ ] `s3+csv`: for extracting, transforming and loading json documents from `S3`.
- [ ] `postgresql`: for extracting, transforming and loading postgres results.
- [ ] `clickhouse`: for extracting, transforming and loading clickhouse results.

Use Cases
=========
- Simple Presto Alternative.
- Simple AWS Athena Alternative.
- Convert your json documents into a SQL db.
- Query your CSV files easily and join them with other data.

How it works?
=============
> internaly `xyr` utilizes `SQLite` as an embeded sql datastore (it may be changed in future and we can add multiple data stores), when you define a table in `XYRCONFIG` file then run `$ xyr import` you will be able to import all defined tables as well querying them via `$ xyr exec "SELECT * FROM TABLE_NAME_HERE"`.