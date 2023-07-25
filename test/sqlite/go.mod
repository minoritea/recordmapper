module github.com/minoritea/recordmapper/test/sqlite

go 1.20

require (
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/minoritea/recordmapper v0.0.0-20230724005140-102f5b78a852
)

require github.com/iancoleman/strcase v0.3.0 // indirect

replace github.com/minoritea/recordmapper => ./../..
