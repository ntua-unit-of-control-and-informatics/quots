package httphandlers

import (
	"encoding/json"
	"net/http"

	// db "quots/database"
	models "quots/models"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var person models.User
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.Credits = 100
	_, err := dbUserImpl.CreateUser(person)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusConflict,
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(person)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	user, err := dbUserImpl.GetUserById(id)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(user)
	}
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	user, err := dbUserImpl.FindUserByEmail(email)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(user)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	min := r.FormValue("min")
	max := r.FormValue("max")
	email := r.FormValue("email")
	if email != "" {
		user, err := dbUserImpl.FindUserByEmail(email)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusNotFound,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(user)
		}
	} else {
		minval, err := strconv.ParseInt(min, 10, 64)
		maxval, err := strconv.ParseInt(max, 10, 64)
		total, users, err := dbUserImpl.GetUsersPaginated(minval, maxval)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusNotFound,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		} else {
			t := strconv.FormatInt(total, 10)
			w.Header().Add("total", t)
			json.NewEncoder(w).Encode(users)
		}
	}

}

func UpdateCredits(w http.ResponseWriter, r *http.Request) {
	var person models.User
	_ = json.NewDecoder(r.Body).Decode(&person)
	user, err := dbUserImpl.UpdateUserCredits(person)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(user)
	}
}

func UserQuotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	appid := r.FormValue("appid")
	usage := r.FormValue("usage")
	size := r.FormValue("size")
	user, err := dbUserImpl.GetUserById(id)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	if appid != "" && usage != "" && size != "" {

		app, err := dbAppImpl.GetApplicationBiId(appid)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		usagecost := app.UsagesCost[usage]
		if usagecost == 0 {
			err := models.ErrorReport{
				Message: "Cound not find usage",
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		fcost, err := strconv.ParseFloat(size, 64)
		costall := usagecost * fcost
		if user.Credits > costall {
			canProceed := models.CanProceed{
				Userid:  user.Id,
				Proceed: true,
			}
			user.Credits = user.Credits - costall
			index := indexOf(appid, user)
			if index == -1 {
				m := make(map[string]float64)
				m[usage] = costall
				spentOn := models.Spent{
					Appid: appid,
					Usage: m,
				}
				user.Spenton = append(user.Spenton, spentOn)
				dbUserImpl.UpdateUsersSpentOn(user)
			} else {
				user.Spenton[index].Usage[usage] += costall
				dbUserImpl.UpdateUsersSpentOn(user)
			}
			dbUserImpl.UpdateUserCredits(user)
			json.NewEncoder(w).Encode(canProceed)
			return
		} else {
			canProceed := models.CanProceed{
				Userid:  user.Id,
				Proceed: false,
			}
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(canProceed)
			return
		}
	} else {
		json.NewEncoder(w).Encode(user)
	}
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	}
}

func calculateCost(appid string) {

}

func indexOf(element string, user models.User) int {
	for k, v := range user.Spenton {
		if element == v.Appid {
			return k
		}
	}
	return -1 //not found.
}
