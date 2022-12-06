# Kelvin
Kelvin is document oriented database written in pure Go.

## Motivation && Reasons
Go has generic support with current versions. \
I want to see if a useful document oriented database can be written for Go using generic types. \
Therefore, this repository contains some experimentation. \
If a good implementation comes out, I hope this repository will allow Go developers to have a tool like [TinyDB](https://github.com/msiemens/tinydb) or [LiteDB](https://github.com/mbdavid/LiteDB).

## Example Code

```go
package main

import "github.com/mertcandav/kelvin"

type Employee struct {
    Name   string
    Title  string
    Salary float64
}

func main() {
    db := kelvin.Open[Employee]("employees.klvn", kelvin.InMemory)
    db.Insert(
        Employee{"Jane SMITH", "Data Scientist", 10250},
        Employee{"Chris DAVIS", "Software Engineer", 8000})
}
```

## Query

**Don't Learn a New Query Language, Keep Writing Go!**

```go
func GetHighSalaryEngineers(db kelvin.Kelvin[Employee]) []Employee {
    return db.Where(func(e Employee) bool {
        return e.Salary > 15000 &&
               (e.Title == "Software Engineer" || e.Title == "Data Engineer")
    })
}
```

## Safe Transactions

Kelvin supports ACID and provides thread-safe functions.

```go
package main

import (
    "fmt"
    "sync"

    "github.com/mertcandav/kelvin"
)

type Count struct { N int }

func Increase(wg *sync.WaitGroup, db kelvin.Kelvin[Count]) {
    db.Map(func(c *Count) { c.N++ })
    wg.Done()
}

func main() {
    db := kelvin.OpenNW[Count]()
    db.Insert(Count{N: 0})

    wg := sync.WaitGroup{}
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go Increase(&wg, db)
    }

    wg.Wait()
    fmt.Println(db.GetCollection()[0].N)
}
```

Kelvin, in the example above, the result of 1000 goroutines will always be N = 1000 consistently and will be very efficient and quick to respond transactions in performance.

## Design Principles
- As efficient and performant as possible
- Take care to keep type safety
- Implement Kelvin as readable && elegant
- Follow ACID principles
- Provide thread-safety
- No 3rd-party dependencies

## Contributing

Even the smallest contribution is greatly appreciated. \
After forking the repository please make your changes and open a PR.

Please follow:
- Your commit messages are explanatory
- Your PR is not huge
- Your changes are tested
- Your changes follows design principles

## License
The Kelvin is distributed under the terms of the BSD 3-Clause license. <br>
[See License Details](LICENSE)
