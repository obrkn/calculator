package calculator

import (
	"fmt"
	"net/http"
	"text/template"
)

func Handler() {
	http.HandleFunc("/", index)
	http.HandleFunc("/result", result)
	http.ListenAndServe(":8080", nil)
	fmt.Println("こんにちは!")
}

func index(w http.ResponseWriter, rq *http.Request) {
	// fmt.Fprintln(w, "こんにちは")
	type Inventory struct {
		Formula string
		Count   uint
	}
	sweaters := Inventory{"wool", 17}
	fmt.Println(rq.FormValue("formula"))
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, sweaters)
	if err != nil {
		panic(err)
	}
}

func result(w http.ResponseWriter, rq *http.Request) {

}
