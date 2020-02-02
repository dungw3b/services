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

/* First Service */

type FirstService struct {
	conn *echo.Echo
}

func (s *FirstService) Name() string {
	return "First Service"
}

func (s *FirstService) Init() {
	s.conn = echo.New()
	s.conn.Logger.SetLevel(log.OFF)
	s.conn.HideBanner = true
	s.conn.HidePort = true
}

func (s *FirstService) ReloadData() {
	glog.Info("Reloaded "+ s.Name(), " data")
}

func (s *FirstService) Start() error {
	port := strconv.Itoa(services.GetConfigInt("firstservice.port"))
	addr := services.GetConfigString("firstservice.listen") +":"+ port
	/*server := &http.Server {
		Addr: addr,
	}*/

	s.conn.GET("/", s.handler)

	glog.Info("Started "+ s.Name() +" on "+ addr)
	/*if err := s.conn.StartServer(server); err != nil {
		return err
	}*/
	return s.conn.Start(addr)
}

func (s *FirstService) Stop() {
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

func (s *FirstService) GetInstance() interface{} {
	return s
}

func (s *FirstService) handler(c echo.Context) error {
	return c.String(http.StatusOK, "First service!")
}

/* END First Service */

/* Second Service */

type SecondService struct {
	conn *echo.Echo
}

func (s *SecondService) Name() string {
	return "Second Service"
}

func (s *SecondService) Init() {
	s.conn = echo.New()
	s.conn.Logger.SetLevel(log.OFF)
	s.conn.HideBanner = true
	s.conn.HidePort = true
}

func (s *SecondService) ReloadData() {
	glog.Info("Reloaded "+ s.Name(), " data")
}

func (s *SecondService) Start() error {
	port := strconv.Itoa(services.GetConfigInt("secondservice.port"))
	addr := services.GetConfigString("secondservice.listen") +":"+ port
	/*server := &http.Server {
		Addr: addr,
	}*/

	s.conn.GET("/", s.handler)

	glog.Info("Started "+ s.Name() +" on "+ addr)
	/*if err := s.conn.StartServer(server); err != nil {
		return err
	}*/
	return s.conn.Start(addr)
}

func (s *SecondService) Stop() {
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

func (s *SecondService) GetInstance() interface{} {
	return s
}

func (s *SecondService) handler(c echo.Context) error {
	return c.String(http.StatusOK, "Second service!")
}

/* END Second Service */


/* Main function */

func NewFirstService() *FirstService {
	return &FirstService{}
}

func NewSecondService() *SecondService {
	return &SecondService{}
}

func init() {
	services.Init(parseConfig)
}

func parseConfig(cfg *config.Config) {
	var (
		strval string
		intval int
	)

	// first service
	strval = cfg.UString("firstservice.listen")
	if len(strval) == 0 {
		glog.Exit("Can not read config firstservice.listen")
	}
	services.SetConfig("firstservice.listen", strval)

	intval = cfg.UInt("firstservice.port")
	if intval <= 0 {
		glog.Exit("Can not read config firstservice.port")
	}
	services.SetConfig("firstservice.port", intval)

	// second service
	strval = cfg.UString("secondservice.listen")
	if len(strval) == 0 {
		glog.Exit("Can not read config secondservice.listen")
	}
	services.SetConfig("secondservice.listen", strval)

	intval = cfg.UInt("secondservice.port")
	if intval <= 0 {
		glog.Exit("Can not read config secondservice.port")
	}
	services.SetConfig("secondservice.port", intval)
}

func main() {
	
	services.Run(
		NewFirstService(),
		NewSecondService(),
	)

}