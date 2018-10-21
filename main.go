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


type igcTrack struct {
	H_date       string  //"H_date": <date from File Header, H-record>,
	Pilot        string  //"pilot": <pilot>,
	Glider       string  //"glider": <glider>,
	Glider_id    string  //"glider_id": <glider_id>,
	Track_length float64 //"track_length": <calculated total track length>
}

type Api struct {
	Uptime time.Time `json:"uptime,omitempty"`
	Info string `json:"info,omitempty"`
	Version string `json:"version,omitempty"`

}
type igcFile struct {
	Url string `json:"url,omitempty"`
}

type igcDB struct {
	igcs map[string]igcFile
}

func (db *igcDB) add(igc igcFile, id string) {
	for _, file := range db.igcs {
		if igc == file {
			return
		}
	}
	db.igcs[id] = igc
}

func (db igcDB) Count() int {
	return len(db.igcs)
}

func (db igcDB) Get(idWanted string) igcFile {
	for id, file := range db.igcs {
		if idWanted == id {
			return file
		}
	}
	return igcFile{}
}

func getApi(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	api := Api{Uptime: time.Now(),
    		 Info: "Service for IGC tracks.",
    		 Version: "v1",
	}

	json.NewEncoder(w).Encode(api)

}

func igcHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		{

			if r.Body == nil {
				http.Error(w, "no JSON body", http.StatusBadRequest)
				return
			}
			var igc igcFile
			err := json.NewDecoder(r.Body).Decode(&igc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			
			Idstr := "id"
			strValue := fmt.Sprintf("%d", idCount)
			newId := Idstr + strValue
			ids = append(ids, newId)
			idCount += 1
			db.add(igc, newId)
			json.NewEncoder(w).Encode(igc.Url+newId)
		}
	case "GET":
		{
			//GET case
			http.Header.Add(w.Header(), "content-type", "application/json")
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
				var igcWanted igcFile
				rgx, _ := regexp.Compile("^id[0-9]*")
				id := parts[4]
				if rgx.MatchString(id) == true {
					igcWanted = db.Get(id)

					//then encode the igcFile
					url := igcWanted.Url
					track, err := igc.ParseLocation(url)
					if err != nil {
						//fmt.Errorf("Problem reading the track", err)
					}
					igcT := igcTrack{}
					igcT.Glider = track.GliderType
					igcT.Glider_id = track.GliderID
					igcT.Pilot = track.Pilot
					igcT.Track_length = track.Task.Distance()
					igcT.H_date = track.Date.String()
					json.NewEncoder(w).Encode(igcT)
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

var db igcDB
var ids []string
var idCount int

func main() {
	db = igcDB{}
	db.igcs = map[string]igcFile{}
	idCount = 0
	ids = nil
	port := os.Getenv("PORT")
	http.HandleFunc("/paragliding/api", getApi)
	http.HandleFunc("/paragliding/api/igc/", igcHandler)
	http.ListenAndServe(":"+port, nil)
}
