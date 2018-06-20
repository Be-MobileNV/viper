// +build integration

package remote

import (
	"fmt"
	"testing"
)

type TestRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
	username      string
	password      string
}

func (rp TestRemoteProvider) Provider() string {
	return rp.provider
}

func (rp TestRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp TestRemoteProvider) Path() string {
	return rp.path
}

func (rp TestRemoteProvider) SecretKeyring() string {
	return rp.secretKeyring
}

func (rp TestRemoteProvider) Username() string {
	return rp.username
}

func (rp TestRemoteProvider) Password() string {
	return rp.password
}

func TestGet(t *testing.T) {
	request := TestRemoteProvider{
		endpoint: "http://127.0.0.1:2379",
		path:     "test",
	}
	w := &remoteConfigProvider{}
	result, err := w.Get(request)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	fmt.Println(result)
}

func TestWatch(t *testing.T) {
	request := TestRemoteProvider{
		endpoint: "http://127.0.0.1:2379",
		path:     "test",
	}
	w := &remoteConfigProvider{}
	result, err := w.Watch(request)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	fmt.Println(result)
}

func TestWatchChannel(t *testing.T) {
	request := TestRemoteProvider{
		endpoint: "http://127.0.0.1:2379",
		path:     "test",
	}
	w := &remoteConfigProvider{}
	result, _, err := w.WatchChannel(request)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i := 0; i < 3; i++ {
		<-result
	}

	fmt.Println(result)
}
