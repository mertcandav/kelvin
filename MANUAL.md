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
db := kelvin.Open[Employee]("employees.klvn", kelvin.InMemory)
```

In the example above, you connect to an unencrypted Kelvin database.
If it is encrypted you will get an error.
If there is no database in the specified file path, it will be created.
Datas are `Employee` structure.

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
nw := db.IsNoWrite()
if nw {
    // NoWrite mode
}
```

## Write Buffer to Disk
In in-memory mode, you have to write the memory to disk yourself. To do this, the ``Commit`` function is used. \
The ``Commit`` function is available for only in-memory mode.

```go
db.Commit()
```

## Get Collection

The ``GetCollection`` function is used to get all data of database. \
Returns immutable copy of collection, but not deep copy.

```go
coll := db.GetCollection()
```

## Map Function

Map iterates into all collection and commits changes. \
Does not nothing if handler is nil.

This method can be useful if you want to manipulate data based on a certain condition.

> Collection copies are not deep immutable copy. \
> So if you make changes any mutable field of T, you can change original collection.

```go
db.Map(func(e *Employee) {
    switch e.Title {
    case "Software Engineer":
        e.Salary = (e.Salary*120) / 100
    case "Computer Scientist":
        e.Salary = (e.Salary*130) / 100
    case "Data Scientist":
        e.Salary = (e.Salary*125) / 100
    }
})
```

The example above gives an increase in the salaries of the employees according to their job titles.

## Insert Data

The ``Insert`` function is used to insert data.

```go
db.Insert(
    Car{Name: "James SMITH", Title: "Software Engineer", Salary: 12500},
    Car{Name: "Linda JONES", Title: "Data Engineer", Salary: 10750})
```

## Filter Data

The ``Where`` function is used to get collection with filter. \
Returns nil if handler is nil.

> Collection copies are not deep immutable copy. \
> So if you make changes any mutable field of T, you can change original collection.

```go
employees := db.Where(func(e Employee) bool { return e.Salary > 8000 })
```

The example above returns a collection that contains only employees with a salary higher than 4000.
