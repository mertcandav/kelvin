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
k, err := kelvin.Open[Car]("cars.klvn", kelvin.InMemory)
if err != nil {
    // error occurs
}
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

## Write Buffer to Disk
In in-memory mode, you have to write the memory to disk yourself. To do this, the ``Commit`` function is used. \
The ``Commit`` function is available for only in-memory mode.

```go
err := k.Commit()
if err != nil {
    // committing is failed
}
```

## Insert Data
The ``Insert`` function is used to insert data.

```go
err := k.Insert(
    Car{Brand: "Ferrari", Model: "330 P4"},
    Car{Brand: "Ford", Model: "GT40"})
```
