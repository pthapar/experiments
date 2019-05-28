package main

import (
	"context"
	"flag"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/glog"
	"go.etcd.io/etcd/clientv3"
)

const (
	RefreshInterval = time.Second * 5
	LeaderKey       = "leader"
)

var (
	nodeID = flag.Int("nodeID", 1, "")
)

// TODO: Revoke lease using context cancellation
func renewLease(cli *clientv3.Client, leaseID clientv3.LeaseID) {
	glog.Infof("[node-%d]renewing lease: %+v", *nodeID, leaseID)
	ch, err := cli.Lease.KeepAlive(context.TODO(), leaseID)
	if err != nil {
		glog.Errorf("keep alive for lease failed")
	}

	for {
		select {
		case ttl := <-ch:
			glog.Infof("ttl resp: %+v", ttl.TTL)
		}
	}
}

// ProposeLeader proposes the current node as the leader
func ProposeLeader(cli *clientv3.Client) {
	rand.Seed(time.Now().UnixNano())
	interval := time.Duration(500 + rand.Intn(500))
	glog.Infof("[node-%d]interval: %d", *nodeID, interval)
	tick := time.NewTicker(time.Millisecond * interval)
	for {
		select {
		case <-tick.C:
			txn := cli.Txn(context.TODO())
			lease, _ := cli.Grant(context.TODO(), 10)
			txresp, err := txn.If(
				clientv3.Compare(clientv3.CreateRevision(LeaderKey), "=", 0),
			).Then(
				clientv3.OpPut(LeaderKey, strconv.Itoa(*nodeID), clientv3.WithLease(lease.ID)),
			).Commit()
			if err != nil {
				glog.Errorf("failed to propose leader. %s", err.Error())
			}
			if txresp.Succeeded {
				go renewLease(cli, lease.ID)
				glog.Infof("resp: %v, %+v", txresp.Succeeded, *txresp.Header)
			}
			for _, resp := range txresp.Responses {
				glog.Infof("[node-%d]resp: %+v", *nodeID, resp.String())
			}
		default:
		}
	}
}

func main() {
	// glog.Info("connecting to etcd")
	flag.Parse()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	// glog.Infof("create client")
	if err != nil {
		glog.Exitf("failed to connect to etcd servers. %s", err.Error())
	}
	ProposeLeader(cli)
	// cli.Txn()
}
