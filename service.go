/*
Services Framework
github.com/dungw3b/services
*/
package services

import (
)

type Service interface {
	Name() string // return service name
	Init() // initial service
	Start() error // start service
	Stop() // graceful stop service
	GetService() interface{} // get instance
	ReloadData() // reload configuration
}