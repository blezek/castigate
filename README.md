# castigate: podcast gateway

`castigate` or the podCAST GATEway, is a Go project to manage Podcasts for
portable MP3 players, or simply to download episodes of
podcasts from the command line.

# Basic Usage

Initialize a new `castigate.yaml` file:

```bash
castigate init

# take a look at the initial file
cat castigate.yaml
podcasts: []
filenametemplate: '{{.episode.Date.Format "2006-01-02-15-04-05" }}-{{.episode.Title}}.mp3'
defaultcounttokeep: 10
```

Add a podcast:

```bash
./castigate add simply_put https://simplyput.ligonier.org/rss
INFO[2024-06-24T20:22:10-05:00] loading castigate.yaml
INFO[2024-06-24T20:22:10-05:00] adding podcast: simply_put with feed https://simplyput.ligonier.org/rss to simply_put directory
```

Download podcasts:

```bash
./castigate sync
INFO[2024-06-24T20:22:41-05:00] loading castigate.yaml
INFO[2024-06-24T20:22:41-05:00] fetching feed from https://simplyput.ligonier.org/rss
INFO[2024-06-24T20:22:41-05:00] synchronizing Simply Put
...
```

# Background

`castigate` is designed to be a simple podcast management system.  
Configuration is handled through a YAML file, usually `castigate.yaml` in the
current directory.  While the `castigate.yaml` may be edited directly,
most functions are directly assessable through command line options.

These include:

```bash
./castigate help
57: ./castigate help
castigate or the podCAST GATEway, is a Go project to manage Podcasts for
portable MP3 players, or simply to download episodes of
podcasts from the command line.

Usage:
  castigate [command]

Available Commands:
  add         Adds a podcast to the configuration
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        initialize the config file
  list        list podcasts
  remove      remove podcasts
  sync        Download and sync podcasts

Flags:
  -c, --config string   path to config file (default "castigate.yaml")
  -d, --debug           set logging to debug
  -h, --help            help for castigate

Use "castigate [command] --help" for more information about a command.
```

Podcasts can be added using the `add` subcommand.  Each podcast has a unique label
associated with the podcast.  The `remove` command deletes the podcast, but leaves 
the `.mp3` files on disk.  The `sync` command updates each podcast.  Podcasts episodes
are downloaded until in order based on the Podcasts `counttokeep` number
or the global `defaultcounttokeep` if equal to 0.  Podcasts are downloaded
`oldest` first or `newest` first based on the `start` setting for each Podcast.
For instance, if `start` is `newest` and `counttokeep` is 10, when the Podcast
is downloaded for the first time, the most recent 10 episodes are downloaded.

`castigate` writes a playlist file in each directory using the title of the podcast and 
the `.m3u` extension.

After listening to episodes, simply delete the files from the corresponding directory, and
a new set of episodes, up to `counttokeep` will be downloaded at the next `sync`.

# Filename format

The `filenametemplate` is a [Go text template](https://pkg.go.dev/text/template).  The variables
exposed to the template include the current `Episdode` and `Item`:

```go
type Item struct {
	Title           string                  
	Description     string                  
	Content         string                  
	Link            string                  
	Links           []string                
	Updated         string                  
	UpdatedParsed   *time.Time              
	Published       string                  
	PublishedParsed *time.Time              
	Author          *Person                 
	Authors         []*Person               
	GUID            string                  
	Image           *Image                  
	Categories      []string                
	Enclosures      []*Enclosure            
	DublinCoreExt   *ext.DublinCoreExtension
	ITunesExt       *ext.ITunesItemExtension
	Extensions      ext.Extensions          
	Custom          map[string]string       
}

type Episode struct {
	GUID         string
	URL          string
	State        EpisodeState
	Filename     string
	Date         time.Time
	PodcastLabel string
}
```

The default format is `{{.episode.Date.Format "2006-01-02-15-04-05" }}-{{.episode.Title}}.mp3`, the
episode date followed by the Episode Title.

# License

BSD 3-clause license.