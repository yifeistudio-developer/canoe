package config

import (
	"canoe/internal/util"
	"github.com/kataras/golog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var client naming_client.INamingClient

func Register(cfg *Config, log *golog.Logger) {
	server := cfg.Server
	var err error
	client, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	})
	address, err := util.GetLocalIpAddress()
	if err != nil {
		log.Fatal("register server failed.", err)
	}
	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: cfg.ApplicationName,
		Ip:          address,
		Port:        server.Port,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Metadata: map[string]string{
			"gRPC_port": "50051",
		},
	})
	if err != nil {
		log.Fatal("register server failed.", err)
	}
	log.Info("register server successfully.")
}

func DeRegister(cfg *Config, log *golog.Logger) {
	server := cfg.Server
	address, _ := util.GetLocalIpAddress()
	instance, err := client.DeregisterInstance(vo.DeregisterInstanceParam{
		ServiceName: cfg.ApplicationName,
		Ip:          address,
		Port:        server.Port,
	})
	if err != nil {
		log.Fatal("deregister server failed.", err)
	}
	log.Infof("deregister server successfully. %v", instance)
}

func GetService(serviceName string) (model.Service, error) {
	return client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
}
