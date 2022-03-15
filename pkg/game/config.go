package game

import "log"
import "github.com/fsnotify/fsnotify"
import "github.com/BurntSushi/toml"

type Config struct {
    // If the file related to this struct should be reloaded into this stucture when its contents change
    HotReloadConfig bool

	// Debug aid
	DrawCollisionShapes bool

	// Resources
	ResourcePath string

	exit chan bool
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
                        toml.DecodeFile("simul_conf.toml", &conf)
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
