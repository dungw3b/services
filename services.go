/*
Services Framework
github.com/dungw3b/services
*/
package services

import (
	"os"
	"fmt"
	"flag"
	"sync"
	"time"
	"syscall"
	"io/ioutil"
	"os/signal"
	"github.com/dungw3b/glog"
	"github.com/dungw3b/config"
)

var (
	Config map[string]interface{}
	Services []Service

	waitgroup sync.WaitGroup
	configPath string
	parseConfig ParseConfigFunc
	iStop bool
)

type ParseConfigFunc func(cfg *config.Config)

func init() {
	Config = make(map[string]interface{})
	iStop = false
	
	// manage signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	//signal.Ignore(syscall.SIGHUP)
	go func() {
		for {
			s := <-sig
			switch s {
			
			// reload configuration
			case syscall.SIGHUP:
				glog.Info("---Reload configuration---")
				startFail := false
				for _, service := range Services {
					iStop = true
					service.Stop()
					service.Init()
					service.ReloadData()
					go func() {
						if err := service.Start(); err != nil {
							if !iStop {
								glog.Error("Start ", service.Name(), " error: ", err)
								startFail = true
							}
						}
					}()
					time.Sleep(200 * time.Millisecond)
					
					if startFail {
						waitgroup.Done()
						return
					}
					iStop = false
				}
			
			// graceful stop
			default:
				glog.Warning("Graceful stop")
				iStop = true
				for _, service := range Services {
					service.Stop()
				}
				waitgroup.Done()
			}
		}
	}()
}

func Init(parser ParseConfigFunc) {
	os.Args = append(os.Args, "-logtostderr=true")
	path := flag.String("c", "", "full path to config file Ex. conf/app.json")
	flag.Parse()
	if len(*path) == 0 {
		fmt.Println("\nUsage:", os.Args[0], "-c conf/app.json");
		os.Exit(1)
	}

	configPath = *path
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		glog.Fatal("Can not read configuration file ", *path)
		os.Exit(1)
	}
	cfg, err := config.ParseJson(string(data))
	if err != nil {
		glog.Fatal("Can not parse JSON configuration file ", *path)
		os.Exit(1)
	}

	parseConfig = parser
	parseConfig(cfg)
}

func GetService(name string) interface{} {
	for _, service := range Services {
		if service.Name() == name {
			//glog.Info(service)
			return service
		}
	}
	return nil
}

func reloadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Error("Can not read configuration file ", path)
		return
	}
	cfg, err := config.ParseJson(string(data))
	if err != nil {
		glog.Error("Can not parse JSON configuration file ", err)
		return
	}

	parseConfig(cfg)
	glog.Info("Reloaded configuration file")
}

func registerService(service Service) {
	Services = append(Services, service)
}

func Run(services ...Service) {
	waitgroup.Add(1)
	startFail := false
	for _, service := range services {
		registerService(service)
		service.Init()
		go func() {
			if err := service.Start(); err != nil {
				if !iStop {
					glog.Error("Start ", service.Name(), " error: ", err)
					startFail = true
				}
			}
		}()
		time.Sleep(200 * time.Millisecond)
		
		if startFail {
			waitgroup.Done()
			return
		}
	}

	// wait
	waitgroup.Wait()
	glog.Flush()
	//time.Sleep(time.Second)
}

func SetConfig(name string, value interface{}) {
	Config[name] = value
}

func GetConfigString(name string) string {
	if _,found := Config[name]; found {
		return Config[name].(string)
	}
	glog.Error("Config ", name, " not found, return empty string")
	return ""
}

func GetConfigInt(name string) int {
	if _,found := Config[name]; found {
		return Config[name].(int)
	}
	glog.Error("Config ", name, " not found, return int 0")
	return 0
}
