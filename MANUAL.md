# Kelvin Manual

Kelvin is a document-oriented database that uses the generic support of Go.


## Getting Started

To use Kelvin in your Go project, you must have the Kelvin module. \
Import Kelvin before you start using it;

```go
import "github.com/mertcandav/kelvin"
```

A Kelvin database is always uses ``.klvn`` extension. \
The ``Open``, ``OpenSafe`` or ``OpenNW`` functions are used to create or use an existing Kelvin database. \
The ``OpenSafe`` function is recommended if you want to encrypt the database. \
The ``OpenNW`` function is shortcut for no-write and in-memory mode.

```go
db := kelvin.Open[Employee]("employees.klvn", kelvin.InMemory)
```

```go
db := kelvin.OpenNW[Employee]()
```

In the example above, you connect to an unencrypted Kelvin database.
If it is encrypted you will get an error.
If there is no database in the specified file path, it will be created.
Datas are `Employee` structure.

There are two modes, the example above uses in-memory mode.

> **Warning** \
> The data type must be structure. \
> Structure or pointer to structure.

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
    Employee{Name: "James SMITH", Title: "Software Engineer", Salary: 12500},
    Employee{Name: "Linda JONES", Title: "Data Engineer", Salary: 10750})
```

## Filter Data

The ``Where`` function is used to get collection with filter. \
Returns nil if handler is nil.

```go
employees := db.Where(func(e *Employee) bool { return e.Salary > 8000 })
```

The example above returns a collection that contains only employees with a salary higher than 8000.

## Unsafe Functions

Unsafe functions are unsafe.
They are mostly the unsafe equivalents of safe functions for the sake of performance.
Their main insecurities come from the fact that the buffer can be changed.
They do not use an immutable copy, so changes made to a mutable data may affect the original data, which may break data consistency.

Unsafe functions starts with `U` prefix.

List of unsafe functions:
- ``UWhere`` is unsafe equavalents of ``Where``

## Static Typing

The original Kelvin database structure is not directly provided for safety reasons. \
Use the ``Kelvin`` interface for static typing.

```go
var db kelvin.Kelvin[Employee]
db = kelvin.OpenNW[Employee]()
```

## How Kelvin Handling Data Safety?

Kelvin is designed to be ACID compliant and thread-safe to guarantee data safety. \
Let's see exactly how Kelvin handling these approaches;


### Atomicity

Kelvin behaves in such a way that all actions in a transaction are valid only when all actions are successful. \
It always uses a copy of the data to behave this way.

While in InMemory mode, a copy of the buffer is created and all actions that cause a change, such as writes, are executed from this copy. \
Any error condition does not make the data inconsistent as the changes are made on a copy. This way all changes are overridden. \
If all transactions are successful, the copy of the buffer should become the current state of the database. Therefore, the buffer is set to this copy.

While in Strict mode, data is read from disk, every operation like in InMemory mode is expected to be successful. \
If all transactions are successful, it is written to disk.

### Isolation

Transactions can occur on the same moment and may need to be executed concurrently. \
Kelvin works to give the impression of concurrency, but it's not exactly simultaneous. Each transaction has to wait for another. \
The transactions do not interfere with each other, but are executed sequentially to give a consistent result. \
Each transaction is independent and isolated from each other.

Kelvin uses Mutex locks to achieve this. \
Each Kelvin database locks the data collection when needed so that it can be accessed by a single transaction.

When reading from the database, the data collection is locked, so it is certain that there will be no change while reading. Only the reading thread can access it. After the read operation is done, the lock is unlocked. \
Even if the read operation includes a different transaction, if it will not cause a change in the database, it unlocks after reading. \
In this way, concurrency is ensured as much as possible.

For write operations, the entire data collection remains locked from the beginning to the end of the transaction. \
Since the data will change, all read operations are also suspended. All other writes are also suspended to use up-to-date data. \
A write operation locks the data collection when it is certain that the write operation will occur, or at least one read operation. \
Because write operations are executed with a copy, no intermediate changes will appear in the original data collection. \
In this way, it tries to make the highest amount of concurrency possible.

Here is a simple example of how concurrency is handled;

The most common scenario in this regard is based on banks, and we will do the same. \
Let's say you have a bank account with account number 9056, there is $900 in this account. \
Two transactions that write to this account were attempted to be executed concurrent.

We have two transactions; T1 and T2

 - T1: deducts $250 from account 9056
 - T2: deducts $300 from accout 9056

To ensure consistency, Kelvin executes these transactions as if they were executed sequentially. \
To do this, it uses the locking mechanisms described above. \
T1 then T2 is executed first, or T2 first then T1. In both cases, the transactions wait for the transaction executed before them to finish. \
It's obvious why it does this: race condition \
In both cases, the account balance must be $350.

If full concurrency was allowed, the data could be inconsistent. \
For example, T1 starts and reads the account balance, result will be $900. \
Immediately after that, T2 starts and reads the account balance, likewise getting $900. \
T1 then calculates that the account balance should be $650 by making 900 - 250. \
Immediately after that T2 then calculates that the account balance should be $600 by making 900 - 300. \
T1 then updates the account balance to $650, after T2 updates to $600. \
At the end of the transaction, $600 is left in the account. It's wrong calculation, inconsistent.

Because Kelvin ranks transactions, consistency is ensured. \
In the same script, the executing transaction is expected to finish first. \
Therefore, the transaction that starts after it is guaranteed to read the current data, so the account balance is always guaranteed to be $350.

### Durability

Kelvin ensures that this change is permanent once the transaction is successful.

It will always write changes to disk as long as it is in Strict mode.
In this way, all changes are preserved. \
The database can be restored and processed again from the last state in which it was saved. \
In this way, it is guaranteed that data will not be lost in case of power failure.

Data is not automatically written to disk while in InMemory mode, but the abilities are the same. \
The last valid database state that you manually saved will be restoreable and operable again.
