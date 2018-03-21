package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Point struct {
	x float64
	y float64
}

type BoundingBox struct {
	TopLeft     Point
	TopRight    Point
	BottomLeft  Point
	BottomRight Point
}

type Tower struct {
	ID          float64 `json:"id"`
	Lat         float64 `json:"est_lat"`
	Lng         float64 `json:"est_lng"`
	Acc         float64 `json:est_acc`
	NetworkName string  `json:"network_name_mapped"`
	PhoneType   string  `json:"phone_type"`
	NetworkID   int     `json:"canonical_network_id"`
	Is2G        bool    `json:"is_2g"`
	Is3G        bool    `json:"test"`
	IsLTE       bool    `json:"is_lte"`
	Confidence  float64 `json:"confidence"`
}

type Store interface {
	GetTower(BoundingBox) ([]*Tower, error)
}

type dbStore struct {
	db *sql.DB
}

var store Store

/* Database */

func InitStore(s Store) {
	store = s
}

func (store *dbStore) GetTower(box BoundingBox) ([]*Tower, error) {
	sql := fmt.Sprintf(
		"SELECT est_lng, est_lat, id, est_acc, network_name_mapped, phone_type, "+
			"canonical_network_id, is_2g, is_3g, is_lte, confidence "+
			"FROM london_towers WHERE ST_Within(geom,ST_GeomFromText("+
			"'POLYGON((%f %f, %f %f, %f %f, %f %f, %f %f))', 4326))",
		box.TopLeft.x, box.TopLeft.y,
		box.TopRight.x, box.TopRight.y,
		box.BottomLeft.x, box.BottomLeft.y,
		box.BottomRight.x, box.BottomRight.y,
		box.TopLeft.x, box.TopLeft.y,
	)

	rows, err := store.db.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	towers := make([]*Tower, 0)
	for rows.Next() {
		t := &Tower{}
		err = rows.Scan(
			&t.Lng, &t.Lat, &t.ID, &t.Acc, &t.NetworkName, &t.PhoneType,
			&t.NetworkID, &t.Is2G, &t.Is3G, &t.IsLTE, &t.Confidence,
		)
		if err != nil {
			return nil, err
		}

		towers = append(towers, t)
	}

	return towers, nil
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Path("/towers").
		HandlerFunc(getTower).
		Queries("lng", "{lng:[+-]?([0-9]*[.])?[0-9]+}", "lat", "{lat:[+-]?([0-9]*[.])?[0-9]+}").
		Methods("GET")

	s := http.StripPrefix("/home/", http.FileServer(http.Dir("./assets/")))
	r.PathPrefix("/home/").Handler(s)

	return r
}

func getTower(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved get tower request")
	vars := mux.Vars(r)

	formlng := vars["lng"]
	formlat := vars["lat"]

	lat, err := strconv.ParseFloat(formlat, 64)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lng, err := strconv.ParseFloat(formlng, 64)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bb := BoundingBox{
		TopLeft:     Point{lat + 0.01, lng - 0.01},
		TopRight:    Point{lat + 0.01, lng + 0.01},
		BottomLeft:  Point{lat - 0.01, lng - 0.01},
		BottomRight: Point{lat - 0.01, lng + 0.01},
	}

	towers, err := store.GetTower(bb)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Got %d results \n", len(towers))

	b, err := json.Marshal(towers)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

func main() {
	dbinfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		"54.166.116.70",
		"otissimon",
		"aQFt8H1uIuP",
		"london_towers",
	)

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	InitStore(&dbStore{db: db})

	r := NewRouter()
	panic(http.ListenAndServe(":8080", r))
}
