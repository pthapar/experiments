package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/glog"
)

const (
	RefreshInterval = time.Second * 5
	LeaderKey       = "leader"
	IfcFilePath     = "/etc/sysconfig/network-scripts/ifcfg-eth0:1"
	ifcEth01        = `
BOOTPROTO=static
DEVICE=eth0:1
ONBOOT=no
TYPE=Ethernet
USERCTL=no
NM_CONTROLLED=no
`
)

var (
	nodeID        = flag.Int("nodeID", 1, "")
	clusterIPCIDR = flag.String("cluster-ip", "", "")
	etcdAddr      = flag.String("etcd-addr", "localhost:2379", "")
	isLeader      = false
	mu            = &sync.Mutex{}
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
	// Fail and exit if this  node was a master to make sure that the
	// exit post step(systemd) for ifc clean up is executed
	// Idea is to send a signal to a channel such that the  exit go routins is kicked in
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

// TODO: return err from this method such that caller can de-elect itself as leader
func bringUpClusterIP() {
	glog.Info("[node-%d] bringing up eth0:1.", *nodeID)
	ifc, err := net.InterfaceByName("eth0")
	if err != nil {
		glog.Errorf("failed to bring up interface. %s", err.Error())
		return
	}
	_, ipNet, err := net.ParseCIDR(*clusterIPCIDR)
	if err != nil {
		glog.Errorf("failed to parse ip mask. %s", err.Error())
		return
	}

	netmask := fmt.Sprintf("%d.%d.%d.%d", ipNet.Mask[0], ipNet.Mask[1], ipNet.Mask[2], ipNet.Mask[3])
	glog.Infof("[node-%d]setting netmask: %s", *nodeID, netmask)

	addtionalAttrs := fmt.Sprintf("HWADDR=%s\nIPADDR=%s\nNETMASK=%s",
		ifc.HardwareAddr.String(), strings.Split(*clusterIPCIDR, "/")[0], netmask,
	)

	ifcContent := ifcEth01 + addtionalAttrs
	glog.Infof("writing: %s to %s", ifcContent, IfcFilePath)
	err = os.RemoveAll(IfcFilePath)
	if err != nil {
		glog.Errorf("failed to  remove ifc fil. %s", err.Error())
	}
	// Create an alias network interface
	err = ioutil.WriteFile(IfcFilePath, []byte(ifcContent), 0644)
	if err != nil {
		glog.Errorf("failed to write %s to ifc-cfg-eth0:1. %s", ifcContent, err.Error())
		return
	}

	cmd := exec.Command("/usr/sbin/ifup", "eth0:1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		glog.Errorf("failed to bring up the interface. %s", err.Error())
	}

	// tick := time.NewTicker(time.Duration(1) * time.Second)
	// ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(5)*time.Second))
	// for {
	// 	select {
	// 	case <-tick.C:
	// 		cmd := exec.Command("/usr/sbin/ifup", "eth0:1")
	// 		cmd.Stdout = os.Stdout
	// 		cmd.Stderr = os.Stderr
	// 		err = cmd.Run()
	// 		if err != nil {
	// 			glog.Errorf("failed to bring up the interface. %s", err.Error())
	// 			continue
	// 		}
	// 		return
	// 	case <-ctx.Done():
	// 		cleanUpNetIfc()
	// 		// TODO: De-elect from leadership
	// 		glog.Errorf("giving up on bringing up eth0:1")
	// 		return
	// 	}

	// }
}

// TODO: clean up netifc in case this node is not a leader. It will take care of split brain problem where
// this node was partitioned and is trying to join back.
func cleanUpNetIfc() {
	glog.Info("[node-%d] bringing down eth0:1.", *nodeID)
	cmd := exec.Command("/usr/sbin/ifdown", "eth0:1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		glog.Warningf("failed to bring down the interface. %s", err.Error())
	}
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

func giveUpLeadership(cli *clientv3.Client, leaseID clientv3.LeaseID) {
	resp, err := cli.Revoke(context.TODO(), leaseID)
	if err != nil {
		glog.Warning("[node-%d]failed to revoke lease. Deleting % key if needed", *nodeID, LeaderKey)
	}
	glog.Infof("[node-%d]resp from revoke:  %+v", *nodeID, resp)
	txn := cli.Txn(context.TODO())
	txresp, err := txn.If(
		clientv3.Compare(clientv3.Value(LeaderKey), "=", *nodeID),
	).Then(
		clientv3.OpDelete(LeaderKey),
	).Commit()
	if txresp.Succeeded {
		glog.Infof("[node-%d]de-elected myself as leader")
	}
}

func main() {
	glog.Info("connecting to etcd")
	flag.Parse()

	if *clusterIPCIDR == "" {
		glog.Exitf("cluster IP CIDR is required")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	// glog.Infof("create client")
	if err != nil {
		glog.Exitf("failed to connect to etcd servers. %s", err.Error())
	}

	// Keep the lease till this node is alive
	lease, _ := cli.Grant(context.TODO(), 2)
	glog.Infof("[node-%d]lease object: %+v", *nodeID, *lease)
	go renewLease(cli, lease.ID)
	if proposeSelfAsLeader(cli, lease.ID) {
		SetLeader()
		bringUpClusterIP()
	}

	// Make sure to bring the ifc down
	cleanUpNetIfc()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	// TODO: Add one more channel to consume exit signal from any failure in this service
	// The idea is to exit the process such that clean up can kick in.
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		cleanUpNetIfc()
		giveUpLeadership(cli, lease.ID)
		os.Exit(0)
	}()
	defer cleanUpNetIfc()
	defer giveUpLeadership(cli, lease.ID)

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
	// cli.Txn()
}
