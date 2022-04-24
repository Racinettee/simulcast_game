package game

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
)

type WindowConfig struct {
	Title         string
	Width, Height int
}

type Config struct {
	// If the file related to this struct should be reloaded into this stucture when its contents change
	HotReloadConfig bool
	// If internal game entities can hotload their assets when changed
	HotReloadAssets bool
	// Debugging aid - will draw collision geometry over entities
	DrawCollisionShapes bool
	// Resources
	ResourcePath string
	// The size of the window and its title
	Window WindowConfig
	// Reference to the game which will delegate out messages
	game *Game
	exit chan bool
}

func DefaultConfig(g *Game) Config {
	return Config{
		HotReloadConfig:     false,
		DrawCollisionShapes: false,
		ResourcePath:        ".",
		Window: WindowConfig{
			Title:  "Default Window Title",
			Width:  240,
			Height: 120,
		},
		game: g,
	}
}

func (conf *Config) StartWatching(what string) {
	if conf.HotReloadConfig {
		conf.exit = make(chan bool)
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println(err)
			return
		}
		watcher.Add(what)
		go func() {
			defer watcher.Close()
			for {
				select {
				// watch for file events
				case event := <-watcher.Events:
					log.Printf("%+v\n", event)
					if event.Op&fsnotify.Write == fsnotify.Write && event.Name == "simul_conf.toml" {
						newConf := DefaultConfig(conf.game)
						_, err := toml.DecodeFile("simul_conf.toml", &newConf)
						if err != nil {
							log.Println(err)
						}
						conf.game.eventQueue <- func() {
							*conf = newConf
							conf.game.ConfigChange()
						}
					}

				case err := <-watcher.Errors:
					log.Println(err)

				case <-conf.exit:
					return
				}
			}
		}()
	}
}

func (conf *Config) StopWatching() {
	if conf.exit != nil {
		conf.exit <- true
	}
}
