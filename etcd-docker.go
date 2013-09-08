/**
 * Created with IntelliJ IDEA.
 * User: chris
 * Date: 08/09/2013
 * Time: 21:50
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	//"time"
	//"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
	"log"
)

func main() {

	client, err := docker.NewClient("http://localhost:4243")
	if err != nil {
		log.Fatal(err)
	}

	opts := docker.ListContainersOptions{All: true}

	containers, err := client.ListContainers(opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(containers)

}
