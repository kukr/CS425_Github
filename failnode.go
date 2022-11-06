package main
import (
"net"
"encoding/gob"
"fmt"
"path/filepath"
)
var membership_list []int

type updatemessage struct {
	msgtype string
	file_name string
	replicas []int
	version int
}
func inform_nodes_about_fail_node(host string) {
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
func nodefailure("host" string){
	/*
	Like failure detector, receiver is an UDP receiver to get all contact message including
	- UPDATE: update idm and fm, sync the status
	- DELETE: alert an deletion operation
	- FAILED_RELAY: receive from failure detector, to know that a replica is down
	- FAILED: multicasted by replicas, to let all other node know a replica is down
	- JOIN: multicasted by sdfs server whose failure detector is introducer (default is node with id 1)
	:return: None
	*/
	fid := host
	if lives[fid] != 1{
		continue
	}
	lives[fid] = 0
	for i:=0; i<len(idm[fid]);i++{
		f:=idm[fid][i]
		replicas = fm[f]["replicas"]
		replicas[fid] = 0
		//check if itself needs to help re-replicate
		if self.id == max(replicas){
			// choice an available source
			var livescp []int
			cnt := 0
			for i:=0;i<len(lives);i++{
				if lives[i]==1 {
					if replicas[i] != 1 {
						livescp[cnt] = i
						cnt += 1
					}
				}
			}
		}
			rid = livescp[randIntn(cnt)]
			//help to re-replicate
			files, err := ioutil.ReadDir("./")
			if err != nil {
				log.Fatal(err)
			}
		
			for _, f := range files {
					file = f.Name()
					file_path := filepath.Join(SDFS_PATH, file)
					if os.Stat(file_path){ 
							if strings.HasPrefix(file, f){
								prefix := "sandhu5" + '@' + getHostFromId(rid)
								cmd := exec.Command("scp", file_name , prefix + ':' + file_path)//get_hostname_from_id
								//update status and new replica message
								replicas[rid] = 1
								inform_nodes_about_failnode("update", f, replicas, fm[f]["version"])
								}
							}
					idm[fid] = {}
}
