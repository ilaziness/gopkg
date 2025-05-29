package main

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// sqlite作为数据库

var db *sqlx.DB
var err error

func main() {
	log.Default().SetFlags(log.Lshortfile | log.LstdFlags)
	// 连接数据库
	db, err = sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalln("open db err:", err)
	}

	// 测试数据库连通性
	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}

	// 第二种方法， open并ping
	// MustConnect失败会panic，Connect失败会返回error
	//db = sqlx.MustConnect("sqlite3", ":memory:")

	schema := `CREATE TABLE place (
		country text,
		city text NULL,
		telcode integer);`

	// 执行查询
	result, err := db.Exec(schema)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("create table", result)

	//插入数据
	cityState := `INSERT INTO place (country, telcode) VALUES (?, ?)`
	countryCity := `INSERT INTO place (country, city, telcode) VALUES (?, ?, ?)`
	// MustExec出错会panic
	db.MustExec(cityState, "Hong Kong", 852)
	db.MustExec(cityState, "Singapore", 65)
	db.MustExec(countryCity, "South Africa", "Johannesburg", 27)
	// 查询占位符
	// mysql: ?
	// PG: $1 $2
	// sqlite: ? 或者$1语法都可以
	// oracle: :name

	// 获取所有记录
	rows, err := db.Query("SELECT country, city, telcode FROM place")

	// 迭代每一行
	for rows.Next() {
		var country string
		// city可能为NULL, 所以使用 NullString 类型
		var city sql.NullString
		var telcode int
		// scan使用反射将sql类型映射到go类型
		_ = rows.Scan(&country, &city, &telcode)
	}
	// 检查错误
	if err = rows.Err(); err != nil {
		log.Println("get rows err:", err)
	}
	// 释放连接回连接池，如果读取了rows的所有数据不用显示调用，会自动close
	rows.Close()

	// scan到struct
	type Place struct {
		Country       string
		City          sql.NullString
		TelephoneCode int `db:"telcode"`
	}

	rows2, _ := db.Queryx("SELECT * FROM place")
	for rows2.Next() {
		var p Place
		_ = rows2.StructScan(&p)
		log.Println(p)
	}

	// 查询一行
	row := db.QueryRow("SELECT * FROM place WHERE telcode=?", 852)
	var telcode int
	_ = row.Scan(&telcode) //scan完之后释放连接
	log.Println("query row result:", telcode)
	var p Place
	_ = db.QueryRowx("SELECT city, telcode FROM place LIMIT 1").StructScan(&p)
	log.Println("query row result2:", p)

	//get和select
	// 单行get，多行select
	p = Place{}
	pp := []Place{}

	// 获取第一行并scan到p
	err = db.Get(&p, "SELECT * FROM place LIMIT 1")

	// scan到slice pp
	err = db.Select(&pp, "SELECT * FROM place WHERE telcode > ?", 50)

	// 简单类型
	var id int
	err = db.Get(&id, "SELECT count(*) FROM place")

	// 10个name
	var names []string
	err = db.Select(&names, "SELECT name FROM place LIMIT 10")

	// 事务
	// tx, err := db.Begin()
	// err = tx.Exec(...)
	// err = tx.Commit()
	//
	// tx := db.MustBegin()
	// tx.MustExec(...)
	// err = tx.Commit()
	// tx.Rollback()

	// 查询辅助
	// in
	var levels = []int{4, 6, 7}
	query, args, err := sqlx.In("SELECT * FROM users WHERE level IN (?);", levels)
	//query = db.Rebind(query) //Rebind重新绑定到后端
	_, err = db.Query(query, args...)
	if err != nil {
		log.Fatalln("query helper in err:", err)
	}
	//命名查询
	// 查询变量绑定使用命名语法绑定结构体字段或者map的key
	p = Place{Country: "South Africa"}
	_, err = db.NamedQuery(`SELECT * FROM place WHERE country=:country`, p)
	m := map[string]interface{}{"city": "Johannesburg"}
	_, err = db.NamedExec(`SELECT * FROM place WHERE city=:city`, m)
}
