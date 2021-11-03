package calculator

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func Handler() {
	http.HandleFunc("/", index)
	http.HandleFunc("/result", result)
	http.ListenAndServe(":8080", nil)
	fmt.Println("こんにちは!")
}

func index(w http.ResponseWriter, rq *http.Request) {
	fr := rq.FormValue("formula")
	cl := rq.FormValue("clear")
	if cl == "true" {
		fr = fr[:(len(fr) - 1)]
	}
	ds := strings.ReplaceAll(fr, "p", "+")
	ds = strings.ReplaceAll(ds, "m", "-")
	ds = strings.ReplaceAll(ds, "t", "×")
	ds = strings.ReplaceAll(ds, "d", "÷")
	item := struct {
		Formula string
		Display string
	}{
		fr,
		ds,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, item)
	if err != nil {
		panic(err)
	}
}

func result(w http.ResponseWriter, rq *http.Request) {
	fr := rq.FormValue("formula")

	// 分割
	var sl1 []string
	for {
		ar := []int{
			strings.Index(fr, "p"),
			strings.Index(fr, "m"),
			strings.Index(fr, "t"),
			strings.Index(fr, "d"),
		}
		nm := ar[0]
		for i := 1; i < len(ar); i++ {
			if ar[i] < nm && ar[i] != -1 {
				nm = ar[i]
			}
		}
		if nm == -1 {
			break
		} else if nm == 0 {
			sl1 = append(sl1, fr[:1])
			fr = fr[1:]
		} else {
			sl1 = append(sl1, fr[:nm])
			fr = fr[nm:]
		}
	}

	// 変換
	var sl2 []string
	var st2 []string
	for i, v := range sl1 {
		if v == "p" || v == "m" || v == "t" || v == "d" {
			st2 = append([]string{v}, st2...)
		} else {
			sl2 = append(sl2, v)
			if st2[0] == "p" || st2[0] == "m" {
				if sl1[i+1] == "t" || sl1[i+1] == "d" {
					continue
				} else {
					sl2 = append(sl2, st2[0])
					st2 = st2[1:]
				}
			} else if st2[0] == "t" || st2[0] == "d" {
				sl2 = append(sl2, st2[0])
				st2 = st2[1:]
			}
		}
	}
	sl2 = append(sl2, st2...)

	// 計算
	var st3 []int
	for _, v := range sl2 {
		if v == "p" {
			st3 = append([]int{(st3[1] + st3[0])}, st3[2:]...)
		} else if v == "m" {
			st3 = append([]int{(st3[1] - st3[0])}, st3[2:]...)
		} else if v == "t" {
			st3 = append([]int{(st3[1] * st3[0])}, st3[2:]...)
		} else if v == "d" {
			st3 = append([]int{(st3[1] % st3[0])}, st3[2:]...)
		} else {
			nm, _ := strconv.Atoi(v)
			st3 = append([]int{nm}, st3...)
		}
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		panic(err)
	}
	item := struct {
		Result int
	}{
		st3[0],
	}
	err = tmpl.Execute(w, item)
	if err != nil {
		panic(err)
	}
}
