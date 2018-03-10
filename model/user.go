package model

type User struct {
	Ip           string         `json:"ip" from:"ip"`
	Name         string         `json:"name" from:"name"`
	Status       int            `json:"status" from:"status"`  //0:离线；1：在线
}
