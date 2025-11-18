package services

import "github.com/kardianos/service"
var serviceConfigs = service.Config{
	Name: "",
	DisplayName: "",
	Description: "",
	UserName: "",
}
var serviceInteface = service.Interface{
	Start: "",
	Stop: ""
}
func CreateService(){
	service.New()
}