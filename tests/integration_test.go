package tests

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"greenlight.m4rk1sov.github.com/internal/data"
	"testing"
	"time"
)

func TestMovieLifecycle(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	// Step 1: Insert a new movie
	mock.ExpectQuery("INSERT INTO movies").
		WithArgs("New Movie", 2021, 120, pq.Array([]string{"Drama"})).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "version"}).
			AddRow(1, time.Now(), 1))

	// Step 2: Retrieve the movie
	columns := []string{"id", "created_at", "title", "year", "runtime", "genres", "version"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, time.Now(), "New Movie", 2021, 120, pq.Array([]string{"Drama"}), 1)
	mock.ExpectQuery("SELECT (.+) FROM movies WHERE").WithArgs(1).WillReturnRows(rows)

	// Step 3: Update the movie
	mock.ExpectQuery("UPDATE movies SET title = \\$1, year = \\$2, runtime = \\$3, genres = \\$4, version = version \\+ 1 WHERE id = \\$5 AND version = \\$6 RETURNING version").
		WithArgs("Updated Movie", 2021, 130, pq.Array([]string{"Drama"}), 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(2))

	// Step 4: Delete the movie
	mock.ExpectExec("DELETE FROM movies WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	m := &data.MovieModel{DB: db}

	// Insert
	err := m.Insert(&data.Movie{Title: "New Movie", Year: 2021, Runtime: 120, Genres: []string{"Drama"}, Version: 1})
	if err != nil {
		t.Fatalf("Failed to insert movie: %s", err)
	}

	// Retrieve
	movie, err := m.Get(1)
	if err != nil || movie == nil {
		t.Fatalf("Failed to retrieve movie: %s", err)
	}

	// Update
	movie.Title = "Updated Movie"
	movie.Runtime = 130
	err = m.Update(movie)
	if err != nil {
		t.Fatalf("Failed to update movie: %s", err)
	}

	// Delete
	err = m.Delete(1)
	if err != nil {
		t.Fatalf("Failed to delete movie: %s", err)
	}
}

func TestUserRegistrationAndLogin(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	newUser := data.UserInfo{
		Name: "Bob", Surname: "Smith", Email: "bob.smith@example.com", Activated: false,
	}
	password := "password123"
	newUser.PasswordHash = data.PasswordHash{}
	if err := newUser.PasswordHash.Set(password); err != nil {
		t.Fatalf("Failed to set password: %s", err)
	}

	// Mock user insertion
	mock.ExpectQuery("INSERT INTO user_info").
		WithArgs(
			newUser.Name, newUser.Surname, newUser.Email, newUser.PasswordHash.Hash, newUser.Activated,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_role", "version"}).
			AddRow(1, time.Now(), time.Now(), "user", 1))

	// Mock user retrieval
	mock.ExpectQuery("SELECT (.+) FROM user_info WHERE email = ?").
		WithArgs(newUser.Email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "fname", "sname", "email", "password_hash", "user_role", "activated", "version",
		}).
			AddRow(1, time.Now(), time.Now(), newUser.Name, newUser.Surname, newUser.Email, newUser.PasswordHash.Hash, "user", newUser.Activated, 1))

	m := &data.UserInfoModel{DB: db}

	// Insert user
	err := m.Insert(&newUser)
	if err != nil {
		t.Fatalf("Failed to insert user: %s", err)
	}

	// Attempt to login
	retrievedUser, err := m.GetByEmail(newUser.Email)
	if err != nil {
		t.Fatalf("Failed to retrieve user: %s", err)
	}

	// Check password
	correct, _ := retrievedUser.PasswordHash.Matches(password)
	if !correct {
		t.Errorf("Password does not match")
	}
}

func TestDepartmentAndModulesAssociation(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	// Mock department insertion with module info
	mock.ExpectQuery("INSERT INTO departmentInfo").
		WithArgs("IT", 20, "Gloria Smith", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	columns := []string{"id", "created_at", "updated_at", "moduleName", "moduleDuration", "examType", "version"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, time.Now(), time.Now(), "Intro to Go", 30, "Multiple Choice", 1)

	mock.ExpectQuery("SELECT id, created_at, updated_at, moduleName, moduleDuration, examType, version FROM module_info WHERE id = \\$1").
		WithArgs(1).WillReturnRows(rows)

	m := &data.DepartmentInfoModel{DB: db}
	moduleModel := &data.Module_infoModel{DB: db}

	// Insert department
	department := data.DepartmentInfo{
		DepartmentName: "IT", StaffQuantity: 20, DepartmentDirector: "Gloria Smith", Module_Info: 1,
	}
	err := m.Insert(&department)
	if err != nil {
		t.Fatalf("Failed to insert department: %s", err)
	}

	// Retrieve associated module
	module, err := moduleModel.Get(1)
	if err != nil {
		t.Fatalf("Failed to retrieve module: %s", err)
	}

	if module.ModuleName != "Intro to Go" {
		t.Errorf("Expected 'Intro to Go', got '%s'", module.ModuleName)
	}
}

func TestDeleteUserAndTokens(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	// Mock user deletion
	mock.ExpectExec("DELETE FROM user_info WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulating that one row was affected.

	// Choose the scope you want to delete, for example, `ScopeAuthentication`
	// Mock token deletion for the specified scope
	mock.ExpectExec("DELETE FROM tokens WHERE scope = \\$1 AND user_id = \\$2").
		WithArgs(data.ScopeAuthentication, 1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulating that one row was affected.

	userInfoModel := &data.UserInfoModel{DB: db}
	tokenModel := &data.TokenModel{DB: db}

	// Delete user
	err := userInfoModel.Delete(1)
	if err != nil {
		t.Fatalf("Failed to delete user: %s", err)
	}

	// Delete tokens
	err = tokenModel.DeleteAllForUser(data.ScopeAuthentication, 1)
	if err != nil {
		t.Fatalf("Failed to delete tokens for user: %s", err)
	}
}
