package api

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/pborman/uuid"
	"k8s.io/client-go/1.4/kubernetes"
	v1types "k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/rest"
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
	Ledger  releaseLedger `json:"-"`
}

// NewApp creates a new application with the given ID. If no ID is supplied, one will be
// automatically generated.
func NewApp(id string) (*App, error) {
	if id == "" {
		id = generateAppName()
	}
	app := &App{
		UUID:    uuid.New(),
		ID:      id,
		Created: time.Now(),
		Updated: time.Now(),
	}
	// create an initial release for the app
	app.NewRelease(nil, nil)
	// create a namespace for the app
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	namespace := &v1types.Namespace{
		ObjectMeta: v1types.ObjectMeta{
			Name: id,
			Labels: map[string]string{
				"heritage": "deis",
			},
		},
	}
	if _, err := clientset.Core().Namespaces().Create(namespace); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) String() string {
	return a.ID
}

// LatestRelease returns the most recent release in the ledger.
func (a *App) LatestRelease() *Release {
	if len(a.Ledger) == 0 {
		return nil
	}
	sort.Sort(sort.Reverse(a.Ledger))
	return a.Ledger[0]
}

// NewRelease appends a new release to the ledger using the provided build and config.
func (a *App) NewRelease(build *Build, config *Config) *Release {
	latestRelease := a.LatestRelease()
	if latestRelease == nil {
		latestRelease = &Release{
			App:     a,
			Version: 0,
		}
	}
	if build == nil {
		build = latestRelease.Build
	}
	if config == nil {
		config = latestRelease.Config
	}
	release := &Release{
		App:     a,
		Build:   build,
		Config:  config,
		Version: latestRelease.Version + 1,
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
