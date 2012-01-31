Latenza API
===========

* `GET key [field]*`
  * Retrieves the fields specified within record `key`. If no field is specified, gets all attributes of record.
* `SET key field value [field value]*`
  * Sets field=value for all specified field/value pairs for the record specified by key
* `DEL key [field]*`
  * Delete all specified fields from the record specified by `key`. If no fields are specified, deletes the entire record.
* `INCR key field delta [field delta]*`
  * Increment the value of `field` in the record `key` by the amount `delta`.
* `ADD key field value [field value]*`
  * Insert value(s) into set `field` within record `key`.
* `SUB key field value [value]*`
  - Delete the values in `field` within record `key`.
* `CARD key field`
  - Returns the cardinality of the set specifed by `field` in the record specified by `key`.
* `MEMBER key field value`
  - Returns whether or not value is a member of set `field` within record `key`.

