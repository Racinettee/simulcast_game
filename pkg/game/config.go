package game

import "log"
import "github.com/fsnotify/fsnotify"
import "github.com/BurntSushi/toml"

type WindowConfig struct {
    Title string
    Width, Height int
}

type Config struct {
    // If the file related to this struct should be reloaded into this stucture when its contents change
    HotReloadConfig bool

	// Debug aid
	DrawCollisionShapes bool

	// Resources
	ResourcePath string

	Window WindowConfig

	game *Game
	exit chan bool
}

func DefaultConfig(g *Game) Config {
    return Config {
        HotReloadConfig: false,
        DrawCollisionShapes: false,
        ResourcePath: ".",
        Window: WindowConfig {
            Title: "Default Window Title",
            Width: 240,
            Height: 120,
        },
        game: g,
    }
}

func (conf *Config) StartWatching(what string) {

    if conf.HotReloadConfig {
        conf.exit = make(chan bool)
        log.Println("creating watcher")
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
            log.Println(err)
            return
        }
        log.Println("starting config watcher")
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
                        conf.game.eventQueue<- func() {
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
        conf.exit<- true
    }
}
