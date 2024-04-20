package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Movies         MovieModel
	Module_info    Module_infoModel
	DepartmentInfo DepartmentInfoModel
	//// Set the Movies field to be an interface containing the methods that both the
	//// 'real' model and mock model need to support.
	//Movies interface {
	//	Insert(movie *Movie) error
	//	Get(id int64) (*Movie, error)
	//	Update(movie *Movie) error
	//	Delete(id int64) error
	//	GetAll()
	//}
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies:         MovieModel{DB: db},
		Module_info:    Module_infoModel{DB: db},
		DepartmentInfo: DepartmentInfoModel{DB: db},
	}
}

//// Create a helper function which returns a Models instance containing the mock models
//// only.
//func NewMockModels() Models {
//	return Models{
//		Movies: MockMovieModel{},
//	}
//}
