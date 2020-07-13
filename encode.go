// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"github.com/BurntSushi/toml"
)

func Decode(data []byte) (Env, error) {
	env := make(Env)
	err := toml.Unmarshal(data, &env)
	return env, err
}
