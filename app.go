package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
)

type releaseLedger []*Release

func (rl releaseLedger) Len() int           { return len(rl) }
func (rl releaseLedger) Swap(i, j int)      { rl[i], rl[j] = rl[j], rl[i] }
func (rl releaseLedger) Less(i, j int) bool { return rl[i].Version < rl[j].Version }

// App represents an application deployed to Deis.
type App struct {
	// UUID is the unique identifier for the app.
	UUID    string        `json:"-"`
	ID      string        `json:"id"`
	Created time.Time     `json:"created"`
	Updated time.Time     `json:"updated"`
	LogPath string        `json:"-"`
	Ledger  releaseLedger `json:"-"`
}

// NewApp creates a new application with the given ID. If no ID is supplied, one will be
// automatically generated.
func NewApp(id string) *App {
	if id == "" {
		id = generateAppName()
	}
	app := &App{
		UUID:    uuid.New(),
		ID:      id,
		Created: time.Now(),
		Updated: time.Now(),
		LogPath: path.Join("/tmp", id+".log"),
	}
	// truncate or create the file
	f, err := os.Create(app.LogPath)
	if err != nil {
		log.WithFields(log.Fields{
			"app": app.ID,
		}).Errorf("could not create log file: %v", err)
		return app
	}
	defer f.Close()
	return app
}

func (a App) String() string {
	return a.ID
}

// Log stores an application message on disk, using the default formats for its operands.
func (a App) Log(message string) {
	f, err := os.OpenFile(a.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{
			"app": a.ID,
		}).Errorf("error opening file: %v", err)
		return
	}
	defer f.Close()
	log.WithFields(log.Fields{
		"app": a.ID,
	}).Info(message)
	buf := bytes.NewBufferString(fmt.Sprintf("%s deis[api]: %s\n", time.Now(), strings.TrimSpace(message)))
	io.Copy(f, buf)
}

// Logf stores an application message on disk, formatting according to a format specifier.
func (a App) Logf(message string, args ...interface{}) {
	a.Log(fmt.Sprintf(message, args))
}

// LatestRelease returns the most recent release in the ledger.
func (a App) LatestRelease() *Release {
	sort.Sort(sort.Reverse(a.Ledger))
	return a.Ledger[0]
}

// NewRelease appends a new release to the ledger using the provided build and config.
func (a *App) NewRelease(build *Build, config *Config) *Release {
	newVersion := a.LatestRelease().Version + 1
	release := &Release{
		Build:   build,
		Config:  config,
		Version: newVersion,
	}
	a.Ledger = append(a.Ledger, release)
	return release
}

// Rollback appends a new release to the ledger using the specified release's build + config.
func (a *App) Rollback(version int) error {
	if version < 1 {
		return errors.New("version cannot be below 0")
	}
	for _, r := range a.Ledger {
		if r.Version == version {
			release := a.NewRelease(r.Build, r.Config)
			return release.Publish()
		}
	}
	return errors.New("release not found")
}

func generateAppName() string {
	adjectives := []string{
		"ablest", "absurd", "actual", "allied", "artful", "atomic", "august",
		"bamboo", "benign", "blonde", "blurry", "bolder", "breezy", "bubbly",
		"candid", "casual", "cheery", "classy", "clever", "convex", "cubist",
		"dainty", "dapper", "decent", "deluxe", "docile", "dogged", "drafty",
		"earthy", "easier", "edible", "elfish", "excess", "exotic", "expert",
		"fabled", "famous", "feline", "finest", "flaxen", "folksy", "frozen",
		"gaslit", "gentle", "gifted", "ginger", "global", "golden", "grassy",
		"hearty", "hidden", "hipper", "honest", "humble", "hungry", "hushed",
		"iambic", "iconic", "indoor", "inward", "ironic", "island", "italic",
		"jagged", "jangly", "jaunty", "jiggly", "jovial", "joyful", "junior",
		"kabuki", "karmic", "keener", "kindly", "kingly", "klutzy", "knotty",
		"lambda", "leader", "linear", "lively", "lonely", "loving", "luxury",
		"madras", "marble", "mellow", "metric", "modest", "molten", "mystic",
		"native", "nearby", "nested", "newish", "nickel", "nimbus", "nonfat",
		"oblong", "offset", "oldest", "onside", "orange", "outlaw", "owlish",
		"padded", "peachy", "pepper", "player", "preset", "proper", "pulsar",
		"quacky", "quaint", "quartz", "queens", "quinoa", "quirky",
		"racing", "rental", "rising", "rococo", "rubber", "rugged", "rustic",
		"sanest", "scenic", "shadow", "skiing", "stable", "steely", "syrupy",
		"taller", "tender", "timely", "trendy", "triple", "truthy", "twenty",
		"ultima", "unbent", "unisex", "united", "upbeat", "uphill", "usable",
		"valued", "vanity", "velcro", "velvet", "verbal", "violet", "vulcan",
		"webbed", "wicker", "wiggly", "wilder", "wonder", "wooden", "woodsy",
		"yearly", "yeasty", "yeoman", "yogurt", "yonder", "youthy", "yuppie",
		"zaftig", "zanier", "zephyr", "zeroed", "zigzag", "zipped", "zircon",
	}

	nouns := []string{
		"anaconda", "airfield", "aqualung", "armchair", "asteroid", "autoharp",
		"babushka", "bagpiper", "barbecue", "bookworm", "bullfrog", "buttress",
		"caffeine", "chinbone", "countess", "crawfish", "cucumber", "cutpurse",
		"daffodil", "darkroom", "doghouse", "dragster", "drumroll", "duckling",
		"earthman", "eggplant", "electron", "elephant", "espresso", "eyetooth",
		"falconer", "farmland", "ferryman", "fireball", "footwear", "frosting",
		"gadabout", "gasworks", "gatepost", "gemstone", "goldfish", "greenery",
		"handbill", "hardtack", "hawthorn", "headwind", "henhouse", "huntress",
		"icehouse", "idealist", "inchworm", "inventor", "insignia", "ironwood",
		"jailbird", "jamboree", "jerrycan", "jetliner", "jokester", "joyrider",
		"kangaroo", "kerchief", "keypunch", "kingfish", "knapsack", "knothole",
		"ladybird", "lakeside", "lambskin", "larkspur", "lollipop", "lungfish",
		"macaroni", "mackinaw", "magician", "mainsail", "mongoose", "moonrise",
		"nailhead", "nautilus", "neckwear", "newsreel", "novelist", "nuthatch",
		"occupant", "offering", "offshoot", "original", "organism", "overalls",
		"painting", "pamphlet", "paneling", "pendulum", "playroom", "ponytail",
		"quacking", "quadrant", "queendom", "question", "quilting", "quotient",
		"rabbitry", "radiator", "renegade", "ricochet", "riverbed", "rucksack",
		"sailfish", "sandwich", "sculptor", "seashore", "seedcake", "stickpin",
		"tabletop", "tailbone", "teamwork", "teaspoon", "traverse", "turbojet",
		"umbrella", "underdog", "undertow", "unicycle", "universe", "uptowner",
		"vacation", "vagabond", "valkyrie", "variable", "villager", "vineyard",
		"waggoner", "waxworks", "waterbed", "wayfarer", "whitecap", "woodshed",
		"yachting", "yardbird", "yearbook", "yearling", "yeomanry", "yodeling",
		"zaniness", "zeppelin", "ziggurat", "zirconia", "zoologer", "zucchini",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf(
		"%s-%s",
		adjectives[r.Intn(len(adjectives))],
		nouns[r.Intn(len(nouns))],
	)
}
