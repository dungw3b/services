/*
Services Framework
github.com/dungw3b/services
*/
package services

import (
	"os"
	"sync"
	"time"
	"syscall"
	"os/signal"
	"github.com/golang/glog"
)

var (
	Config map[string]interface{}
	waitgroup sync.WaitGroup
	Services []Service
)

func init() {
	waitgroup.Add(1)
	
	// manage signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	//signal.Ignore(syscall.SIGHUP)
	go func() {
		for {
			s := <-sig
			switch s {
			/*
			// reload configuration
			case syscall.SIGHUP:
				glog.Info("---Reload configuration---")
				ReloadConfig()
				glog.Info("--------------------------")
			*/
			// graceful stop
			default:
				glog.Warning("Graceful stop")
				for _, service := range Services {
					service.Stop()
				}
				waitgroup.Done()
			}
		}
	}()
}

func registerService(service Service) {
	Services = append(Services, service)
}

func Run(services ...Service) {
	for _, service := range services {
		registerService(service)
		service.Init()
		go func() {
			if err := service.Start(); err != nil {
				glog.Error("Service ", service.Name(), " error ", err)
				waitgroup.Done()
			}
		}()
		time.Sleep(100 * time.Millisecond)
	}

	// wait
	waitgroup.Wait()
	glog.Flush()
	time.Sleep(time.Second)
}

func SetConfig(name string, value interface{}) {
	Config[name] = value
}

func GetConfigString(name string) string {
	return Config[name].(string)
}

func GetConfigInt(name string) int {
	return Config[name].(int)
}