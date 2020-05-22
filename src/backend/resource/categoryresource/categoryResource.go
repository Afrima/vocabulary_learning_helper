package categoryresource

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/afrima/vocabulary_learning_helper/src/backend/entity/categoryentity"
	"github.com/afrima/vocabulary_learning_helper/src/backend/resource"
	"github.com/afrima/vocabulary_learning_helper/src/backend/utility"
)

func Init(r *mux.Router) {
	const path = "/category"
	r.HandleFunc(path, get).Methods(http.MethodGet)
	r.HandleFunc(path+"/{id}", getByID).Methods(http.MethodGet)
	r.HandleFunc(path, insert).Methods(http.MethodPost)
	r.Handle(path+"/{id}", resource.IsAuthorized(deleteCategory)).Methods(http.MethodDelete)
}

func getByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.Header().Set(utility.ContentType, utility.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	categoryList, err := categoryentity.GetCategoryByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
		return
	}
	if err = json.NewEncoder(w).Encode(categoryList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(utility.ContentType, utility.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	claims, _ := resource.GetTokenClaims(r)
	userName := claims["userName"].(string)
	categoryList, err := categoryentity.GetCategory(userName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
		return
	}
	if err = json.NewEncoder(w).Encode(categoryList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
	}
}

func insert(w http.ResponseWriter, r *http.Request) {
	body, err := getCategoryFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		log.Print(err)
		return
	}
	claims, _ := resource.GetTokenClaims(r)
	body.Owner = strings.ToLower(claims["userName"].(string))
	if err = body.Insert(); err != nil {
		switch err.(type) {
		case categoryentity.Error:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Print(err)
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set(utility.ContentType, utility.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
		return
	}
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := categoryentity.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
		return
	}
	w.Header().Set(utility.ContentType, utility.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
		return
	}
}

func getCategoryFromBody(r *http.Request) (categoryentity.Category, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return categoryentity.Category{}, err
	}
	var returnValue categoryentity.Category
	if err = json.Unmarshal(reqBody, &returnValue); err != nil {
		return categoryentity.Category{}, err
	}
	return returnValue, nil
}
