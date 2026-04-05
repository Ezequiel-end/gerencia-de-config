package user

import (
	"errors"
	"testing"
)

type fakeRepo struct {
	users  []User
	nextID int
	fail   error
}

func (f *fakeRepo) GetAll() ([]User, error) {
	if f.fail != nil {
		return nil, f.fail
	}
	return f.users, nil
}

func (f *fakeRepo) GetByID(id int) (*User, error) {
	if f.fail != nil {
		return nil, f.fail
	}
	for _, u := range f.users {
		if u.ID == id {
			cp := u
			return &cp, nil
		}
	}
	return nil, errors.New("user not found")
}

func (f *fakeRepo) Create(u User) (*User, error) {
	if f.fail != nil {
		return nil, f.fail
	}
	if f.nextID == 0 {
		f.nextID = 1
	}
	u.ID = f.nextID
	f.nextID++
	f.users = append(f.users, u)
	cp := u
	return &cp, nil
}

func (f *fakeRepo) Update(id int, u User) (*User, error) {
	if f.fail != nil {
		return nil, f.fail
	}
	for i, x := range f.users {
		if x.ID == id {
			u.ID = id
			f.users[i] = u
			cp := u
			return &cp, nil
		}
	}
	return nil, errors.New("user not found")
}

func (f *fakeRepo) Delete(id int) error {
	if f.fail != nil {
		return f.fail
	}
	for i, u := range f.users {
		if u.ID == id {
			f.users = append(f.users[:i], f.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func TestUserService_ListUsers(t *testing.T) {
	repo := &fakeRepo{users: []User{{ID: 1, Name: "a", Email: "a@x"}}}
	s := NewUserService(repo)
	list, err := s.ListUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "a" {
		t.Fatalf("got %+v", list)
	}
}

func TestUserService_CreateAndGet(t *testing.T) {
	repo := &fakeRepo{}
	s := NewUserService(repo)
	created, err := s.CreateUser(User{Name: "b", Email: "b@x"})
	if err != nil || created.ID != 1 {
		t.Fatalf("create: %+v err=%v", created, err)
	}
	got, err := s.GetUser(1)
	if err != nil || got.Email != "b@x" {
		t.Fatalf("get: %+v err=%v", got, err)
	}
}
