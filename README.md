Latenzad
========

A key-value store.

Installation
------------

First, make sure you have [Go](http://golang.org) installed and `$GOROOT` set properly. You will need Go1 or a weekly release to build; Go r60.3 or earlier is too old.

    $ git clone git://github.com/codeslinger/latenza.git
    $ cd latenza
    $ ./do

API
---

* `CREATE table name`
  - Create `table` using the field `name` as the table's primary key.
* `DROP table`
  - Delete `table` and all of its contents.
* `GET table key [field]*`
  * Retrieves the fields specified within record `key`. If no field is specified, gets all attributes of record.
* `MGET table key [key]*`
  * Retrieves multiple full records in one call.
* `SET table key field value [field value]*`
  * Sets field=value for all specified field/value pairs for the record specified by key
* `DEL table key [field]*`
  * Delete all specified fields from the record specified by `key`. If no fields are specified, deletes the entire record.
* `INCR table key field delta [field delta]*`
  * Increment the value of `field` in the record `key` by the amount `delta`.
* `ADD table key field value [field value]*`
  * Insert value(s) into set `field` within record `key`.
* `SUB table key field value [value]*`
  - Delete the values in `field` within record `key`.
* `CARD table key field`
  - Returns the cardinality of the set specifed by `field` in the record specified by `key`.
* `MEMBER table key field value`
  - Returns whether or not value is a member of set `field` within record `key`.
* `STATS`
  - Returns status report for running instance.

Data Types
----------

* A given field in a record can be one of the following:
  * String (UTF-8 encoded)
  * Integral (64-bit signed integer)
  * Set of strings
  * Set of integers
* Data types are specified at the time of field creation.
* Attempts to call `INCR` on a string value will result in an error response.

Constraints
-----------

* Table, key and field names are interpreted as UTF-8 byte strings. They must be 1-256 bytes in length.
* Total size of a given record must be less than or equal to 1MB over all field values.

