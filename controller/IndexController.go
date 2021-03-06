package controller

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"go_movies/services"
	heroTpl "go_movies/views/hero"
	"net/http"
	"strconv"
)

// 首页
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	path := r.URL.Path
	cate := r.URL.Query()["cate"]
	_start := r.URL.Query()["start"]
	_stop := r.URL.Query()["stop"]

	// 需要展示的数据
	show := make(map[string]interface{})

	// 所有类别/导航
	Categories := services.AllCategoryDate()

	key := "detail_links:hk" // 默认首页

	start := int64(0)
	stop := int64(14)

	if len(_start) > 0 {
		StartInt64, _ := strconv.ParseInt(_start[0], 10, 64)
		start = StartInt64
	}

	if len(_stop) > 0 {
		StopInt64, _ := strconv.ParseInt(_stop[0], 10, 64)
		stop = StopInt64
	}

	prev := path + "?start=" + strconv.FormatInt(start-15, 10) + "&stop=" + strconv.FormatInt(stop-15, 10)
	next := path + "?start=" + strconv.FormatInt(start+15, 10) + "&stop=" + strconv.FormatInt(stop+15, 10)

	prevStatus := "1"
	nextStatus := "1"

	navLink := "/"
	if len(cate) > 0 {
		key = "detail_links:id:" + services.TransformCategoryId(cate[0])
		navLink = cate[0]
		prev = path + "?cate=" + cate[0] + "&start=" + strconv.FormatInt(start-15, 10) + "&stop=" + strconv.FormatInt(stop-15, 10)
		next = path + "?cate=" + cate[0] + "&start=" + strconv.FormatInt(start+15, 10) + "&stop=" + strconv.FormatInt(stop+15, 10)
	}

	if start > stop || stop-start > 15 || start < 0 {
		start = 0
		stop = 15
	}

	MovieLists := services.MovieListsRange(key, start, stop)

	LenMovieLists := len(MovieLists)

	if start-15 < 0 {
		prevStatus = "0"
	}

	if LenMovieLists < 15 || LenMovieLists == 0 {
		nextStatus = "0"
	}

	NewFilmKey := "detail_links:id:1"
	NewTVKey := "detail_links:id:2"

	NewFilm := services.MovieListsRange(NewFilmKey, 0, 49)
	NewTV := services.MovieListsRange(NewTVKey, 0, 49)

	show["categories"] = Categories
	show["page"] = "页面"
	show["movieLists"] = MovieLists
	show["new_film"] = NewFilm
	show["new_tv"] = NewTV
	show["prev"] = prev
	show["next"] = next
	show["prev_status"] = prevStatus
	show["next_status"] = nextStatus
	show["nav_link"] = navLink

	buffer := new(bytes.Buffer)
	heroTpl.Index(show, buffer)
	w.Write(buffer.Bytes())

}

func Movie(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	link := r.URL.Query()["link"][0]

	// 需要展示的数据
	show := make(map[string]interface{})

	// 所有类别/导航
	Categories := services.AllCategoryDate()

	MovieDetail := services.MovieDetail(link)

	NewFilmKey := "detail_links:id:1"
	NewTVKey := "detail_links:id:2"

	NewFilm := services.MovieListsRange(NewFilmKey, 0, 49)
	NewTV := services.MovieListsRange(NewTVKey, 0, 49)

	show["categories"] = Categories
	show["MovieDetail"] = MovieDetail
	show["new_film"] = NewFilm
	show["new_tv"] = NewTV
	show["nav_link"] = "/"

	buffer := new(bytes.Buffer)
	heroTpl.MDetail(show, buffer)
	w.Write(buffer.Bytes())
}

func Play(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	show := make(map[string]interface{})

	PlayUrl := r.URL.Query()["play_url"][0]
	PlayType := r.URL.Query()["play_type"][0]

	show["play_url"] = PlayUrl
	show["type"] = PlayType

	buffer := new(bytes.Buffer)

	if PlayType == "mp4" {
		heroTpl.Mp4(show, buffer)
	} else if PlayType == "m3u8" {
		heroTpl.M3u8(show, buffer)
	} else {
		heroTpl.Kuyun(show, buffer)
	}

	w.Write(buffer.Bytes())
}

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	show := make(map[string]interface{})

	q := ""
	_q := r.URL.Query()["q"]
	if len(_q) > 0 {
		q = _q[0]
	}

	if q == "" {
		q = "[]"
	}

	MovieLists := services.SearchMovies(q)

	// 所有类别/导航
	Categories := services.AllCategoryDate()

	NewFilmKey := "detail_links:id:1"
	NewTVKey := "detail_links:id:2"

	NewFilm := services.MovieListsRange(NewFilmKey, 0, 49)
	NewTV := services.MovieListsRange(NewTVKey, 0, 49)

	show["movieLists"] = MovieLists
	show["categories"] = Categories
	show["new_film"] = NewFilm
	show["new_tv"] = NewTV
	show["nav_link"] = "/"

	buffer := new(bytes.Buffer)
	heroTpl.Search(show, buffer)
	w.Write(buffer.Bytes())
}

func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// 需要展示的数据
	show := make(map[string]interface{})

	// 所有类别/导航
	Categories := services.AllCategoryDate()

	show["categories"] = Categories

	show["nav_link"] = "/about"

	buffer := new(bytes.Buffer)
	heroTpl.About(show, buffer)
	w.Write(buffer.Bytes())
}
