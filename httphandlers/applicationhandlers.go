package httphandlers

import (
	"encoding/json"
	"net/http"
	models "quots/models"
	utils "quots/utils"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateApplication(w http.ResponseWriter, r *http.Request) {
	var application models.Application
	_ = json.NewDecoder(r.Body).Decode(&application)
	secret := utils.RandStringBytesMaskImprSrcUnsafe(18)
	application.AppSecret = secret
	_, err := dbAppImpl.CreateApplication(application)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusConflict,
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(application)
	}
}

func GetApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	application, err := dbAppImpl.GetApplicationBiId(id)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(application)
	}
}

func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	deletedCount, err := dbAppImpl.DeleteApplicationBiId(id)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(deletedCount)
	}
}

func UpdateApplicationSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	application, err := dbAppImpl.GetApplicationBiId(id)
	application.AppSecret = utils.RandStringBytesMaskImprSrcUnsafe(16)
	applicatio, err := dbAppImpl.UpdateApp(application)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(applicatio)
	}
}

func GetApplications(w http.ResponseWriter, r *http.Request) {
	min := r.FormValue("min")
	max := r.FormValue("max")
	minval, err := strconv.ParseInt(min, 10, 64)
	maxval, err := strconv.ParseInt(max, 10, 64)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	counted, applications, err := dbAppImpl.GetAllApps(minval, maxval)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		t := strconv.FormatInt(counted, 10)
		w.Header().Add("total", t)
		// w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, Authorization, Client-id, Client-secret, Total, total")
		// w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, DELETE, POST, OPTIONS")
		json.NewEncoder(w).Encode(applications)
	}
}

func UpdateApplication(w http.ResponseWriter, r *http.Request) {
	var applicatio models.Application
	_ = json.NewDecoder(r.Body).Decode(&applicatio)
	_, err := dbAppImpl.UpdateApp(applicatio)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(applicatio)
	}
}
