/*
Services Framework
github.com/dungw3b/services
*/
package main

import (
	"os"
	"fmt"
	"flag"
	"strconv"
	"io/ioutil"
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

func (s *SimpleService) Start() error {
	port := strconv.Itoa(GetConfigInt("api.port"))
	addr := services.GetConfigString("api.listen") +":"+ port
	s := &http.Server {
		Addr: addr,
	}

	s.conn.GET("/", s.handler)

	glog.Info("Start "+ a.Name() +" on "+ addr)
	if err := s.conn.StartServer(s); err != nil {
		return err
	}
	return nil
}

func (s *SimpleService) Stop() {
	if s.conn != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.conn.Shutdown(ctx); err != nil {
			glog.Error("Shutdown "+ a.Name() +" error ")
			glog.Error(err)
		}
		glog.Info("Stopped "+ a.Name())
	}
}

func (s *SimpleService) GetService() *SimpleService {
	return s
}

func (s *SimpleService) handler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func NewSimpleService() *SimpleService {
	return &SimpleService{}
}

func init() {
	os.Args = append(os.Args, "-logtostderr=true")
	path := flag.String("c", "", "full path to config file Ex. conf/app.json")
	flag.Parse()
	if len(*path) == 0 {
		fmt.Println("\nUsage:", os.Args[0], "-c conf/app.json");
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(*path)
	if err != nil {
		fmt.Println("Can not read configuration file "+ *path)
		os.Exit(1)
	}
	cfg, err := config.ParseJson(string(data))
	if err != nil {
		fmt.Println("Can not parse JSON configuration file "+ *path)
		os.Exit(1)
	}
	
	parseConfig(cfg)
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
		NewSimpleService()
	)
	defer services.Close()
	
}