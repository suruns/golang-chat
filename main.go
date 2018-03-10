package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"app/dao"
	"app/model"
	"encoding/json"
	"time"
	"app/common"
	"path"
	"strconv"
	"os"
)

const (
	TEXT        = 1
	HEARTBEAT   = 2
	LOGIN       = 3
	LOGOUT      = 4
	IMAGE       = 5
)

type Client struct {
	User         model.User
	WsConn       *websocket.Conn
	LastTime     int64
}
type Message struct {
	User         model.User           `json:"user"`
	Message      string			      `json:"message"`
	Type         int                  `json:"type"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:1024,
	WriteBufferSize:1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make([]Client,0)

func wshandler(w http.ResponseWriter, r *http.Request,db *sql.DB)  {
	conn, err := wsupgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	var message Message
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		json.Unmarshal(msg,&message)

		if message.Type == TEXT || message.Type == LOGIN{
			hasIp := false
			for _,v := range clients{
				if v.User.Ip == message.User.Ip {
					hasIp = true
				}
			}

			if !hasIp {
				client := Client{message.User,conn,common.GetCurrentTime()}
				clients = append(clients,client)
			}
			err := dao.UserLogin(db,message.User.Ip)
			if err != nil{
				log.Fatalln(err)
				panic(err)
			}
			SendAllMessage(t,msg)
		} else if message.Type == HEARTBEAT {
			for i,v := range clients {
				if v.User.Ip == message.User.Ip {
					clients[i].LastTime = common.GetCurrentTime()
				}
			}
		}
	}
}
func websocketHeartbeat(db *sql.DB)  {
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for _ = range ticker.C {
			for i,v := range clients {
				err := v.WsConn.WriteMessage(websocket.PingMessage,nil)
				if err!=nil{
					logoutMessage := Message{v.User,"",LOGOUT}
					l,err := json.Marshal(logoutMessage)
					if err != nil{
						log.Fatalln(err)
						panic(err)
					}
					err = dao.UserLogout(db,v.User.Ip)
					if err != nil{
						log.Fatalln(err)
						panic(err)
					}
					clients = append(clients[:i],clients[i+1:]...)
					SendAllMessage(websocket.TextMessage,l)
				}
			}
		}
	}()
}
func SendAllMessage(messageType int, data []byte)  {
	for _,k := range clients{
		k.WsConn.WriteMessage(messageType,data)
	}
}
func main() {
	pathName := "f:/data/go"
	db,err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/chat?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil{
		log.Fatalln(err)
	}

	router := gin.Default()
	router.LoadHTMLGlob("view/*")
	router.Static("/static","./static")
	router.Static("/upload",pathName)
	router.GET("/list", func(c *gin.Context) {
		ip := c.Query("ip")
		users, err := dao.GetUserList(db,ip)
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusOK,gin.H{
				"status":404,
			})
		}else {
			c.JSON(http.StatusOK,gin.H{
				"status":200,
				"entities":users,
			})
		}
	})
	//router.GET("/index", func(c *gin.Context) {
	//	c.HTML(http.StatusOK,"index.html",gin.H{
	//		"users":users,
	//	})
	//})
	websocketHeartbeat(db)
	router.GET("/checkIp", func(c *gin.Context) {
		ip := c.Query("ip")
		user,isExited := dao.GetUserByIp(db,ip)
		if isExited {
			c.JSON(http.StatusOK,gin.H{
				"status":200,
				"entity":user,
			})
		}else {
			c.JSON(http.StatusOK,gin.H{
				"status":404,
				"msg":"ip is not exist",
			})
		}
	})
	router.POST("/register", func(c *gin.Context) {
		name := c.PostForm("name")
		ip := c.PostForm("ip")
		log.Println(name,">>>>>>>>>",ip,">>>>>>>",c.Query("name"))
		_,ok := dao.GetUserByName(db,name);
		if ok {
			c.JSON(http.StatusOK,gin.H{
				"status":400,
				"msg":"name is exist",
			})
		} else {
			err := dao.UpdateUserByIp(db,ip,name)
			if err != nil {
				c.JSON(http.StatusOK,gin.H{
					"status":400,
					"msg":err.Error(),
				})
			} else{
				c.JSON(http.StatusOK,gin.H{
					"status":200,
					"msg":"register successful",
				})
			}
		}
	})
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		fileType := c.PostForm("type")
		ip := c.PostForm("ip")
		if err != nil {
			c.JSON(http.StatusOK,gin.H{
				"status":400,
				"msg":"get file error:" + err.Error(),
			})
			return
		}
		os.MkdirAll(pathName,0777)
		fileName := strconv.FormatInt(time.Now().Unix(),10)   + path.Ext(file.Filename)
		filePathName := pathName + "/" + fileName
		if err:=c.SaveUploadedFile(file, filePathName);err != nil{
			c.JSON(http.StatusOK,gin.H{
				"status":400,
				"msg":"save file error:"+err.Error(),
			})
		}else{
			scheme := "http://"
			if c.Request.TLS != nil{
				scheme = "https://"
			}
			user,_ := dao.GetUserByIp(db,ip)
			for _, k := range clients {
				msg := Message{
					User: *user,
					Message:scheme + c.Request.Host+"/upload/"+fileName,
					Type:IMAGE,
				}
				text,err := json.Marshal(msg);
				if err!=nil{
					log.Fatalln(err)
				}
				k.WsConn.WriteMessage(websocket.TextMessage,text)
			}
			c.JSON(http.StatusOK,gin.H{
				"status":200,
				"msg":"upload file successful",
				"url":scheme + c.Request.Host+"/upload/"+fileName,
				"type":fileType,
			})
		}
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK,"index.html",c.ClientIP())
	})
	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK,"test.html",nil)
	})
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer,c.Request,db)
	})
	router.Run(":8000")
}

