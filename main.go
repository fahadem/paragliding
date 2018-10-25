package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/marni/goigc"
)


type Track struct {
	H_date       string  `json:"H_date,omitempty"` //"H_date": <date from File Header, H-record>,
	Pilot        string  `json:"pilot,omitempty"` //"pilot": <pilot>,
	Glider       string  `json:"glider,omitempty"` //"glider": <glider>,
	Glider_id    string  `json:"glider_id,omitempty"` //"glider_id": <glider_id>,
	Track_length float64 `json:"track_length,omitempty"` //"track_length": <calculated total track length>
	Track_src_url string  `json:"track_src_url,omitempty"`  
}

type Ticker struct {
	T_latest string `json:"t_latest,omitempty"` //"t_latest": <latest added timestamp>,
	T_start string `json:"t_start,omitempty"` //"t_start": <the first timestamp of the added track>, this will be the oldest track recorded
	T_stop string `json:"t_stop,omitempty"` //"t_stop": <the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
	Tracks []Track `json:"tracks,omitempty"` //"tracks": [<id1>, <id2>, ...]
	Processing time.Time `json:"processing,omitempty"` //"processing": <time in ms of how long it took to process the request>
}

type Api struct {
	Uptime time.Time `json:"uptime,omitempty"`
	Info string `json:"info,omitempty"`
	Version string `json:"version,omitempty"`

}
type File struct {
	Url string `json:"url,omitempty"`
}

type igcDB struct {
	igcs map[string]File
}

func (db *igcDB) add(file File, id string) {
	for _, f := range db.igcs {
		if file == f {
			return
		}
	}
	db.igcs[id] = file
}

func (db igcDB) Count() int {
	return len(db.igcs)
}

func (db igcDB) Get(idWanted string) File {
	for id, file := range db.igcs {
		if idWanted == id {
			return file
		}
	}
	return File{}
}

func getApi(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	api := Api{Uptime: time.Now(),
    		 Info: "Service for IGC tracks.",
    		 Version: "v1",
	}

	json.NewEncoder(w).Encode(api)

}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	switch r.Method {
	case "POST":
		{

			if r.Body == nil {
				http.Error(w, "no JSON body", http.StatusBadRequest)
				return
			}
			var file File
			err := json.NewDecoder(r.Body).Decode(&file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			//json.NewEncoder(w).Encode(url)
			Idstr := "id"
			strValue := fmt.Sprintf("%d", idCount)
			newId := Idstr + strValue
			ids = append(ids, newId)
			idCount += 1
			db.add(file, newId)
			json.NewEncoder(w).Encode(newId)
		}
	case "GET":
		{
			//GET case

			parts := strings.Split(r.URL.Path, "/")
			
			if len(parts) < 5 || len(parts) > 6 {
				//deal with errors
				json.NewEncoder(w).Encode("404")
				return
			}
			if parts[4] == "" {
				//deal with the array
				json.NewEncoder(w).Encode(ids)

			}
			if parts[4] != "" {
				//deal with the id
				//var Wanted File
				rgx, _ := regexp.Compile("^id[0-9]*")
				id := parts[4]
				if rgx.MatchString(id) == true {
					//Wanted = db.Get(id)

					//encode the File
					//url := Wanted.Url
					s:="http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
					track, err := igc.ParseLocation(s)
					if err != nil {
						//fmt.Errorf("Problem reading the track", err)
					}
					T := Track{}
					T.Glider = track.GliderType
					T.Glider_id = track.GliderID
					T.Pilot = track.Pilot
					T.Track_length = track.Task.Distance()
					T.H_date = track.Date.String()

					latestT = time.Now()
					json.NewEncoder(w).Encode(T)


				}
				if rgx.MatchString(id) == false {
					fmt.Fprintln(w, "Use format id0 or id21 for exemple")
				}
			}

		}
	default:

		http.Error(w, "Only GET and POST methods are supported", http.StatusNotImplemented)

	}
}

func latestTicker(w http.ResponseWriter, r *http.Request) {
	//http.Header.Add(w.Header(), "content-type", "application/json")
	//parts := strings.Split(r.URL.Path, "/")
	json.NewEncoder(w).Encode(time.Since(latestT).String())
}

var db igcDB
var ids []string
var idCount int
var latestT time.Time

func main() {
	db = igcDB{}
	db.igcs = map[string]File{}
	idCount = 0
	ids = nil
	port := os.Getenv("PORT")
	http.HandleFunc("/paragliding/api", getApi)
	http.HandleFunc("/paragliding/api/track/", trackHandler)
	http.HandleFunc("/paragliding/api/ticker/latest", latestTicker)
	http.ListenAndServe(":"+port, nil)
}
