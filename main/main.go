/*
Services Framework
github.com/dungw3b/services
*/
package main

import (
	"time"
	"context"
	"strconv"
	"net/http"
	"github.com/golang/glog"
	"github.com/olebedev/config"
	"github.com/dungw3b/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type SimpleService struct {
	conn *echo.Echo
}

func (s *SimpleService) Name() string {
	return "Simple Service"
}

func (s *SimpleService) Init() {
	s.conn = echo.New()
	s.conn.Logger.SetLevel(log.OFF)
	s.conn.HideBanner = true
	s.conn.HidePort = true
}

func (s *SimpleService) ReloadData() {
	glog.Info("Reloaded "+ s.Name(), " data")
}

func (s *SimpleService) Start() error {
	port := strconv.Itoa(services.GetConfigInt("service.port"))
	addr := services.GetConfigString("service.listen") +":"+ port
	/*server := &http.Server {
		Addr: addr,
	}*/

	s.conn.GET("/", s.handler)

	glog.Info("Started "+ s.Name() +" on "+ addr)
	/*if err := s.conn.StartServer(server); err != nil {
		return err
	}*/
	s.conn.Start(addr)
	return nil
}

func (s *SimpleService) Stop() {
	if s.conn != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.conn.Shutdown(ctx); err != nil {
			glog.Error("Shutdown "+ s.Name() +" error ")
			glog.Error(err)
		}
		glog.Info("Stopped "+ s.Name())
	}
}

func (s *SimpleService) GetService() interface{} {
	return s
}

func (s *SimpleService) handler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func NewSimpleService() *SimpleService {
	return &SimpleService{}
}

func init() {
	services.Init(parseConfig)
}

func parseConfig(cfg *config.Config) {
	var (
		strval string
		intval int
	)

	// service
	strval = cfg.UString("service.listen")
	if len(strval) == 0 {
		glog.Exit("Can not read config service.listen")
	}
	services.SetConfig("service.listen", strval)

	intval = cfg.UInt("service.port")
	if intval <= 0 {
		glog.Exit("Can not read config service.port")
	}
	services.SetConfig("service.port", intval)
}

func main() {
	
	services.Run(
		NewSimpleService(),
	)

}