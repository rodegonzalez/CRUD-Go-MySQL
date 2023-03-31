package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// --------------------------------------------------------------------
// structs
// --------------------------------------------------------------------
type Item struct {
	Id          int
	Name        string
	Description string
}

// --------------------------------------------------------------------
// vars and configuration
// --------------------------------------------------------------------

// templates location
var templates = template.Must(template.ParseGlob("views/*"))

// server port
var serverport string = ":8080"

// logging
var logging bool = true

// database
var dbdriver string = "mysql"
var dbuser string = "test"
var dbpass string = "test"
var dbname string = "test"

//var dbhost string = "127.0.0.1"

// --------------------------------------------------------------------
// logger
// --------------------------------------------------------------------
func logger(msg any) {
	if logging {
		log.Println(msg)
	}
}

// --------------------------------------------------------------------
// Main
// --------------------------------------------------------------------
func main() {
	//fmt.Println("Hello World!")

	// common files like css, js, images, icons
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// routing
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/list", ListHandler)
	http.HandleFunc("/create", CreateHandler)
	http.HandleFunc("/insert", InsertHandler)
	http.HandleFunc("/delete", DeleteHandler)
	http.HandleFunc("/edit", EditHandler)
	http.HandleFunc("/update", UpdateHandler)

	// start the server
	logger("Server running on port " + serverport + ".")
	http.ListenAndServe(serverport, nil)
}

// --------------------------------------------------------------------
// views and handlers
// --------------------------------------------------------------------
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(w, "home", nil)
}

func ListHandler(w http.ResponseWriter, req *http.Request) {

	connDB := connDB()
	recordset, err := connDB.Query("SELECT * FROM items")

	if err != nil {
		panic(err.Error())
	}

	item := Item{}
	itemArray := []Item{}

	for recordset.Next() {
		var id int
		var name, description string
		err = recordset.Scan(&id, &name, &description)

		if err != nil {
			panic(err.Error())
		}

		item.Id = id
		item.Name = name
		item.Description = description

		itemArray = append(itemArray, item)
	}

	logger("Item list:")
	logger(itemArray)
	templates.ExecuteTemplate(w, "list", itemArray)
}

func CreateHandler(w http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(w, "create", nil)
}

func InsertHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		name := req.FormValue("name")
		description := req.FormValue("description")
		logger("Insert item:")
		logger("name: " + name)
		logger("description: " + description)

		connDB := connDB()
		insertRecords, err := connDB.Prepare("INSERT INTO items (name, description) VALUES (?,?)")

		if err != nil {
			panic(err.Error())
		}

		insertRecords.Exec(name, description)
		http.Redirect(w, req, "/list", http.StatusMovedPermanently)
	}
}

func DeleteHandler(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")
	logger("Deleting id: " + id)

	connDB := connDB()
	deleteRecords, err := connDB.Prepare("DELETE FROM items WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	deleteRecords.Exec(id)
	http.Redirect(w, req, "/list", http.StatusMovedPermanently)
}

func EditHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	logger("Edit id: " + id)

	connDB := connDB()
	recordset, err := connDB.Query("SELECT * FROM items WHERE id=?", id)

	if err != nil {
		panic(err.Error())
	}

	item := Item{}
	for recordset.Next() {
		var id int
		var name, description string
		err = recordset.Scan(&id, &name, &description)

		if err != nil {
			panic(err.Error())
		}
		item.Id = id
		item.Name = name
		item.Description = description
	}
	logger("Editing item. Values:")
	logger(item)
	templates.ExecuteTemplate(w, "edit", item)
}

func UpdateHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		id := req.FormValue("id")
		name := req.FormValue("name")
		description := req.FormValue("description")
		logger("Updating:")
		logger("Update id: " + id)
		logger("Update name: " + name)
		logger("Update description: " + description)

		connDB := connDB()
		updateRecords, err := connDB.Prepare("UPDATE items SET name=?, description=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		updateRecords.Exec(name, description, id)
		http.Redirect(w, req, "/list", http.StatusMovedPermanently)
	}
}

// --------------------------------------------------------------------
// database
// --------------------------------------------------------------------
func connDB() (connDB *sql.DB) {

	//var connstring string = dbuser + ":" + dbpass + "@" + dbhost + "/" + dbname
	var connstring string = dbuser + ":" + dbpass + "@/" + dbname
	connDB, err := sql.Open(dbdriver, connstring)

	if err != nil {
		//log.Fatal("Fallo al conectar a la base de datos." + err.Error())
		panic(err.Error())
	}
	return connDB
}
