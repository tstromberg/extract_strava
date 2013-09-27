// Simple tool to extract Strava rides from the Strava Android database.
//
// Usage:
//
// extract_strava </path/to/strava/database> <output directory>
//
// For more information, please see README.md.

package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path"
	"text/template"
	"time"
)

var (
	// cache the loaded TCX template
	tcxTemplate = loadTemplate("tcx.tmpl")
)

// Waypoint maps to the "waypoints" table in the strava database.
type Waypoint struct {
	RideId    string `db:"ride_id"`
	Pos       int    `db:"pos"`
	Timestamp int    `db:"timestamp"`
	// Yes, this is misspelled in the db
	Latitude    float64 `db:"latiude"`
	Longitude   float64 `db:"longitude"`
	Altitude    float64 `db:"altitude"`
	HAccuracy   float64 `db:"h_accuracy"`
	VAccuracy   float64 `db:"v_accuracy"`
	Command     string  `db:"command"`
	Speed       float64 `db:"speed"`
	Bearing     float64 `db:"bearing"`
	DeviceTime  int     `db:"device_time"`
	Filtered    int     `db:"filtered"`
	ElapsedTime int     `db:"elapsed_time"`
	Distance    float64 `db:"distance"`
}

// Waypoint.ExportFilename returns a usable Export filename for a waypoint.
func (w *Waypoint) ExportFilename() string {
	return fmt.Sprintf("%s-%s.tcx", w.TimeString(), w.RideId)
}

// Waypoint.TimeString returns the timestamp in RFC 3339 form.
func (w *Waypoint) TimeString() string {
	t := time.Unix(int64(w.Timestamp/1000), 0)
	return t.Format(time.RFC3339)
}

// TemplateContext has the fields to be sent to the TCX template.
type TemplateContext struct {
	StartTime string
	Waypoints []Waypoint
}

// loadTemplate reads and configures a set of templates
func loadTemplate(path string) *template.Template {
	t := template.New(path)
	_, err := t.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	return t
}

// readWaypoints reads a Strava database and returns Waypoint structs.
func readWaypoints(path string) (waypoints []Waypoint, err error) {
	db, err := sqlx.Open("sqlite3", os.Args[1])
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Queryx("select * FROM waypoints")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	waypoint := Waypoint{}
	for rows.Next() {
		rows.StructScan(&waypoint)
		waypoints = append(waypoints, waypoint)
	}
	rows.Close()
	return
}

// saveRide saves a series of waypoints representing a single ride to a TCX file at outPath.
func saveRide(waypoints []Waypoint, outPath string) error {
	log.Printf("Saving %d waypoints to %s", len(waypoints), outPath)
	context := &TemplateContext{
		StartTime: waypoints[0].TimeString(),
		Waypoints: waypoints,
	}

	f, err := os.Create(outPath)
	defer f.Close()
	if err != nil {
		return err
	}
	return tcxTemplate.ExecuteTemplate(f, "tcx.tmpl", context)
}

// exportWaypoints takes a series of waypoints, splits them, then dumps the resulting TCX files into a directory.
func exportWaypoints(waypoints []Waypoint, outDirectory string) (outPaths []string, err error) {
	ride := make([]Waypoint, 0)
	last := Waypoint{}

	for _, waypoint := range waypoints {
		if waypoint.RideId != last.RideId {
			if last.RideId != "" {
				outPath := path.Join(outDirectory, last.ExportFilename())
				if err := saveRide(ride, outPath); err != nil {
					return outPaths, err
				}
				outPaths = append(outPaths, outPath)
				ride = make([]Waypoint, 0)
			}
			last = waypoint
		}
		ride = append(ride, waypoint)
	}

	// Don't forget to save the last set of records.
	// TODO(tstromberg): Remove duplicate code
	outPath := path.Join(outDirectory, last.ExportFilename())
	if err := saveRide(ride, outPath); err != nil {
		return outPaths, err
	}
	outPaths = append(outPaths, outPath)
	return
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: ./extract_strava </path/to/databases/strava> <output directory>")
	}

	waypoints, err := readWaypoints(os.Args[1])
	if err != nil {
		log.Fatal("Error reading waypoints: %s", err)
	}
	log.Printf("%d records found", len(waypoints))

	files, err := exportWaypoints(waypoints, os.Args[2])
	if err != nil {
		log.Fatal("Error exporting waypoints: %s", err)
	}
	log.Printf("Files written: %s", files)
}
