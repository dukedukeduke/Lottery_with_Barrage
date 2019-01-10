package main

import (
	"./settings/jsonParse"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const MaxLengthBubble int = 100
const WebsocketMax = 10
const GetWebsocketFailed = -1
const QueueMax = 100
const MessageMax = 100
const HeartbeatGap = 60

var staticPath string
var templatePath string
var Employees settings.EmployeeList
var msgPool MessagePool
var pool WebsocketPool
var Winner = make(map[string]bool)
var Exclude = make(map[string]bool)

type MessagePool struct{
	sync.Mutex
	Message []string
}

type MessagePoolFullError string
type MessageLengthError string
type MessagePoolEmptyError string
type ConnNotFound string
type WsNotFound string

type WebsocketUtil struct {
	connId string
	conn *websocket.Conn
	status bool
	activeTime time.Time
}

type WebsocketPool struct {
	sync.Mutex
	util []WebsocketUtil
	poolId uuid.UUID
}

type AwardOutput struct{
	Employee settings.EmployeeList
	Winner []string}

type SaveLotteryOutput struct{
	Status string
	Winner string}

func (e MessagePoolFullError) Error() string   { return "message pool full for message: " + string(e) }

func (e MessagePoolEmptyError) Error() string   { return "message pool empty for pop, length is:" + string(e)}

func (e ConnNotFound) Error() string   { return "Conn not found, id is:" + string(e)}

func (e WsNotFound) Error() string   { return string(e)}

func (msg *MessagePool) Push(message string) (err error){
	var(
		lenTmp int
	)

	msg.Lock()
	defer msg.Unlock()
	lenTmp = len(msg.Message)
	if lenTmp < QueueMax{
		msg.Message = append(msg.Message, message)
		return nil
	}else{
		return MessagePoolFullError(message)
	}
}

func (msg *MessagePool) Pop() (string, error){
	var(
		lenTmp int
		msgTmp string
	)
	msg.Lock()
	defer msg.Unlock()
	lenTmp = len(msg.Message)
	if lenTmp > 0{
		msgTmp = msg.Message[lenTmp - 1]
		msg.Message = msg.Message[0: lenTmp - 1]
		return msgTmp, nil
	}else{
		return "", MessagePoolEmptyError(lenTmp)
	}
}

func (cr *WebsocketPool)GetWebsocketConn(ws *websocket.Conn) (string, error){
	var (
		IdTmp uuid.UUID
		errUuid error
	)
	cr.Lock()
	activeTime := time.Now()
	defer cr.Unlock()
	if len(cr.util) >= WebsocketMax{
		log.Println("Reach the max and will pop the previous conn")
		cr.util = cr.util[1:]
		if IdTmp , errUuid = uuid.NewV4();errUuid != nil{
			log.Println("Generate uuid failed")
			return "", errUuid
		}
		cr.util = append(cr.util, WebsocketUtil{IdTmp.String(), ws, true, activeTime})
	}else{
		if IdTmp , errUuid = uuid.NewV4();errUuid != nil{
			log.Println("Generate uuid failed")
			return "", errUuid
		}
		cr.util = append(cr.util, WebsocketUtil{IdTmp.String(), ws, true, activeTime})
	}
	return IdTmp.String(), nil
}

func (cr *WebsocketPool)FindWebsocketByConnId(id string) (*WebsocketUtil, error){
	cr.Lock()
	defer cr.Unlock()

	var tmp int = -1
	for index, value :=range cr.util{
		if value.connId == id{
			tmp = index
		}
	}
	if tmp > -1{
		return &cr.util[tmp], nil
	}else{
		return nil, WsNotFound("")
	}

}

func (cr *WebsocketPool)FindWebsocketByConn(ws *websocket.Conn) (conn *WebsocketUtil, err error){
	cr.Lock()
	defer cr.Unlock()

	var tmp int = -1
	for index, value :=range cr.util{
		if value.conn == ws{
			tmp = index
		}
	}
	if tmp != -1{
		return &cr.util[tmp], nil
	}else{
		return nil, WsNotFound("Ws not found")
	}

}

func (cr *WebsocketPool)ReleaseWebsocketConn(id string) error{
	cr.Lock()
	defer cr.Unlock()

	var tmp int = -1
	for index, value :=range cr.util{
		if value.connId == id{
			tmp = index
		}
	}
	if tmp > -1{
		tmpUtil := cr.util[tmp+1:]
		cr.util = cr.util[:tmp]
		for i,_ := range tmpUtil[:]{
			cr.util = append(cr.util, tmpUtil[i])
		}
		return nil
	}else{
		return ConnNotFound(id)
	}
}

func (cr *WebsocketPool)DeletePool(){
	pool = WebsocketPool{}
}

func init(){
	var Config settings.Config
	var ConfigData settings.JsonParse = &Config

	var Employee settings.EmployeeList
	var EmployeeData settings.JsonParse = &Employee

	// parse config.json
	errConfig := ConfigData.Load("settings/config.json")
	if errConfig != nil {
		panic("Parse config json failed")
	}
	staticPath = Config.Static
	templatePath = Config.Template

	// parse employee.json
	employeeConfig := EmployeeData.Load("static/employee.json")
	if employeeConfig != nil {
		panic("Parse config json failed")
	}
	Employees = Employee

	go func() {
		for{
			var(
				data string
				errGetData error
			)
		MainStart:
			for index, _ws := range pool.util {
				fmt.Println(time.Now().Sub(pool.util[index].activeTime))
				fmt.Println(60*time.Second)
				if time.Now().Sub(pool.util[index].activeTime) > HeartbeatGap*time.Second{
					log.Println("Timeout point")
					log.Println("Timeout for client:", _ws.connId)
					pool.ReleaseWebsocketConn(_ws.connId)
					fmt.Println(pool.util)
				}
			}
			bubbleLength := len(msgPool.Message)
			if bubbleLength < 1{
				time.Sleep(10*time.Second)
				goto MainStart
			}
			if data, errGetData = msgPool.Pop();errGetData != nil{
				log.Println(errGetData)
				goto MainStart
			}
			for _, _ws := range pool.util {
				if time.Now().Sub(_ws.activeTime) < HeartbeatGap*time.Second{
					if err := websocket.Message.Send(_ws.conn, string(data)); err != nil {
						log.Println(err)
						log.Println("Connect error for client:", _ws.connId)
						pool.ReleaseWebsocketConn(_ws.connId)
					}
				}else{
					log.Println("Timeout for client:", _ws.connId)
					pool.ReleaseWebsocketConn(_ws.connId)
				}
			}
		}
	}()

}

func award(w http.ResponseWriter, r *http.Request) {
	var(
		t *template.Template
		err error
	)
	if r.Method == "GET" {
		var winnerList []string
		for key, value := range Winner{
			if value{
				winnerList = append(winnerList, key)
			}
		}

		if t, err = template.ParseFiles(templatePath + "/award.gtpl");err !=nil{
			log.Print(err)
			return
		}
		t.Execute(w, AwardOutput{Employees, winnerList})
	}
}

func message(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles(templatePath + "/message.gtpl")
		t.Execute(w, nil)
	}else if r.Method == "POST"{
        tmpMessage :=r.Form["message"][0]
        if len(tmpMessage) > MessageMax{
			fmt.Println("Message length is not less than " + string(MessageMax))
			t, _ := template.ParseFiles(templatePath + "/message.gtpl")
			t.Execute(w, nil)
			return
		}
		msgPushError := msgPool.Push(tmpMessage)
		if msgPushError != nil{
			fmt.Println(msgPushError)
		}
		t, _ := template.ParseFiles(templatePath + "/message.gtpl")
		t.Execute(w, nil)
	}else{
		fmt.Println("method is:" + r.Method)
		fmt.Fprintf(w, "Method illegal")
		return
	}
}


func bubble(ws *websocket.Conn) {
	var (
		newId string
		NewConnError error
	)

	if newId, NewConnError = pool.GetWebsocketConn(ws); NewConnError != nil{
		if errSend := websocket.Message.Send(ws, "status:connect failed, please try again later");errSend != nil{
			log.Println(errSend)
			log.Println("Client connect error")
			return
		}
	}else{
		log.Println("Create new conn successfully, id is:", newId)
	}

	log.Println(pool.util)

	defer func() {
		log.Println("Error happend and disconnect")
		pool.ReleaseWebsocketConn(newId)
	}()

	for {
		START:
		var (
			reply string
			wsWebsocketUtil *WebsocketUtil
			errWs error
			)
		if errRecv := websocket.Message.Receive(ws, &reply);errRecv != nil{
			log.Println(errRecv)
			log.Println("Client disconnect for conn:", newId)
			return
		}

		log.Println("Heartbeat from:", newId)

		if wsWebsocketUtil, errWs = pool.FindWebsocketByConnId(newId);errWs != nil{
			log.Println(errWs)
			return
		}
		wsWebsocketUtil.activeTime = time.Now()

		time.Sleep(10*time.Second)
		goto START
	}
}


func SetTimer(){
	time.Sleep(3* time.Second)
}

func saveLottery(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dataList []string
	t := template.New("saveLottery")
	if r.Method == "GET" {
		tmpList, ok := r.URL.Query()["list"]
		if !ok || len(tmpList) < 1{
			t.Execute(w, SaveLotteryOutput{"failed", ""})
		}else{
			dataList = strings.Split(tmpList[0], ",")
		}
		if len(dataList) > 0{
			for _, value :=range dataList{
				if value != ""{
					fmt.Println(value)
					Winner[value] = true
				}
			}
			var winnerList []string
			var winnerString string
			for key, value := range Winner{
				if value{
					winnerList = append(winnerList, key)
				}
			}
			winnerString = strings.Join(winnerList, ",")
			fmt.Println(SaveLotteryOutput{"success", winnerString})
			t.Parse("{{.Status}}:" + "{{.Winner}}")
			t.Execute(w, SaveLotteryOutput{"success", winnerString})
		}else{
			//fmt.Fprintf(w, "failed")
			t.Execute(w, SaveLotteryOutput{"failed", ""})
		}
	}else{
		fmt.Println("method is:" + r.Method)
		//fmt.Fprintf(w, "Method illegal")
		t.Execute(w, SaveLotteryOutput{"failed", ""})
	}
}


func main() {
	var errUUID error
	pool = WebsocketPool{}
	pool.poolId , errUUID= uuid.NewV4()
	if errUUID != nil{
		panic("Apply Pool failed")
		return
	}

	// static files url
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir(templatePath))))

	// award page, just show the page
	http.HandleFunc("/lottery/", award)

	// award page, just show the page
	http.HandleFunc("/save_lottery/", saveLottery)

	// login page
	http.HandleFunc("/message/", message)

	// get the bubble text and send to award page
	http.Handle("/bubble", websocket.Handler(bubble))
	fmt.Print("Start Server...")
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
