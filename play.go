package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/jameycribbs/hare"
	"github.com/jameycribbs/hare/datastores/disk"
	"github.com/jameycribbs/hare/examples/crud/models"
)

func main() {
	ds, err := disk.New("./data", ".json")
	if err != nil {
		panic(err)
	}

	db, err := hare.New(ds)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//----- CREATE -----

	recID, err := db.Insert("episodes", &models.Episode{
		Season:           6,
		Episode:          19,
		Film:             "Red Zone Cuba",
		Shorts:           []string{"Speech:  Platform, Posture, and Appearance"},
		YearFilmReleased: 1966,
		DateEpisodeAired: time.Date(1994, 12, 17, 0, 0, 0, 0, time.UTC),
		HostID:           2, // See associated Host model
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("New record id is:", recID)

	//----- READ -----

	rec := models.Episode{}

	err = db.Find("episodes", 4, &rec)
	if err != nil {
		panic(err)
	}

	// Notice that this is using the benefits of the associated
	// Host model to print the host's name.
	fmt.Printf("Found record is %v and it was hosted by %v\n", rec.Film, rec.Host.Name)

	//----- UPDATE -----

	rec.Film = "The Skydivers - The Final Cut"
	if err = db.Update("episodes", &rec); err != nil {
		panic(err)
	}

	//----- DELETE -----

	err = db.Delete("episodes", 2)
	if err != nil {
		panic(err)
	}

	//----- QUERYING -----

	results, err := models.QueryEpisodes(db, func(r models.Episode) bool {
		// Notice that we are taking advantage of the
		// code we put in the Episode AfterFind method
		// to be able to do the query by the associated
		// host's name.
		return r.Host.Name == "Joel"
	}, 0)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		// Again, we are able to automatically use the host's name, because the
		// embedd Host struct was populated in the AfterFind method.
		fmt.Printf("%v hosted the season %v episode %v film, '%v'\n", r.Host.Name, r.Season, r.Episode, r.Film)

		// Here we are once again taking advantage of the code we put in the Episode
		// AfterFind method that automatically populates the episode's Comments struct
		// with associated records from the comments table.
		for _, c := range r.Comments {
			fmt.Printf("\t-- Comment for episode %v: %v\n", r.Episode, c.Text)
		}
	}
}

func init() {
	cmd := exec.Command("cp", "./data/episodes_default.txt", "./data/episodes.json")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
