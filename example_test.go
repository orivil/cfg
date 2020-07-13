// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg_test

import (
	"fmt"
	"github.com/orivil/cfg"
	"github.com/orivil/service"
	"os"
)

func ExampleEnv_Unmarshal() {
	var data = `
[mysql]
host = "127.0.0.1"
port= "3306"
`

	// Config structure
	type mysql struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	}

	// Decode data to cfg.Env
	envs, err := cfg.Decode([]byte(data))
	if err != nil {
		panic(err)
	}

	// Set os environment. OPTIONAL
	err = os.Setenv("mysql.host", "localhost")
	if err != nil {
		panic(err)
	}

	// Get sub environment values, it use OS environment value if the value exist
	mysqlEnvs, err := envs.GetSub("mysql")
	if err != nil {
		panic(err)
	}
	fmt.Println(mysqlEnvs["host"])
	fmt.Println(mysqlEnvs["port"])

	config := mysql{}
	// Marshal mysqlEnvs to TOML data, and then Unmarshal data into config
	err = mysqlEnvs.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Host)
	fmt.Println(config.Port)

	config = mysql{}
	_ = os.Setenv("mysql.host", "127.0.0.1")
	// UnmarshalSub is shorthand for GetSub and Unmarshal
	err = envs.UnmarshalSub("mysql", &config)
	fmt.Println(config.Host)

	// Output:
	// localhost
	// 3306
	// localhost
	// 3306
	// 127.0.0.1
}

// NewService provide a configuration service, usually many services is depend on
// configuration service, this example shows how to inject configuration service
// to other services
func ExampleNewService() {
	var data = `
[redis]
addr = "127.0.0.1:3306"
`

	// Configuration service
	configService := cfg.NewService(cfg.NewMemoryStorageService(data))

	// Redis service, redis service depend on configuration service
	rdsService := newRedisService("redis", configService)

	// Dependency injection container
	container := service.NewContainer()

	// Get redis client from dependency injection container
	client, err := rdsService.Get(container)
	if err != nil {
		panic(err)
	}
	fmt.Println(client.addr)
	// Output:
	// 127.0.0.1:3306
}

type redisService struct {
	cfgService *cfg.Service
	namespace  string
	self       service.Provider
}

// Implement the service.Provider interface
func (r *redisService) New(ctn *service.Container) (value interface{}, err error) {
	// Get the singleton envs
	envs, err := r.cfgService.Get(ctn)
	if err != nil {
		return nil, nil
	}
	type config struct {
		Addr string `toml:"addr"`
	}
	c := config{}
	err = envs.UnmarshalSub(r.namespace, &c)
	if err != nil {
		return nil, err
	}
	return &redisClient{addr: c.Addr}, nil
}

// Get singleton redis client
func (r *redisService) Get(ctn *service.Container) (client *redisClient, err error) {
	// This functions means:
	// If ctn contains redisClient, return this redisClient
	// If ctn not contains redisClient, then execute r.New function, save the result to
	// container and return the result
	c, er := ctn.Get(&r.self)
	if er != nil {
		return nil, er
	} else {
		return c.(*redisClient), nil
	}
}

func newRedisService(namespace string, cfgService *cfg.Service) *redisService {
	s := &redisService{
		cfgService: cfgService,
		namespace:  namespace,
		self:       nil,
	}
	s.self = s
	return s
}

type redisClient struct {
	addr string
}
