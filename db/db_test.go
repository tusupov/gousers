package db

import (
	"testing"
)

var testStore = NewDB()

func TestDb_NewGet(t *testing.T) {

	usrCnt := int(1e3)
	usrAmount := uint64(1e3)

	for i := 0; i < usrCnt; i++ {
		_, err := testStore.New(usrAmount)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < usrCnt; i++ {
		usr, ok := testStore.Get(uint64(i + 1))
		if !ok {
			t.Fatalf("Could not find the user, but should have been in the database")
		}
		if usr.Amount != usrAmount {
			t.Fatalf("Invalid balance")
		}
	}

}

func TestDb_Transfer(t *testing.T) {

	userAmount := uint64(1e3)

	fromUser, err := testStore.New(userAmount)
	if err != nil {
		t.Fatal(err)
	}

	toUser, err := testStore.New(userAmount)
	if err != nil {
		t.Fatal(err)
	}

	err = testStore.Transfer(fromUser.Id, toUser.Id, userAmount)
	if err != nil {
		t.Fatal(err)
	}

	err = testStore.Transfer(fromUser.Id, toUser.Id, userAmount)
	if err == nil {
		t.Fatal("An error was expected about amount not enough")
	}

}

func BenchmarkDb_New(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testStore.New(1e3)
		}
	})

}

func BenchmarkDb_Get(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testStore.Get(1)
		}
	})

}

func BenchmarkDb_Transfer(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testStore.Transfer(10, 20, 10)
		}
	})

}

func BenchmarkDb_Parallel(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testStore.New(1e3)
			testStore.Get(1)
			testStore.Transfer(1, 2, 10)
		}
	})

}
