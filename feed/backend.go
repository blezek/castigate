package feed

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type BackendInterface interface {
	Init(filename string)
	Load() (Config, error)
	Save(Config) error
}

type FileBackend struct {
	Filename string
}
type SqliteBackend struct {
	Database *sql.DB
	Filename string
}

func (b *SqliteBackend) Init(filename string) {
	b.Filename = filename
	var err error
	b.Database, err = sql.Open("sqlite3", b.Filename)
	if err != nil {
		log.Fatalf("could not open database %s: %v", b.Filename, err)
	}
	// make the config table
	b.Database.Exec("create table if not exists config ( FilenameTemplate string, DefaultCountToKeep int )")
	b.Database.Exec("create table if not exists podcast ( Label string PRIMARY KEY, Feed string, Directory string, CountToKeep int, Start string )")
	b.Database.Exec("create table if not exists episode ( GUID string, URL string, State int, Filename string, Date timestamp, PodcastLabel string )")

}

func (b *SqliteBackend) Load() (Config, error) {
	config := Config{
		Podcasts:           make([]*Podcast, 0),
		FilenameTemplate:   `{{.episode.Date.Format "2006-01-02-15:04:05" }}-{{.item.Title}}.mp3`,
		DefaultCountToKeep: 10,
	}
	b.Database.QueryRow("select FilenameTemplate, DefaultCountToKeep from config").Scan(&config.FilenameTemplate, &config.DefaultCountToKeep)
	rows, err := b.Database.Query("select Label, Feed, Directory, CountToKeep, Start from podcast")
	if err != nil {
		log.Fatalf("could not read from config table %v", err)
	}
	for rows.Next() {
		podcast := Podcast{Episodes: make(map[string]*Episode, 0)}
		rows.Scan(&podcast.Label, &podcast.Feed, &podcast.Directory, &podcast.CountToKeep, &podcast.Start)
		config.Podcasts = append(config.Podcasts, &podcast)
	}

	labelIndex := map[string]int{}
	for i, podcast := range config.Podcasts {
		labelIndex[podcast.Label] = i
	}

	rows, err = b.Database.Query("select GUID, URL, State, Filename, Date, PodcastLabel from episode")
	if err != nil {
		log.Fatalf("could not read from episode table %v", err)
	}
	for rows.Next() {
		episode := Episode{}
		rows.Scan(&episode.GUID, &episode.URL, &episode.State, &episode.Filename, &episode.Date, &episode.PodcastLabel)
		//id := labelIndex[episode.PodcastLabel]
		// config.Podcasts[id].Episodes = append(config.Podcasts[id].Episodes, &episode)
	}
	return config, nil
}
func (b *SqliteBackend) Save(config Config) error {
	transaction, err := b.Database.Begin()
	if err != nil {
		log.Fatalf("could not start transaction: %v", err)
	}

	_, err = b.Database.Exec("delete from config")
	if err != nil {
		transaction.Rollback()
		return err
	}
	_, err = b.Database.Exec("delete from podcast")
	if err != nil {
		transaction.Rollback()
		return err
	}
	_, err = b.Database.Exec("delete from episode")
	if err != nil {
		transaction.Rollback()
		return err
	}

	_, err = b.Database.Exec("insert into config ( FilenameTemplate, DefaultCountToKeep ) values ( ?, ? )", config.FilenameTemplate, config.DefaultCountToKeep)
	if err != nil {
		transaction.Rollback()
		return err
	}
	for _, podcast := range config.Podcasts {
		_, err = b.Database.Exec("insert into podcast ( Label, Feed, Directory, CountToKeep, Start )", podcast.Label, podcast.Feed, podcast.Directory, podcast.CountToKeep, podcast.Start)
		if err != nil {
			transaction.Rollback()
			return err
		}
		for _, episode := range podcast.Episodes {
			episode.PodcastLabel = podcast.Label
			_, err = b.Database.Exec("insert into episode ( GUID, URL, State, Filename, Date, PodcastLabel ) values ( ?, ?, ?, ?, ?, ? )",
				episode.GUID, episode.URL, episode.State, episode.Filename, episode.Date, episode.PodcastLabel)
			if err != nil {
				transaction.Rollback()
				return err
			}
		}
	}
	return transaction.Commit()
}

func (b *FileBackend) Init(filename string) {
	b.Filename = filename
}
func (b *FileBackend) Load() (Config, error) {

	log.Infof("loading %s", b.Filename)
	// read the yaml file
	contents, err := os.ReadFile(b.Filename)
	if err != nil {
		log.Errorf("failed to read input: %s", err)
		return Config{}, err
	}
	config := Config{
		Podcasts:           make([]*Podcast, 0),
		FilenameTemplate:   `{{.episode.Date.Format "2006-01-02-15:04:05" }}-{{.item.Title}}.mp3`,
		DefaultCountToKeep: 10,
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		log.Errorf("failed to load YAML: %s", err)
		return Config{}, err
	}
	return config, nil
}

func (b *FileBackend) Save(config Config) error {
	//
	//log.Infof("saving backup to %s", b.Filename+".bak")
	//err = os.WriteFile(configFile+".bak", b, 0644)
	//if err != nil {
	//	log.Fatalf("failed to write backup config: %s", err)
	//}
	buffer, err := yaml.Marshal(config)
	if err != nil {
		log.Errorf("failed to marshal YAML: %s", err)
		return err
	}

	// Write the new config
	err = os.WriteFile(b.Filename, buffer, 0644)
	if err != nil {
		log.Errorf("failed to save YAML: %s", err)
		return err
	}
	return err
}
