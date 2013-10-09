package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
	"reflect"
	"runtime"
	"strings"
)

type T struct{}

// prepared statements go here
var (
	Statements = make(map[string]*autorc.Stmt, 0)
)

func GetFunctionName(i interface{}) string {
	arr := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	if strings.Contains(arr, "database.") {
		strArr := strings.Split(arr, "database.")
		return strArr[1]
	}
	return arr
}

func PrepareAll() error {

	//makes a channel
	c := make(chan int)

	go PrepareCurtDev(c)

	<-c

	return nil
}

func GetStatement(key string) (stmt *autorc.Stmt, err error) {
	stmt, ok := Statements[key]
	if !ok {
		qry := expvar.Get(key)
		if qry == nil {
			err = errors.New("Invalid query reference")
		}
	}
	return
}

func PrepareCurtDev(ch chan int) {

	chn := make(chan int, 0)

	// TESTING
	var t T
	val := reflect.ValueOf(&t)
	for i := 0; i < val.NumMethod(); i++ {
		m := val.Method(i)
		if iface := m.Interface(); iface != nil {
			go func(method func(), c chan int) {
				method()
				c <- 1
			}(iface.(func()), chn)
		}
	}

	for i := 0; i < val.NumMethod(); i++ {
		m := val.Method(i)
		if iface := m.Interface(); iface != nil {
			<-chn
			// me := iface.(func())
			log.Printf(" ~ \033[32;1m [OK]\033[0m %s Statements Completed", GetFunctionName(iface))
		}
	}

	ch <- 1
}

func (t *T) WebsiteStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["GetAllContacts"] = "select * from Contact order by created desc"
	UnPreparedStatements["InsertContact"] = "insert into Contact(fname, lname, email, message, created) values(?,?,?,?,?)"
	UnPreparedStatements["GetMessageById"] = "select message from Contact where id=?"

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
	return
}

func PrepareCurtDevStatement(name string, sql string, ch chan int) {
	stmt, err := CurtDevDb.Prepare(sql)
	if err == nil {
		Statements[name] = stmt
	} else {
		log.Println(err)
	}
	ch <- 1
}
