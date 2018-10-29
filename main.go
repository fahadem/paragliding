package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"strconv"
	"github.com/marni/goigc"
)

type Api struct {
	Uptime time.Time `json:"uptime,omitempty"`
	Info string `json:"info,omitempty"`
	Version string `json:"version,omitempty"`
}

type Track struct {
	H_date       string  `json:"H_date,omitempty"` //"H_date": <date from File Header, H-record>,
	Pilot        string  `json:"pilot,omitempty"` //"pilot": <pilot>,
	Glider       string  `json:"glider,omitempty"` //"glider": <glider>,
	Glider_id    string  `json:"glider_id,omitempty"` //"glider_id": <glider_id>,
	Track_length float64 `json:"track_length,omitempty"` //"track_length": <calculated total track length>
	Track_src_url string  `json:"track_src_url,omitempty"`  //<the original URL used to upload the track, ie. the URL used with POST>
}

type Ticker struct {
	T_latest string `json:"t_latest,omitempty"` //"t_latest": <latest added timestamp>,
	T_start string `json:"t_start,omitempty"` //"t_start": <the first timestamp of the added track>, this will be the oldest track recorded
	T_stop string `json:"t_stop,omitempty"` //"t_stop": <the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
	Tracks []string `json:"tracks,omitempty"` //"tracks": [<id1>, <id2>, ...]
	Processing float64 `json:"processing,omitempty"` //"processing": <time in ms of how long it took to process the request>
}

type Webhook struct {
	Webhook_Url string `json:"webhookURL,omitempty"`
	Min_Trigger_Value int64 `json:"minTriggerValue,omitempty"`
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

func (db igcDB) alreadyInDb(fileW File) bool {
	for _, file := range db.igcs {
		if file.Url == fileW.Url {
			return true
		}
	}
	return false
}
func getApi(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if parts[2] == "" {
		http.Redirect(w,r,"https://fahadem2.herokuapp.com/paragliding/api",301)
	}
	http.Header.Add(w.Header(), "content-type", "application/json")
	api := Api{Uptime: time.Now(),
    		 Info: "Service for IGC tracks.",
    		 Version: "v1",
	}
	json.NewEncoder(w).Encode(api)
}

func trackHandler(w http.ResponseWriter, r *http.Request) {

	url:="http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, err := igc.ParseLocation(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	switch r.Method {

	case "POST":
		{
			start := time.Now()
			http.Header.Add(w.Header(), "content-type", "application/json")
			parts := strings.Split(r.URL.Path, "/")
			if r.Body == nil {
				http.Error(w, "no JSON body", http.StatusBadRequest)
				return
			}

			err := json.NewDecoder(r.Body).Decode(&url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			//send message to webhooks
			t = time.Now().Nanosecond()
			e = time.Since(start).Seconds()
			t_conv := strconv.Itoa(t)
			e_conv := fmt.Sprintf("%f", e)
			text := "{\"text\": \"Timestamp :" + t_conv + ", new track is " + parts[5] + " (processing time is " + e + ")\"}"
			payload := strings.NewReader(text)
			for _, wh := range dbWh {
				client := &http.Client{Timeout: (time.Second * 30)}
				req, err := http.NewRequest("POST", wh.Webhook_Url, payload)
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					fmt.Print(err.Error())
				}
				fmt.Println(resp.Status)
			}

		}
	case "GET":
		{
			//GET case

			parts := strings.Split(r.URL.Path, "/")
			
			if strings.HasPrefix(parts[3],"track") && parts[4] == "" {
				//deal with the array
				http.Header.Add(w.Header(), "content-type", "application/json")
				json.NewEncoder(w).Encode(ids)
			}
			if strings.HasPrefix(parts[4],"id") && parts[5] == "" { 
				//deal with the id
				rgx, _ := regexp.Compile("^id[0-9]*")
				id := parts[4]
				ids = append(ids,id)

				if rgx.MatchString(id) == true {
					http.Header.Add(w.Header(), "content-type", "application/json")
					var f File
					db.add(f,id) 
					idCount += 1
					T := Track{}
					T.Glider = track.GliderType
					T.Glider_id = id
					T.Pilot = track.Pilot
					T.Track_length = track.Task.Distance()
					T.H_date = track.Date.String()
					T.Track_src_url = url

					timestamp = append(timestamp,time.Now())
					json.NewEncoder(w).Encode(T)
				}
				if rgx.MatchString(id) == false {
					fmt.Fprintln(w, "Use format id0 or id21 for exemple")
				}
				if strings.HasPrefix(parts[5],"field") {

					fmt.Fprintf(w,"Pilot: %s, gliderType: %s, gliderId: %s,track_length: %f, H_date: %s, track_src_url: %s", track.Pilot, track.GliderType,track.GliderID,track.Task.Distance(), track.Date.String(),url)
				}
			}
		}
	default:
		http.Error(w, "Only GET and POST methods are supported", http.StatusNotImplemented)
	}
}

func latestTicker(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(time.Since(timestamp[len(timestamp)-1]).String())
}

func getApiTicker(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	
	start := time.Now()
	ticker := Ticker{
			T_latest: timestamp[len(timestamp)-1].String(),
			T_start: timestamp[0].String(),
			T_stop: timestamp[len(timestamp)-1].String(),
			Tracks: ids,
			Processing: time.Since(start).Seconds()*1e3,
		  }
	json.NewEncoder(w).Encode(ticker)
}

func webhookNewTrack(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	switch r.Method {
		case "POST":
		{
			webhook := Webhook {}	
			err := json.NewDecoder(r.Body).Decode(&webhook)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			idCountWh += 1
			idWh = parts[5]
			dbWh[idWh] = wh
		}
		case "GET":
		{
			if strings.HasPrefix(parts[5],"id") {
				idWanted := parts[5]
				for id, file := range dbWh {
					if id == idWanted {
						json.NewEncoder(w).Encode(file)
					}
				}
			}
		}
		case "DELETE":
		{
			if strings.HasPrefix(parts[5],"id") {
				idWanted := parts[5]
				for id, file := range dbWh {
					if id == idWanted {
						json.NewEncoder(w).Encode(file)
						delete(dbWh, idWanted)
					}
				}
			}
		}
		default:
			http.NotFound(w, r)

	}	
}

func adminCount(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(db.Count())
}

func adminDel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		nbDel := 0
		ids = []string{}
		for id, _ := range db.igcs {
			delete(db.igcs, id)
			nbDel += 1
		}
		json.NewEncoder(w).Encode(nbDel)
	}
}
var db igcDB
var ids []string
var dbId []string 
var idCount int
var idCountWh int
var idWh string
var dbWh map[string]Webhook
var timestamp []time.Time
var t int
var e float64
//webhookURL:= "https://hooks.slack.com/services/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"

func main() {
	db = igcDB{}
	db.igcs = map[string]File{}
	idCount = 0
	ids = nil
	port := os.Getenv("PORT")
	http.HandleFunc("/paragliding/", getApi)
	http.HandleFunc("/paragliding/api", getApi)
	http.HandleFunc("/paragliding/api/track/", trackHandler)
	http.HandleFunc("/paragliding/api/ticker/latest", latestTicker)
	http.HandleFunc("/paragliding/api/ticker/", getApiTicker)
	http.HandleFunc("/paragliding/api/webhook/new_track/", webhookNewTrack)
	http.HandleFunc("/admin/api/tracks_count", adminCount)
	http.HandleFunc("/admin/api/tracks", adminDel)
	http.ListenAndServe(":"+port, nil)
}
