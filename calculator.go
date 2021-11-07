package calculator

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func Handler() {
	http.HandleFunc("/", index)
	http.HandleFunc("/result", result)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, rq *http.Request) {
	fr := rq.FormValue("formula")
	cl := rq.FormValue("clear")
	op := ""
	ce := ""
	if cl == "true" {
		fr = fr[:(len(fr) - 1)]
	}
	if len(fr) == 0 {
		ce = "disabled"
	}
	if len(fr) == 0 || fr[len(fr)-1:] == "p" || fr[len(fr)-1:] == "m" || fr[len(fr)-1:] == "t" || fr[len(fr)-1:] == "d" {
		op = "disabled"
	}
	ds := strings.ReplaceAll(fr, "p", "+")
	ds = strings.ReplaceAll(ds, "m", "-")
	ds = strings.ReplaceAll(ds, "t", "×")
	ds = strings.ReplaceAll(ds, "d", "÷")
	item := struct {
		Formula    string
		Display    string
		Operations string
		Cancel     string
	}{
		fr,
		ds,
		op,
		ce,
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
	rs := "Error"

	// 分割
	var sl1 []string
	for {
		ar := []int{
			strings.Index(fr, "p"),
			strings.Index(fr, "m"),
			strings.Index(fr, "t"),
			strings.Index(fr, "d"),
		}
		nm := -1
		for i := 0; i < len(ar); i++ {
			if nm == -1 || ar[i] < nm && ar[i] != -1 {
				nm = ar[i]
			}
		}
		if nm == -1 {
			sl1 = append(sl1, fr)
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
			if len(sl1) < i {
				st2 = append([]string{v}, st2...)
			} else if sl1[i+1] == "p" || sl1[i+1] == "m" || sl1[i+1] == "t" || sl1[i+1] == "d" {
				// 演算子が連続
				sl2 = []string{}
				st2 = []string{}
				break
			} else {
				st2 = append([]string{v}, st2...)
			}
		} else {
			sl2 = append(sl2, v)
			if len(st2) == 0 {
				continue
			} else if st2[0] == "p" || st2[0] == "m" {
				if len(sl1) == i+1 {
					sl2 = append(sl2, st2[0])
					st2 = st2[1:]
				} else if sl1[i+1] == "t" || sl1[i+1] == "d" {
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
	var st3 []float64
	for _, v := range sl2 {
		if v == "p" {
			st3 = append([]float64{(st3[1] + st3[0])}, st3[2:]...)
		} else if v == "m" {
			st3 = append([]float64{(st3[1] - st3[0])}, st3[2:]...)
		} else if v == "t" {
			st3 = append([]float64{(st3[1] * st3[0])}, st3[2:]...)
		} else if v == "d" {
			st3 = append([]float64{(st3[1] / st3[0])}, st3[2:]...)
		} else {
			nm, _ := strconv.ParseFloat(v, 64)
			st3 = append([]float64{nm}, st3...)
		}
	}
	if len(st3) >= 1 {
		rs = strconv.FormatFloat(st3[0], 'f', -1, 64)
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		panic(err)
	}
	item := struct {
		Result string
	}{
		rs,
	}
	err = tmpl.Execute(w, item)
	if err != nil {
		panic(err)
	}
}
