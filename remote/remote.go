// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package remote integrates the remote features of Viper.
package remote

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"github.com/Be-MobileNV/viper"
	"github.com/coreos/etcd/clientv3"
)

type remoteConfigProvider struct {
	client *clientv3.Client
}

func (rc remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	if rc.client == nil {
		var err error
		rc.client, err = getConfigManager(rp)
		if err != nil {
			return nil, err
		}
	}
	r, err := rc.client.Get(context.Background(), rp.Path())
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(r.Kvs[0].Value), nil
}

func (rc remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	if rc.client == nil {
		var err error
		rc.client, err = getConfigManager(rp)
		if err != nil {
			return nil, err
		}
	}
	w := rc.client.Watch(context.Background(), rp.Path())
	resp := <-w
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	val := resp.Events[0].Kv.Value
	return bytes.NewReader(val), nil
}

func (rc remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (responseChannel <-chan *viper.RemoteResponse, quitwc chan bool, err error) {
	quitwc = make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)
	if rc.client == nil {
		var err error
		rc.client, err = getConfigManager(rp)
		if err != nil {
			return nil, nil, err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	w := rc.client.Watch(ctx, rp.Path())
	go func(etcdResponseChannel clientv3.WatchChan, vr chan<- *viper.RemoteResponse, quitwc <-chan bool, cancel context.CancelFunc) {
		for {
			select {
			case <-quitwc:
				cancel()
				return
			case resp := <-etcdResponseChannel:
				vr <- &viper.RemoteResponse{
					Error: resp.Err(),
					Value: resp.Events[0].Kv.Value,
				}

			}

		}
	}(w, viperResponsCh, quitwc, cancel)
	return viperResponsCh, quitwc, nil
}

func getConfigManager(rp viper.RemoteProvider) (*clientv3.Client, error) {
	etcdConfig := clientv3.Config{
		Endpoints:   strings.Split(rp.Endpoint(), ","),
		DialTimeout: 5 * time.Second,
		Username:    rp.Username(),
		Password:    rp.Password(),
	}
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func init() {
	viper.RemoteConfig = &remoteConfigProvider{}
}
