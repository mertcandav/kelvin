# Kelvin Manual

Kelvin is a document-oriented database that uses the generic support of Go.


## Getting Started

To use Kelvin in your Go project, you must have the Kelvin module. \
Import Kelvin before you start using it;

```go
import "github.com/mertcandav/kelvin"
```

A Kelvin database is always uses ``.klvn`` extension. \
The ``Open`` or ``OpenSafe`` functions are used to create or use an existing Kelvin database. \
The ``OpenSafe`` function is recommended if you want to encrypt the database.

```go
k := kelvin.Open[Car]("cars.klvn", kelvin.InMemory)
```

In the example above, you connect to an unencrypted Kelvin database.
If it is encrypted you will get an error.
If there is no database in the specified file path, it will be created.
Datas are `Car` structure.

There are two modes, the example above uses in-memory mode.

## Modes

Kelvin has two modes: in-memory and strict. \
The in-memory mode stores all content in memory. If you write content to disk, you must write manual. \
The strict mode stores all content in disk.

- InMemory
  - Stores data in memory
  - More performance and faster access times
  - Increases memory usage
  - Data is written to disk manually
- Strict
  - Stores data in disk
  - Reduced performance and access times due to read and write operations
  - More efficient memory usage
  - Data is written to disk every time

## Check Database is NoWrite Mode

```go
nw := k.IsNoWrite()
if nw {
    // NoWrite mode
}
```

## Write Buffer to Disk
In in-memory mode, you have to write the memory to disk yourself. To do this, the ``Commit`` function is used. \
The ``Commit`` function is available for only in-memory mode.

```go
k.Commit()
```

## Get Collection

The ``GetCollection`` function is used to get all data of database. \
Returns immutable copy of collection, but not deep copy.

```go
coll := k.GetCollection()
```

## Insert Data
The ``Insert`` function is used to insert data.

```go
k.Insert(
    Car{Brand: "Ferrari", Model: "330 P4"},
    Car{Brand: "Ford", Model: "GT40"})
```
