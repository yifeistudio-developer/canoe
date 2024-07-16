package config

import (
	"github.com/kataras/golog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var client naming_client.INamingClient

func Register(cfg *Config, logger *golog.Logger) {
	server := cfg.Server
	var err error
	client, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	})
	if err != nil {
		logger.Fatal("register server failed.", err)
	}
	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: cfg.ApplicationName,
		Ip:          "192.168.80.1",
		Port:        server.Port,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Metadata:    map[string]string{},
	})
	if err != nil {
		logger.Fatal("register server failed.", err)
	}
	logger.Info("register server successfully.")
}

func DeRegister(cfg *Config, logger *golog.Logger) {
	server := cfg.Server
	instance, err := client.DeregisterInstance(vo.DeregisterInstanceParam{
		ServiceName: cfg.ApplicationName,
		Ip:          "192.168.80.1",
		Port:        server.Port,
	})
	if err != nil {
		logger.Fatal("deregister server failed.", err)
	}
	logger.Info("deregister server successfully. %v", instance)
}

func GetService(serviceName string) (model.Service, error) {
	return client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
}
