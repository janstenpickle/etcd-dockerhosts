package main

import (
	"flag"
	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
	"log"
	"net"
	"time"
)

func localIp() string {
	var addr string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("Oops: " + err.Error() + "\n")
	} else {

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
				addr = ipnet.IP.String()
			}
		}
	}
	return addr
}

func inspectAndSet(c chan string, etcdClient *etcd.Client, dockerClient *docker.Client, keyname *string, ttl *uint64) {
	for {
		id := <-c
		container, err := dockerClient.InspectContainer(id)
		if err != nil {
			log.Fatal(err)
		} else {

			var addr string

			//if the container config contains ports then set the IP as the host IP
			if len(container.Config.PortSpecs) > 0 {
				addr = localIp()
			} else {
				addr = container.NetworkSettings.IPAddress
			}

			etcdClient.TestAndSet(*keyname+"/"+container.Config.Hostname, addr, addr, *ttl)
			log.Print(container.Config.Hostname)

		}

	}
}

func loop(c chan string, dockerClient *docker.Client, interval *time.Duration) {
	opts := docker.ListContainersOptions{All: false}

	for {
		containers, err := dockerClient.ListContainers(opts)
		if err != nil {
			log.Fatal(err)
		}

		for _, container := range containers {
			id := container.ID
			c <- id
		}

		time.Sleep(time.Second * *interval)
	}

}

func main() {
	keyname := flag.String("keyname", "hosts", "Etcd keyname under which to record containers' hostnames/IP")
	ttl := flag.Uint64("ttl", 172800, "Time to live of the host entry")
	dockerAPIPort := flag.String("port", "4243", "Docker API Port")
	interval := flag.Duration("interval", 10, "Docker API to Etcd sync interval")
	//etcdHost := flag.String("etcd_host", "127.0.0.1", "Etcd host")
	//etcdPort := flag.String("etcd_port", "4001", "Etcd port")
	concurrency := flag.Int("concurrency", 1, "Number of worker threads")

	flag.Parse()

	//etcdCluster := []string{"http://" + *etcdHost + ":" + *etcdPort}
	etcdClient := etcd.NewClient()
	//etcdClient.SetCluster(etcdCluster)

	dockerClient, err := docker.NewClient("http://127.0.0.1:" + *dockerAPIPort)
	if err != nil {
		log.Fatal(err)
	}

	var c = make(chan string, *concurrency)

	for i := 0; i < *concurrency; i++ {
		go inspectAndSet(c, etcdClient, dockerClient, keyname, ttl)
	}

	loop(c, dockerClient, interval)

}
