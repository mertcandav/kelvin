package kelvin

import (
	"testing"
)

type Employee struct {
	Name    string
	Surname string
	Title   string
	Age     uint8
}

func BenchmarkWhere(b *testing.B) {
	const AGE = 40
	db := OpenNW[Employee]()
	db.Insert(
		Employee{"Jane", "Ace", "Computer Engineer", 36},
		Employee{"Fred", "Addison", "Computer Scientist", 35},
		Employee{"Peter", "Ferrell", "UI/UX Designer", 23},
		Employee{"Michael", "Baker", "Frontend Developer", 63},
		Employee{"Boris", "Ginott", "DevOps", 35},
		Employee{"Philip", "Donne", "Frontend Developer", 63},
		Employee{"Susan", "Davis", "Backend Developer", 23},
		Employee{"Emma", "Walker", "Software Engineer", 55},
		Employee{"Mike", "Miller", "Backend Developer", 46},
		Employee{"Tommy", "Thompson", "DevOps", 26},
		Employee{"Angel", "Turner", "Data Scientist", 21},)
	for i := 0; i < b.N; i++ {
		db.Where(func(t *Employee) bool { return t.Age > AGE })
	}
}

func BenchmarkUWhere(b *testing.B) {
	const AGE = 40
	db := OpenNW[Employee]()
	db.Insert(
		Employee{"Jane", "Ace", "Computer Engineer", 36},
		Employee{"Fred", "Addison", "Computer Scientist", 35},
		Employee{"Peter", "Ferrell", "UI/UX Designer", 23},
		Employee{"Michael", "Baker", "Frontend Developer", 63},
		Employee{"Boris", "Ginott", "DevOps", 35},
		Employee{"Philip", "Donne", "Frontend Developer", 63},
		Employee{"Susan", "Davis", "Backend Developer", 23},
		Employee{"Emma", "Walker", "Software Engineer", 55},
		Employee{"Mike", "Miller", "Backend Developer", 46},
		Employee{"Tommy", "Thompson", "DevOps", 26},
		Employee{"Angel", "Turner", "Data Scientist", 21},)
	for i := 0; i < b.N; i++ {
		db.UWhere(func(t *Employee) bool { return t.Age > AGE })
	}
}
