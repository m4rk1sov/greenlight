package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"greenlight.m4rk1sov.github.com/internal/validator"
	"time"
)

type Module_info struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	ModuleName     string    `json:"moduleName"`
	ModuleDuration Runtime   `json:"moduleDuration,omitempty"`
	ExamType       string    `json:"examType"`
	Version        int32     `json:"version"` // The version starts at number 1
}

// To prevent duplication, we can collect the validation checks for a movie into a standalone
// ValidateMovie() function
func ValidateModule(v *validator.Validator, module *Module_info) {
	v.Check(module.ModuleName != "", "moduleName", "must be provided")
	v.Check(len(module.ModuleName) <= 500, "moduleName", "must not be more than 500 bytes long")
	v.Check(module.ModuleDuration != 0, "moduleDuration", "must be provided")
	v.Check(module.ModuleDuration > 0, "moduleDuration", "must be a positive integer")
	v.Check(module.ExamType != "", "examType", "must be provided")
}

// Define a Module_infoModel struct type which wraps a sql.DB connection pool.
type Module_infoModel struct {
	DB *sql.DB
}

// The Insert() method accepts a pointer to a movie struct, which should contain the
// data for the new record.
func (m Module_infoModel) Insert(module_info *Module_info) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
INSERT INTO module_info (moduleName, moduleDuration, examType)
VALUES ($1, $2, $3)
RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the module_info struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{module_info.ModuleName, module_info.ModuleDuration, module_info.ExamType}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the
	// system-generated id, created_at and version values into the module_info struct.
	// Use QueryRowContext() and pass the context as the first argument.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&module_info.ID, &module_info.CreatedAt, &module_info.Version)
}

func (m Module_infoModel) Get(id int64) (*Module_info, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// 1) Define the SQL query for retrieving the movie data.
	// 2) Update the query to return pg_sleep(10) as the first value.
	// 3) Remove the pg_sleep(10) clause.

	query := `
SELECT id, created_at, updated_at, moduleName, moduleDuration, examType, version
FROM module_info
WHERE id = $1`

	// Declare a Movie struct to hold the data returned by the query.
	var module_info Module_info

	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	// 1) Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.

	// 2) Importantly, update the Scan() parameters so that the pg_sleep(10) return value
	// is scanned into a []byte slice.

	// 3) Use the QueryRowContext() method to execute the query, passing in the context
	// with the deadline as the first argument.

	// 4) Remove &[]byte{} from the first Scan() destination.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&module_info.ID,
		&module_info.CreatedAt,
		&module_info.UpdatedAt,
		&module_info.ModuleName,
		&module_info.ModuleDuration,
		&module_info.ExamType,
		&module_info.Version,
	)

	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &module_info, nil
}

func (m Module_infoModel) Update(module_info *Module_info) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	// Add the 'AND version = $6' clause to the SQL query.
	query := `
UPDATE module_info
SET moduleName = $1, moduleDuration = $2, examType = $3, version = version + 1
WHERE id = $4 AND version = $5
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		module_info.ModuleName,
		module_info.ModuleDuration,
		module_info.ExamType,
		module_info.ID,
		module_info.Version, // Add the expected version.
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//// Use the QueryRow() method to execute the query, passing in the args slice as a
	//// variadic parameter and scanning the new version value into the movie struct.
	//return m.DB.QueryRow(query, args...).Scan(&movie.Version)

	// Execute the SQL query. If no matching row could be found, we know the movie
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.

	// Use QueryRowContext() and pass the context as the first argument.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&module_info.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil

}

func (m Module_infoModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
DELETE FROM module_info
WHERE id = $1`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.

	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Create a new GetAll() method which returns a slice of movies. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
// Update the function signature to return a Metadata struct.
func (m Module_infoModel) GetAllModules(moduleName string, examType string, filters Filters) ([]*Module_info, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	// Update the SQL query to include the filter conditions.

	// Use full-text search for the title filter.

	// Add an ORDER BY clause and interpolate the sort column and direction. Importantly
	// notice that we also include a secondary sort on the movie ID to ensure a
	// consistent ordering.

	// Update the SQL query to include the LIMIT and OFFSET clauses with placeholder
	// parameter values.

	// Update the SQL query to include the window function which counts the total
	// (filtered) records.
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, updated_at, moduleName, moduleDuration, examType, version
FROM module_info
WHERE (to_tsvector('simple', moduleName) @@ plainto_tsquery('simple', $1))
ORDER BY %s %s, id ASC
LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []any{moduleName, examType, filters.limit(), filters.offset()}

	// And then pass the args slice to QueryContext() as a variadic parameter.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}

	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	// Declare a totalRecords variable.
	totalRecords := 0
	modules_info := []*Module_info{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var module_info Module_info
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&module_info.ID,
			&module_info.CreatedAt,
			&module_info.UpdatedAt,
			&module_info.ModuleName,
			&module_info.ModuleDuration,
			&module_info.ExamType,
			&module_info.Version,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}

		// Add the Movie struct to the slice.
		modules_info = append(modules_info, &module_info)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.

	}

	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// Include the metadata struct when returning.
	// If everything went OK, then return the slice of movies.
	return modules_info, metadata, nil
}
