package user

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgresUserRepository_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip()
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	if err := MigratePostgres(db); err != nil {
		t.Fatal(err)
	}
	_, _ = db.Exec(`TRUNCATE users RESTART IDENTITY`)

	repo := NewPostgresUserRepository(db)

	created, err := repo.Create(User{Name: "int", Email: "int@test.local"})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID < 1 {
		t.Fatalf("id: %d", created.ID)
	}

	all, err := repo.GetAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("GetAll: %+v err=%v", all, err)
	}

	got, err := repo.GetByID(created.ID)
	if err != nil || got.Name != "int" {
		t.Fatalf("GetByID: %+v err=%v", got, err)
	}

	upd, err := repo.Update(created.ID, User{Name: "int2", Email: "int2@test.local"})
	if err != nil || upd.Name != "int2" {
		t.Fatalf("Update: %+v err=%v", upd, err)
	}

	if err := repo.Delete(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(created.ID); err == nil {
		t.Fatal("expected error after delete")
	}
}
