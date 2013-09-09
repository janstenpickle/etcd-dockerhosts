package main

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
	"log"
	"flag"
	"time"
)


func main() {
	keyname := flag.String("keyname", "hosts", "Etcd keyname under which to record containers' hostnames/IP")
	ttl := flag.Uint64("ttl", 172800, "Time to live of the host entry")
	dockerAPIPort := flag.String("port", "4243", "Docker API Port")
	interval := flag.Duration("interval", 10, "Docker API to Etcd sync interval")
	etcdHost := flag.String("etcd_host", "127.0.0.1", "Etcd host")
	etcdPort := flag.String("etcd_port", "4001", "Etcd port")

	flag.Parse()
	
	etcdCluster := []string{"http://"+*etcdHost+":"+*etcdPort}
	etcdClient := etcd.NewClient()	
	etcdClient.SetCluster(etcdCluster)
	
	client, err := docker.NewClient("http://127.0.0.1:"+*dockerAPIPort)
	if err != nil {
		log.Fatal(err)
	}

	opts := docker.ListContainersOptions{All: false}

	for true {
		containers, err := client.ListContainers(opts)
		if err != nil {
			log.Fatal(err)
		}
	
		for _, container := range containers {
			id:= container.ID
			container, err := client.InspectContainer(id)
			if err != nil {
				log.Fatal(err)
			}
			
			config := container.Config
			network := container.NetworkSettings
			
			etcdClient.TestAndSet(*keyname + "/" + config.Hostname, network.IPAddress, network.IPAddress, *ttl)
		}
		time.Sleep(time.Second**interval)
	}
}


