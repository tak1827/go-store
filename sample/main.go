package main

import (
	"fmt"

	"github.com/tak1827/go-store/store"
	"github.com/tak1827/go-store/store/sample/pb"
)

var PREFIX_PERSON = []byte(".person")

func main() {
	leveldb, err := store.NewLevelDB("")
	if err != nil {
		panic(err.Error())
	}

	personStore := store.NewPrefixStore(leveldb, PREFIX_PERSON)

	person := pb.Person{
		Name:  "tom",
		Id:    123,
		Email: "tom@mail.com",
		Phones: []*pb.Person_PhoneNumber{
			&pb.Person_PhoneNumber{
				Number: "090-1111-2222",
				Type:   pb.Person_MOBILE,
			},
			&pb.Person_PhoneNumber{
				Number: "012-345-678",
				Type:   pb.Person_HOME,
			},
		},
	}

	key := person.StoreKey()
	value, err := person.Marshal()
	if err != nil {
		panic(err.Error())
	}

	if err = personStore.Put(key, value); err != nil {
		panic(err.Error())
	}

	v, err := personStore.Get(key)
	if err != nil {
		panic(err.Error())
	}

	var p pb.Person
	err = p.Unmarshal(v)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("unmarshalled person: %v", p)

	if !p.Equal(person) {
		panic(fmt.Sprintf("unexpected get value, get: %v, want: %v", p, person))
	}

	if g, w := pb.PersonIdFromStoreKey(key), person.GetId(); g != w {
		panic(fmt.Sprintf("unexpected recoverd id, get: %d, want: %d", g, w))
	}
}
