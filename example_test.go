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

func ExampleEnv_UnmarshalSub() {
	// TOML data
	var data = `[mysql]
host = "127.0.0.1"
port= "3306"`
	type mysql struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	}
	envs, err := cfg.Decode([]byte(data))
	if err != nil {
		panic(err)
	} else {
		env := mysql{}

		// Set os environment
		err = os.Setenv("mysql.host", "localhost")
		if err != nil {
			panic(err)
		}
		// UnmarshalSub or GetSub will load os environment values

		err = envs.UnmarshalSub("mysql", &env)
		if err != nil {
			panic(err)
		}
		fmt.Println(env.Host)
		fmt.Println(env.Port)
	}

	// Output:
	// localhost
	// 3306
}

func ExampleNewService() {

	// TOML data
	var configData = `
[mysql]
host = "127.0.0.1"
port= "3306"
`
	configService := cfg.NewService(cfg.NewMemoryStorageService(configData))

	container := service.NewContainer()

	envs, err := configService.Get(container)
	if err != nil {
		panic(err)
	}
	envs, err = envs.GetSub("mysql")
	for s, i := range envs {
		fmt.Printf("%s = %s\n", s, i)
	}
	// Output:
	// host = 127.0.0.1
	// port = 3306
}
