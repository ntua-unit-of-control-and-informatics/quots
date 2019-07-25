package models

type User struct {
	Id       string  `json:"id" bson:"_id"`
	Email    string  `json:"email" bson:"email"`
	Username string  `json:"username" bson:"username"`
	Credits  float64 `json:"credits" bson:"credits"`
	Spenton  []Spent `json:"spenton" bson:"spenton"`
}

type CanProceed struct {
	Userid  string `json:"userid" bson:"userid"`
	Proceed bool   `json:"proceed"`
}

type Spent struct {
	Appid string             `json:"appid" bson:"appid"`
	Usage map[string]float64 `json:"usage"  bson:"usage"`
}

type Application struct {
	Id         string             `json:"id" bson:"_id"`
	UsagesCost map[string]float64 `json:"usagescost" bson:"usagescost"`
	UsageTypes []string           `json:"usagetypes" bson:"usagetypes"`
	AppSecret  string             `json:"clientsecret" bson:"clientsecret"`
	Enabled    bool               `json:"enabled" bson:"enabled"`
	BaseURLS   []string           `json:"baseURLS" bson:"baseURLS"`
}

type ErrorReport struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type AppUsed struct {
	Appid  string `json:"appid" bson:"appid"`
	Usage  Usage  `json:"usage" bson:"usage"`
	Userid string `json:"userid" bson:"userid"`
}

type Usage struct {
	Usagetype string  `json:"usagetype" bson:"usagetype"`
	Used      float64 `json:"used" bson:"used"`
}

type JWTSecret struct {
	Signature string `bson:"signature"`
}
