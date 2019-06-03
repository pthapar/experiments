package main

import (
	"context"
	"flag"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/glog"
)

const (
	RefreshInterval = time.Second * 5
	LeaderKey       = "leader"
)

var (
	nodeID   = flag.Int("nodeID", 1, "")
	isLeader = false
	mu       = &sync.Mutex{}
)

// SetLeader ..
func SetLeader() {
	mu.Lock()
	defer mu.Unlock()
	isLeader = true

}

// UnsetLeader ..
func UnsetLeader() {
	mu.Lock()
	defer mu.Unlock()
	isLeader = false
}

func IsLeader() bool {
	mu.Lock()
	defer mu.Unlock()
	return isLeader
}

// TODO: Revoke lease using context cancellation
func renewLease(cli *clientv3.Client, leaseID clientv3.LeaseID) {
	glog.Infof("[node-%d]renewing lease: %+v", *nodeID, leaseID)
	ch, err := cli.Lease.KeepAlive(context.TODO(), leaseID)
	if err != nil {
		glog.Errorf("keep alive for lease failed")
	}
	// TODO: close the lease before returning
	for {
		select {
		case ttl := <-ch:
			glog.Infof("ttl resp: %+v", ttl.TTL)
		}
	}
}

func bringUpClusterIP() {
	glog.Info("[node-%d] bring up cluster ip noop done.")
}

// returns true if the node is elected as leader
func proposeSelfAsLeader(cli *clientv3.Client, leaseID clientv3.LeaseID) bool {
	glog.Infof("proposing self as leader")
	txn := cli.Txn(context.TODO())
	txresp, err := txn.If(
		clientv3.Compare(clientv3.CreateRevision(LeaderKey), "=", 0),
	).Then(
		clientv3.OpPut(LeaderKey, strconv.Itoa(*nodeID), clientv3.WithLease(leaseID)),
	).Commit()
	if err != nil {
		glog.Errorf("failed to propose leader. %s", err.Error())
	}

	// Only trigerred when this node is elected
	if txresp.Succeeded {
		glog.Infof("[node-%d]leader elected", *nodeID)
		glog.Infof("resp: %v, %+v", txresp.Succeeded, *txresp.Header)
	}
	return txresp.Succeeded
}

func main() {
	glog.Info("connecting to etcd")
	flag.Parse()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	// glog.Infof("create client")
	if err != nil {
		glog.Exitf("failed to connect to etcd servers. %s", err.Error())
	}
	// Keep the lease till alive till this node is alive
	lease, _ := cli.Grant(context.TODO(), 5)
	glog.Infof("[node-%d]lease object: %+v", *nodeID, *lease)
	go renewLease(cli, lease.ID)
	if proposeSelfAsLeader(cli, lease.ID) {
		SetLeader()
		bringUpClusterIP()
	}

	w := cli.Watch(context.Background(), LeaderKey)

	for ws := range w {
		for _, e := range ws.Events {
			if e.Type == mvccpb.DELETE {
				glog.Infof("leader key deleted")
				// Key deleted
				if proposeSelfAsLeader(cli, lease.ID) {
					SetLeader()
					bringUpClusterIP()
				}
			}
		}
	}
	// TODO: Handle SIGINT to improve failover faster

	// cli.Txn()
}

// ProposeLeader proposes the current node as the leader
// func ProposeLeader(cli *clientv3.Client) {
// 	rand.Seed(time.Now().UnixNano())
// 	interval := time.Duration(500 + rand.Intn(500))
// 	glog.Infof("[node-%d]interval: %d", *nodeID, interval)
// 	tick := time.NewTicker(time.Millisecond * interval)

// 	for {
// 		select {
// 		case <-tick.C:
// 			txn := cli.Txn(context.TODO())
// 			txresp, err := txn.If(
// 				clientv3.Compare(clientv3.CreateRevision(LeaderKey), "=", 0),
// 			).Then(
// 				clientv3.OpPut(LeaderKey, strconv.Itoa(*nodeID), clientv3.WithLease(lease.ID)),
// 			).Commit()
// 			if err != nil {
// 				glog.Errorf("failed to propose leader. %s", err.Error())
// 			}
// 			// Only trigerred when this node is elected
// 			if txresp.Succeeded {
// 				glog.Infof("[node-%d]leader elected", *nodeID)
// 				SetLeader()
// 				go renewLease(cli, lease.ID)
// 				bringUpClusterIP()
// 				glog.Infof("resp: %v, %+v", txresp.Succeeded, *txresp.Header)
// 			}
// 			for _, resp := range txresp.Responses {
// 				glog.Infof("[node-%d]resp: %+v", *nodeID, resp.String())
// 			}
// 		default:
// 		}
// 	}
// }
