package main
import (
"net"
"encoding/gob"
"fmt"
)
var membership_list []int

type message struct {
	host int
	lives []int
}
func inform_join_node(lives int[]){
	for i:=0;i<len(membership_list);i++ {
		if membership_list[i]!= self.host{
		nodeid := membership_list[i]
		conn, err := net.Dial("tcp", getIdFromHost(nodeid)+":8081")
		if err != nil {
			log.Fatal("Connection error", err)
		}
		encoder := gob.NewEncoder(conn)
		p := &message{host, lives}
		encoder.Encode(p)
		conn.Close()
	}
	}
}
func receive_join_msg(conn net.Conn){
        /*
        Like failure detector, receiver is an UDP receiver to get all contact message including
        - UPDATE: update idm and fm, sync the status
        - DELETE: alert an deletion operation
        - FAILED_RELAY: receive from failure detector, to know that a replica is down
        - FAILED: multicasted by replicas, to let all other node know a replica is down
        - JOIN: multicasted by sdfs server whose failure detector is introducer (default is node with id 1)
        :return: None
        */
		dec := gob.NewDecoder(conn)
    	msg := &message{}
    	dec.Decode(msg)
		if failure_detector.is_introducer(){ //? could not find failure detector
			//multicast join message, sync structure of lives
			jid := getIdFromHost(msg.host)
			lives[jid] = 1
			inform_join_node(lives)
		} else{
			// if receiver is not introducer, just update lives
			for i:=0;i<len(msg.lives);i++{
				if msg.lives[i] == 1{
					lives[i] = 1
				}
			} //? Update lives
		}
	}
	func main() {
		fmt.Println("start");
	   ln, err := net.Listen("tcp", ":8081")
		if err != nil {
			// handle error
		}
		for {
			conn, err := ln.Accept() // this blocks until connection or error
			if err != nil {
				// handle error
				continue
			}
			go receive_join_msg(conn) // a goroutine handles conn so that the loop can accept other connections
		}
	}
