package main

import (
	"errors"
	"testing"
)

// Relationship определяет положение в семье.
type Relationship string

// Возможные роли в семье.
const (
	Father      = Relationship("father")
	Mother      = Relationship("mother")
	Child       = Relationship("child")
	GrandMother = Relationship("grandMother")
	GrandFather = Relationship("grandFather")
)

// Family описывает семью.
type Family struct {
	Members map[Relationship]Person
}

// Person описывает конкретного человека в семье.
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

var (
	// ErrRelationshipAlreadyExists возвращает ошибку, если роль уже занята.
	// Подробнее об ошибках поговорим в девятой теме: «Errors, log».
	ErrRelationshipAlreadyExists = errors.New("relationship already exists")
)

// AddNew добавляет нового члена семьи.
// Если в семье ещё нет людей, создаётся пустой map.
// Если роль уже занята, метод выдаёт ошибку.
func (f *Family) AddNew(r Relationship, p Person) error {
	if f.Members == nil {
		f.Members = map[Relationship]Person{}
	}
	if _, ok := f.Members[r]; ok {
		return ErrRelationshipAlreadyExists
	}
	f.Members[r] = p
	return nil
}

func TestAddNew(t *testing.T) {
	tests := []struct {
		name            string
		value           Family
		addPerson       Person
		addRelationship Relationship
		want            error
	}{
		{
			name: "positive",
			value: Family{
				Members: map[Relationship]Person{
					Child: {
						FirstName: "Vasya",
						LastName:  "Petrov",
						Age:       12,
					}},
			},
			addPerson: Person{
				FirstName: "Vitya",
				LastName:  "Petrov",
				Age:       35,
			},
			addRelationship: Father,
			want:            nil,
		},
		{
			name: "with error",
			value: Family{
				Members: map[Relationship]Person{
					Child: {
						FirstName: "Vasya",
						LastName:  "Petrov",
						Age:       12,
					}},
			},
			addPerson: Person{
				FirstName: "Vitya",
				LastName:  "Petrov",
				Age:       13,
			},
			addRelationship: Child,
			want:            ErrRelationshipAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := tt.value.AddNew(tt.addRelationship, tt.addPerson); res != tt.want {
				t.Errorf("AddNew() = %v, want %v", res, tt.want)
			}
		})
	}
}
