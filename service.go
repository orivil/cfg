// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"github.com/orivil/service"
)

type Service struct {
	storageService StorageService
	self           service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var store Storage
	store, err = s.storageService.Get(ctn)
	if err != nil {
		return nil, err
	}
	var data []byte
	data, err = store.GetTomlData()
	if err != nil {
		return nil, err
	}
	var env Env
	env, err = Decode(data)
	if err != nil {
		return nil, err
	}
	return env, nil
}

func (s *Service) Get(ctn *service.Container) (envs Env, err error) {
	es, er := ctn.Get(&s.self)
	if er != nil {
		return nil, er
	} else {
		return es.(Env), nil
	}
}

func NewService(storageService StorageService) *Service {
	s := &Service{storageService: storageService}
	s.self = s
	return s
}
