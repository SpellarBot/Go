package utils
import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)
type MysqlClient struct {
	DatabaseName string
	UserName string
	PassWrod string
	DB *sql.DB
}
func(M *MysqlClient) Init(database string, username string, password string){
	M.PassWrod = password
	M.DatabaseName = username
	M.DatabaseName = database
	S := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8",
		username,password,database)
	DB,err := sql.Open("mysql",S)
	if err==nil {
		//fmt.Println("Connect to Mysql success")
		M.DB = DB
	}
}

func (M MysqlClient) Exec(Sql string){
	M.DB.Exec(Sql)
}


func MysqlTest(){
	Client := MysqlClient{}
	Client.Init("test","pig","123456")
	db := Client.DB
	stmt, err := db.Prepare(`INSERT into test values (?,?)`)

	if err==nil{
		fmt.Println("begin to insert")
		stmt.Exec("wangxiao",19)
		fmt.Println("succ to insert")
	}
	S,err := db.Query("select * from test")
	for S.Next(){
		var name string
		var age int
		S.Scan(&name,&age)
		fmt.Println(name,age)
	}

	//db.Exec("DELETE  from test")
	N,err := db.Query(`select count(*) from test`)
	var NN int
	for N.Next(){
		N.Scan(&NN)
	}
	fmt.Println(NN)

	F := db.QueryRow("select count(*) from test where name =?","wangxiao")
	err = F.Scan(&NN)
	fmt.Println(err,NN)

	db.Close()


}

