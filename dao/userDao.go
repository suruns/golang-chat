package dao

import (
	"log"
	"database/sql"
	"app/model"
)

func GetUserList(db *sql.DB,ip string) ([]model.User,error) {

	rows, err := db.Query("select ip,name,status from user WHERE status = 1 and ip <> ?",ip)

	if err != nil{
		log.Fatalln(err)
	}
	defer rows.Close()

	users := make([]model.User,0)
	for rows.Next() {
		var user model.User
		rows.Scan(&user.Ip,&user.Name,&user.Status)
		users = append(users,user)
	}
	log.Print(users)
	if err =rows.Err();err!=nil{
		log.Fatalln(err)
		return nil,err
	}
	return users,nil
}

func GetUserByIp(db *sql.DB,ip string) (*model.User,bool)  {
	var user model.User
	row := db.QueryRow("SELECT ip,name,status from user where ip = ?",ip)
	err := row.Scan(&user.Ip,&user.Name,&user.Status)
	if err != nil {
		return nil,false
	}else {
		return &user,true
	}
}
func GetUserByName(db *sql.DB,name string) (*model.User,bool)  {
	var user model.User
	row := db.QueryRow("SELECT ip,name,status from user where name = ?",name)
	err := row.Scan(&user.Ip,&user.Name,&user.Status)
	if err != nil {
		return nil,false
	}else {
		return &user,true
	}
}

func UpdateUserByIp(db *sql.DB,ip,name string) (error)  {
	_, err := db.Exec("insert into user(ip,name,status) VALUES (?,?,?)",ip,name,1)
	if err!=nil{
		log.Fatalln(err)
		return err
	} else {
		return nil
	}
}

func UpdateUserStatus(db *sql.DB,ip string,status int) error {
	_,err := db.Exec("update user set status = ? where ip = ?",status,ip)
	if err != nil {
		log.Fatalln(err)
		return err
	} else {
		return nil
	}
}

func UserLogin(db *sql.DB,ip string) error {
	return UpdateUserStatus(db,ip,1)
}

func UserLogout(db *sql.DB,ip string) error {
	return UpdateUserStatus(db,ip,0)
}