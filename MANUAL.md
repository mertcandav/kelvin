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
k, err := kelvin.Open("employees.klvn", kelvin.InMemory)
if err != nil {
    // error occurs
}
```

In the example above, you connect to an unencrypted Kelvin database.
If it is encrypted you will get an error.
If there is no database in the specified file path, it will be created.

There are two modes, the example above uses in-memory mode.

## Modes

Kelvin has two modes: in-memory and strict. \
The in-memory mode stores all content in memory. If you write content to disk, you must write manual. \
The strict mode stores all content in disk.

- In-Memory
  - Stores data in memory
  - More performance and faster access times
  - Increases memory usage
- Strict
  - Stores data in disk
  - Reduced performance and access times due to read and write operations
  - More efficient memory usage
