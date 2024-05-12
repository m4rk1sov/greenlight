package tests

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"greenlight.m4rk1sov.github.com/internal/data"
	"greenlight.m4rk1sov.github.com/internal/validator"
	"testing"
	"time"
)

func TestMovieModelGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"id", "created_at", "title", "year", "runtime", "genres", "version"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, time.Now(), "Test Movie", 2019, 123, pq.Array([]string{"Action"}), 1)

	mock.ExpectQuery("SELECT (.+) FROM movies WHERE").WithArgs(1).WillReturnRows(rows)
	m := &data.MovieModel{DB: db}

	movie, err := m.Get(1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if movie == nil || movie.ID != 1 {
		t.Errorf("Expected movie with ID 1, got %v", movie)
	}
}

func TestValidateMovieYear(t *testing.T) {
	v := validator.New()
	movie := &data.Movie{
		Title:   "Test",
		Year:    2026,
		Runtime: 200,
		Genres:  []string{"testing"},
	}
	data.ValidateMovie(v, movie)
	if v.Valid() {
		t.Errorf("Expected invalid due to year being in future")
	}
}

func TestMovieModelDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM movies WHERE").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	m := &data.MovieModel{DB: db}

	err = m.Delete(1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestCreateUserInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	newUser := data.UserInfo{
		Name:      "Alice",
		Surname:   "Smith",
		Email:     "alice@example.com",
		Activated: false,
	}

	password := "password123"
	newUser.PasswordHash = data.PasswordHash{}
	if err := newUser.PasswordHash.Set(password); err != nil {
		t.Fatalf("Failed to set password: %s", err)
	}

	// Prepare the expected return values from the RETURNING clause
	mock.ExpectQuery("INSERT INTO user_info").
		WithArgs(
			newUser.Name, newUser.Surname, newUser.Email, newUser.PasswordHash.Hash, newUser.Activated,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_role", "version"}).
			AddRow(1, time.Now(), time.Now(), "user", 1))

	m := &data.UserInfoModel{DB: db}
	err = m.Insert(&newUser)
	if err != nil {
		t.Errorf("Failed to insert mock data. Unexpected error: %s", err)
	}
	if newUser.ID != 1 { // Assuming your method sets the ID on the UserInfo object
		t.Errorf("Expected user ID to be 1, got %d", newUser.ID)
	}
}

func TestPasswordHashingAndVerification(t *testing.T) {
	// Define the plaintext password
	plaintextPassword := "correcthorsebatterystaple"

	// Create an instance of passwordHash
	p := &data.PasswordHash{}

	// Hash the plaintext password
	err := p.Set(plaintextPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %s", err)
	}

	// Check if the hash matches the plaintext password
	match, err := p.Matches(plaintextPassword)
	if err != nil {
		t.Errorf("Error while checking password match: %s", err)
	}
	if !match {
		t.Errorf("Hash does not match the original password")
	}

	// Additionally, verify that the hash does not match a wrong password
	wrongPassword := "wrongpassword"
	wrongMatch, err := p.Matches(wrongPassword)
	if err != nil {
		t.Errorf("Error while checking password mismatch: %s", err)
	}
	if wrongMatch {
		t.Errorf("Hash matches an incorrect password")
	}
}

func TestInsertDepartment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the query
	mock.ExpectQuery("INSERT INTO departmentInfo").
		WithArgs("Development", 15, "Jaden Smith", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	m := &data.DepartmentInfoModel{DB: db}

	// Call the function that performs the insert
	department := data.DepartmentInfo{
		DepartmentName:     "Development",
		StaffQuantity:      15,
		DepartmentDirector: "Jaden Smith",
		Module_Info:        1,
	}
	err = m.Insert(&department)
	if err != nil {
		t.Errorf("Error when inserting department info: %s", err)
	}
	if department.ID != 1 { // Assuming your method sets the ID on the department object
		t.Errorf("Expected department ID to be 1, got %d", department.ID)
	}
}

func TestGetModule_infoByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	columns := []string{"id", "created_at", "updated_at", "moduleName", "moduleDuration", "examType", "version"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, time.Now(), time.Now(), "Intro to Go", 30, "Multiple Choice", 1)

	mock.ExpectQuery("SELECT (.+) FROM module_info WHERE id = ?").WithArgs(1).WillReturnRows(rows)
	m := &data.Module_infoModel{DB: db}

	module, err := m.Get(1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if module == nil || module.ID != 1 {
		t.Errorf("Expected module with ID 1, got %v", module)
	}
}

func TestDeleteModule_info(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("DELETE FROM module_info WHERE id = ?").WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulating one row affected

	m := &data.Module_infoModel{DB: db}
	err := m.Delete(1)
	if err != nil {
		t.Errorf("Unexpected error during deletion: %s", err)
	}
}
