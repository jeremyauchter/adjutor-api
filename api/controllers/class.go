package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jeremyauchter/adjutor/api/responses"
	"github.com/jeremyauchter/adjutor/models/products"
	"github.com/jeremyauchter/adjutor/utils/formaterror"
)

func (server *Server) Classes(w http.ResponseWriter, r *http.Request) {

	class := products.Class{}

	classes, err := class.AllClasses(server.database.Product)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, classes)
}

type NewClass struct {
	Name           string `json:"name"`
	DepartmentName string `json:"departmentName"`
}

func (server *Server) CreateClass(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	class := NewClass{}
	err = json.Unmarshal(body, &class)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	department := products.Department{}
	dbClass := products.Class{
		Name: class.Name,
	}
	departmentId, _ := department.DepartmentByName(server.database.Product, class.DepartmentName)

	if departmentId.ID > 0 {

		dbClass.DepartmentID = departmentId.ID
	} else {
		dbClass.Department = products.Department{Name: class.DepartmentName}
	}

	dbClass.PrepareClass()
	err = dbClass.ValidateClass()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	dbClassCreated, err := dbClass.CreateClass(server.database.Product)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, dbClassCreated.ID))
	responses.JSON(w, http.StatusCreated, dbClassCreated)
}

func (server *Server) UpdateClass(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	id := uint32(id64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	// Check if the post exist
	class := products.Class{}
	err = server.database.Product.Debug().Model(products.Class{}).Where("id = ?", id).Take(&class).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("class not found"))
		return
	}

	// Read the data classed
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	classUpdate := products.Class{}
	err = json.Unmarshal(body, &classUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	classUpdate.PrepareClass()
	err = classUpdate.ValidateClass()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	classUpdated, err := classUpdate.UpdateClass(server.database.Product, id)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, classUpdated)
}

func (server *Server) GetClassById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 64)
	id := uint32(id64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	class := products.Class{}

	classReceived, err := class.ClassById(server.database.Product, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, classReceived)
}

func (server *Server) DeleteClass(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid class id given to us?
	id64, err := strconv.ParseUint(vars["id"], 10, 64)
	id := uint32(id64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// // Is this user authenticated?
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	// Check if the class exist
	class := products.Class{}
	err = server.database.Product.Debug().Model(products.Class{}).Where("id = ?", id).Take(&class).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = class.DeleteClass(server.database.Product, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
}
