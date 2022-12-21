//go:build !windows
// +build !windows

package keyring

import (
	"context"
	"encoding/json"

	"github.com/99designs/keyring/cmd/k8sbridge"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	supportedBackends[K8Secret] = opener(func(cfg Config) (Keyring, error) {
		var err error

		pass := &k8secretKeyring{
			server: cfg.Server,
		}

		conn, err := grpc.Dial(cfg.Server, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err

		}
		//TODO connection close
		//defer conn.Close()

		pass.k8SHelperClient = k8sbridge.NewK8SHelperClient(conn)
		return pass, nil
	})
}

type k8secretKeyring struct {
	dir             string
	passcmd         string
	prefix          string
	server          string
	k8SHelperClient k8sbridge.K8SHelperClient
}

func (k *k8secretKeyring) Get(key string) (Item, error) {
	var in k8sbridge.K8SHelperRequest
	in.Keyname = key
	keypair, err := k.k8SHelperClient.GetKey(context.Background(), &in)
	data, err := json.Marshal(keypair)
	return Item{Key: key, Data: data}, err
}

func (k *k8secretKeyring) GetMetadata(key string) (Metadata, error) {
	return Metadata{}, nil
}

func (k *k8secretKeyring) Set(i Item) error {
	var in k8sbridge.K8SHelperRequest
	in.Keyname = i.Key

	_, err := k.k8SHelperClient.CreateKeypair(context.Background(), &in)

	return err
}

func (k *k8secretKeyring) Remove(key string) error {

	return nil
}

func (k *k8secretKeyring) Keys() ([]string, error) {
	return nil, nil
}
