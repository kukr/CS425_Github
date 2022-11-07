package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"bufio"
	"os/exec"
	"strconv"
	"strings"
	"net"
	//"hash/crc32"
	//"encoding/gob"
	"hash/fnv"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"math/rand"

	process "cs425/mp/process"
	globob "cs425/mp/glob"
)

var (
	log_process_port  = flag.Int("log_process_port", 50052, "The logger process port")
	devmode           = flag.Bool("devmode", false, "Develop locally?")
	logtofile         = true
	introducerAddress = "172.22.158.181"
	introducerPort    = 50053
	udpserverport     = flag.Int("udpserverport", 20000, "Port of the UDP server")
)

type SDFSFileManager struct {
	versionMan map[string]int
	replicaMan map[string][]int
	idm map[int]map[string]struct{}
}

	func (fileMan *SDFSFileManager) insertFile(sdfsFileName string, ids []int ) {
		for _, id := range ids {
			fileMan.idm[id][sdfsFileName] = struct{}{}
		}
		_, isPresent := fileMan.versionMan[sdfsFileName]
		if !isPresent && len(ids) > 0 {
			fileMan.versionMan[sdfsFileName] = 0
		}
		//fmt.Println("Length of ids: ", len(ids))
		//fmt.Println(ids)
		if len(ids) > 0 {
			fileMan.versionMan[sdfsFileName] += 1
		}
		//fmt.Println(fileMan.versionMan[sdfsFileName])
		for _, id := range ids {
			replicaMatch := false
			for j := 0; j < len(fileMan.replicaMan[sdfsFileName]); j++ {
				if fileMan.replicaMan[sdfsFileName][j] == id {
					replicaMatch = true
					break
				}
				
			}
			if !replicaMatch {
				fileMan.replicaMan[sdfsFileName] = append(fileMan.replicaMan[sdfsFileName], id)
				fmt.Println(id)
			}
		}

	}

	func (fileMan *SDFSFileManager) deleteFile(sdfsFileName string) {
		for i:= 1; i <= 10; i++ {
			_, isPresent := fileMan.idm[i][sdfsFileName]
			if isPresent {
				delete(fileMan.idm[i], sdfsFileName)
			}
		}

		_, isPresent := fileMan.versionMan[sdfsFileName]
		if isPresent {
			delete(fileMan.versionMan, sdfsFileName)
			delete(fileMan.replicaMan, sdfsFileName)
		}
	}

type Address struct {
	host string
	port int
}

type Server struct {
	fileTab SDFSFileManager
	host string
	port int
	id int
	addr Address
	lives []int

	//failureDetector FailureDetector
}

// type updateMessage struct {
// 	Filename string
// 	replicas []int
// 	version int
// }

// func (s Server) newFileManager () {
// 	s.fileTab = SDFSFileManager{}
// 	s.fileTab.versionMan = make(map[string]int)
// 	s.fileTab.replicaMan = make(map[string][]int)
// 	s.fileTab.idm = make(map[int]map[string]struct{})
// 	for i:=1; i<=10; i++ {
// 		s.fileTab.idm[i] = make(map[string]struct{})
// 	}
	
// }

// func (s Server) init() {
// 	s.newFileManager()
// 	thisHostName, err := os.Hostname()

// 	if err!=nil {
// 		fmt.Println("Error in getting this Host Name: ", err)
// 	}
// 	s.host = thisHostName
// 	s.port = globob.SDFS_PORT
// 	s.id = s.getIdFromHost(s.host)
// 	s.addr.host = s.host
// 	s.addr.port = s.port
// 	s.lives = []int{}

// 	// fileMan.fm = make(map[string]map[string][]interface{})
// 	// fileMan.fm[sdfsFileName]["replicas"][1] = "fa22-cs425-5405.cs.illinois.edu"
// 	// fileMan.fm[sdfsFileName]["version"] = 2
// }

// func NewServer() *Server {
// 	server := &Server{}
// 	//server.loadConfigFromJSON(jsonFile)
// 	server.init()
// 	return server
// }

// func (s *Server) ListenUDP() error {
// 	/* Lets prepare a address at any address at port s.config.Port*/
// 	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.po))
// 	if err != nil {
// 		return err
// 	}

// 	/* Now listen at selected port */
// 	s.ServerConn, err = net.Listen("udp", serverAddr)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s Server) handleUpdateConnection (conn net.Conn) {

// 	fileTabVerMan := s.fileTab.versionMan
// 	fileTabRepMan := s.fileTab.replicaMan

// 	indexMan := s.fileTab.idm

// 	//fmt.Println("Handle new update message")
// 			dec := gob.NewDecoder(conn)
//    			p := &updateMessage{}
//     		dec.Decode(p)

// 			fileName := p.Filename
// 			replicaSet := p.replicas
// 			 for _, replica := range replicaSet {
// 				indexMan[replica][fileName] = struct{}{} 
// 			 }

// 			 _, isPresent := fileTabVerMan[fileName]
// 			if !isPresent {
// 				fileTabVerMan[fileName] = 1
// 			}

// 			fileTabVerMan[fileName] = p.version
			 
// 			for _, i := range replicaSet {
// 				match := false
// 				for j := 0; j < len(fileTabRepMan[fileName]); j++ {
// 					if i == fileTabRepMan[fileName][j] {
// 						match = true
// 					}
// 				}

// 				if !match {
// 					fileTabRepMan[fileName] = append(fileTabRepMan[fileName], i)
// 				}
// 			}
// 			conn.Close()

// }

// func (s Server) handleUpdate () {

// 	ln, err := net.Listen("tcp", ":53333")
	
//     if err != nil {
//         // handle error
// 		fmt.Println("Cannot start handleUpdate Server. Error: ", err)
//     }

// 	defer ln.Close()
	
// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
//             panic(err)
//         }

// 		go s.handleUpdateConnection(conn)

// 	}

func generateFailedBuffer(hostIpName string) []byte {
	var messageEncode uint8

	messageEncode = 3
	replyBuf := []byte{byte(messageEncode)}
	replyBuf = append(replyBuf, ':')
	replyBuf = append(replyBuf, []byte(hostIpName)...)
	return replyBuf
}

func (s *Server) receiver () {

	//go s.handleUpdate ()

	 fileTabVerMan := s.fileTab.versionMan
	 fileTabRepMan := s.fileTab.replicaMan

	 indexMan := s.fileTab.idm

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.port))

	if err != nil {
		fmt.Println("err: ", err)
	}

	conn, err := net.ListenUDP("udp", serverAddr)

	if err != nil {
		fmt.Println("Listen err: ", err)
	}

	defer conn.Close()

	recBuf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(recBuf)

		//fmt.Println("To check if it is an infinite loop")

		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		buf := recBuf[:n]

		if len(buf) == 0 {
			continue
		}

		bufList := bytes.Split(buf, []byte(":"))

		messageTypeInt := uint8(bufList[0][0])
		fmt.Println(messageTypeInt)

		if messageTypeInt == 1 {
			fmt.Println("[INFO] Received Update Message")
			var replicaSet []int
			fileName := string(bufList[1])

			fmt.Println(fileName)

			repliDecStr := string(bufList[2])

			replicaList := strings.Split(repliDecStr,",")
			//fmt.Println(replicaList)
			for _, replicaId := range replicaList {
				replicaIdInt,_ := strconv.Atoi(replicaId)
				replicaSet = append(replicaSet, replicaIdInt)
			}

			//fmt.Println(replicaSet)

			for _, replica := range replicaSet {
				indexMan[replica][fileName] = struct{}{} 
			 }

			 _, isPresent := fileTabVerMan[fileName]
			if !isPresent {
				fileTabVerMan[fileName] = 1
			}

			fileTabVerMan[fileName] = int(bufList[3][0])

			for _, i := range replicaSet {
				match := false
				for j := 0; j < len(fileTabRepMan[fileName]); j++ {
					if i == fileTabRepMan[fileName][j] {
						match = true
					}
				}

				if !match {
					fileTabRepMan[fileName] = append(fileTabRepMan[fileName], i)
					fmt.Println(i)
				}
				//fmt.Println(fileTabRepMan)
			}

		} else if messageTypeInt == 2 {
			fmt.Println("[INFO] Received join message.")
			livesList := bytes.Split(bufList[1], []byte(","))
			for _, liveId := range livesList {
				listMatch := false
				for _, j := range s.lives {
					if int(liveId[0]) == j {
						listMatch = true
						break
					}
				}

				if !listMatch {
					s.lives = append(s.lives, int(liveId[0]))
				}
			}

		} else if messageTypeInt == 3 {
			fmt.Println("[INFO] Failed Message")
			failIdName := string(bufList[1])
			failId := getIdFromHost(failIdName)

			failIdexist := false
			for it := 0; it < len(s.lives); it++ {

				if failId == s.lives[it] {
					failIdexist = true
					break
				}

			}
			if !failIdexist {
				continue
			}
			for _, curHostName := range globob.ALL_HOSTS {
				if curHostName != s.host {
					currAddrwithPort := fmt.Sprintf("%s:%d", curHostName, globob.SDFS_PORT)
					conn, err := net.Dial("udp", currAddrwithPort)

					if err != nil {
						log.Fatal("Connection error", err)
					}

					failedMsgSend := generateFailedBuffer(failIdName)

					_, err = conn.Write(failedMsgSend)
		
					if err != nil {
						log.Println("update States 2nd error: ", err)
					}
		
					conn.Close()
				}
			}

			var newLives []int

			for it := 0; it < len(s.lives); it++ {
				if failId == s.lives[it] {
					continue
				} else {
					newLives = append(newLives, s.lives[it])
				}
			}

			s.lives = newLives

			for key, _ := range indexMan[failId] {

				var replicas []int

				for it := 0; it < len(fileTabRepMan[key]); it++ {
					if failId == fileTabRepMan[key][it] {
						continue
					} else {
						replicas = append(replicas, fileTabRepMan[key][it])
					}
				}
				fileTabRepMan[key] = replicas

				maxReplica := replicas[0]

				for _, replica:= range replicas {
					if replica > maxReplica {
						maxReplica = replica
					}
				}

				if s.id == maxReplica {
					var availSour []int

					for it := 0; it < len(s.lives); it++ {
						replicaMat := false
						for aj := 0; aj < len(replicas); aj++ {
							if s.lives[it] == replicas[aj] {
								replicaMat = true
								break
							}
						}
						if !replicaMat {
							availSour = append(availSour, s.lives[it])
						}
					}

					randomIndex := rand.Intn(len(availSour))
					rid := availSour[randomIndex]

					files, err := ioutil.ReadDir(globob.SDFS_PATH)
    				if err != nil {
        				log.Fatal("Read Directory Failed: ", err)
    				}

					for _, f := range files {
						filePath := filepath.Join(globob.SDFS_PATH, f.Name())
				
						if !f.IsDir() && strings.HasPrefix(f.Name(), key) {
							fmt.Println("[INFO] Re-replica file %s to %d", key, rid)
							rHost := s.getHostFromId(rid)
							prefix := "uk3" + "@" + rHost
							cmd := exec.Command("scp", filePath , prefix + ":" + filePath)//get_hostname_from_id
							err := cmd.Run()

							if err != nil{
								fmt.Println("Re-replicate error: ", err)
							}
						}
					}

					newReplicas := append(fileTabRepMan[key], rid)

					updateMsg := generateUpdateBuffer(key, newReplicas, fileTabVerMan[key])

					for _, curHostName := range globob.ALL_HOSTS {
						hostAlive := false
						for i := 0; i < len(s.lives); i++ {
							if getIdFromHost(curHostName) == s.lives[i] {
								hostAlive = true
								break
							}
						}
						if hostAlive {
							fmt.Println("CurrentHostName: ", curHostName)
							thisHostAddrwithPort := fmt.Sprintf("%s:%d", curHostName, globob.SDFS_PORT)
							conne, err := net.Dial("udp", thisHostAddrwithPort)
							if err != nil {
								log.Fatal("Connection error", err)
							}
				
							_, err = conne.Write(updateMsg)
				
							if err != nil {
								log.Println("update States 2nd error: ", err)
							}
				
							conne.Close()
						
						}
					}

				}



			}

			for key, _ := range indexMan[failId] {
				delete(indexMan[failId], key)
			}

		} else if messageTypeInt == 4 {
			fmt.Println("[INFO] Receive Delete Message")
			fileName := string(bufList[1])
			s.fileTab.deleteFile(fileName)

			files, err := ioutil.ReadDir(globob.SDFS_PATH)
    		if err != nil {
        		log.Fatal(err)
    		}

			for _, f := range files {
				filePath := filepath.Join(globob.SDFS_PATH, f.Name())
				
				if !f.IsDir() && strings.HasPrefix(f.Name(), fileName) {
					os.Remove(filePath)
				}
			}

		}

	}

}

func getIdFromHost(host string) int {
        // """
        // e.g. fa18-cs425-g33-01.cs.illinois.edu -> 1
        // :param host: host name
        // :return: an integer id
        // """
        hostSplit := (strings.Split(host, "."))[0]
        hostIdReturn := (strings.Split(hostSplit, "-"))
        hostIdInt,_ := strconv.Atoi(hostIdReturn[len(hostIdReturn) - 1])
        hostIdInt = hostIdInt%10
		if hostIdInt == 0 {
			return 10
		} else {
        	return hostIdInt
		}
}
  
func (s *Server) getHostFromId (hostId int) string {
        // """
        // e.g. 1 -> fa18-cs425-g33-01.cs.illinois.edu
        // :param host_id: int id
        // :return: host str
        // """
        return fmt.Sprintf("fa22-cs425-54%02d.cs.illinois.edu", hostId)
}

func hashCode(s string) int {

	h := fnv.New32a()
    h.Write([]byte(s))
	md5 := int(h.Sum32())
	sum := 0
	for i :=0; i<10; i++ {
		sum = sum + md5%10
		md5 = md5/10
	}

	return sum%10 + 1
}

// func (s Server) informNodesInsertIntoSdfs(sdfsFileName string, targetIds []int, verNo int){
// 	p := &updateMessage{sdfsFileName,targetIds, verNo}
// 	for _, curHostName := range globob.ALL_HOSTS {
// 		hostAlive := false
// 		for i := 0; i < len(s.lives); i++ {
// 			if getIdFromHost(curHostName) == s.lives[i] {
// 				hostAlive = true
// 				break
// 			}
// 		}
// 		if hostAlive {
// 			fmt.Println("CurrentHostName: ", curHostName)
// 		conn, err := net.Dial("tcp", curHostName+":53333")
// 		if err != nil {
// 			log.Fatal("Connection error", err)
// 		}
// 		encoder := gob.NewEncoder(conn)
// 		fmt.Println("Sending update message")
// 		encoder.Encode(p)
// 		conn.Close()
// 		}
// 	}
// }



func (s *Server) getDefaultReplicas (pid int) []int {
	var setReplicas []int
	count := 1
	for i := 0; i < 10; i++ {
		mod := ((pid+i)%10) + 1
		aliveMatch := false
		for j := 0; j < len(s.lives); j++ {
			if mod == s.lives[j] {
				aliveMatch = true
				break
			}
		}
		if aliveMatch && count <= 4 {
			count += 1
			setReplicas = append(setReplicas, mod)
		}
	}
	return setReplicas
}

func generateUpdateBuffer(sdfsFileName string, targetIds []int, verNo int) []byte {
	
	var messageEncode uint8

	messageEncode = 1
	replyBuf := []byte{byte(messageEncode)}
	replyBuf = append(replyBuf, ':')
	replyBuf = append(replyBuf, []byte(sdfsFileName)...)
	replyBuf = append(replyBuf, ':')
	replicaString := ""
	for i := 0; i < len(targetIds) - 1; i++ {
		replicaString = replicaString + strconv.Itoa(targetIds[i])
		replicaString = replicaString + "," 
	}
	replicaString = replicaString + strconv.Itoa(targetIds[len(targetIds) - 1])
	fmt.Println(replicaString)
	replyBuf = append(replyBuf, []byte(replicaString)...)
	replyBuf = append(replyBuf, ':')
	replyBuf = append(replyBuf, byte(verNo))
	
	return replyBuf
}

func (s *Server) putFile(localFileName string, sdfsFileName string){
	
	if _, err := os.Stat(globob.SDFS_PATH); os.IsNotExist(err) {
		fmt.Println("[ERROR] No such local file: %s", localFileName)
		return
	}

	fileTabVerMan := s.fileTab.versionMan
	fileTabRepMan := s.fileTab.replicaMan

	hashVal := hashCode(localFileName)
	_, ok := fileTabVerMan[sdfsFileName]

	var targetIds []int
	var versionNo int
    if ok {
		fmt.Println("Came Here")
		targetIds = fileTabRepMan[sdfsFileName]
		versionNo = fileTabVerMan[sdfsFileName]
		fmt.Println("Version No: ", versionNo)
	} else {
		versionNo = 0;
		targetIds = s.getDefaultReplicas(hashVal)
		fmt.Println("Target Ids: ", targetIds)
	}

	fmt.Println("Target ids: ", targetIds)
	versionNo += 1
	vFileName := sdfsFileName + "," + strconv.Itoa(versionNo)
	fmt.Printf("[INFO] Put file %s to %s \n", localFileName, vFileName)
	var targetActual []int
	for _, i := range targetIds {
		targetHost := s.getHostFromId(i)
		prefix := "uk3" + "@" + targetHost
		fmt.Println("Target Host: ", targetHost)
		cmd := exec.Command("scp", localFileName , prefix + ":" + filepath.Join(globob.SDFS_PATH, vFileName))//get_hostname_from_id
		err := cmd.Run()
		if err != nil{
			fmt.Println("putfile copy error: ", err, i)
		} else{
			targetActual = append(targetActual, i)
		}
	}
	fmt.Println("Target Actual: ", targetActual)
	(&(s.fileTab)).insertFile(sdfsFileName,targetActual)

	updateMsg := generateUpdateBuffer(sdfsFileName, targetActual, fileTabVerMan[sdfsFileName])

	for _, curHostName := range globob.ALL_HOSTS {
		hostAlive := false
		for i := 0; i < len(s.lives); i++ {
			if getIdFromHost(curHostName) == s.lives[i] {
				hostAlive = true
				break
			}
		}
		if hostAlive {
			// fmt.Println("CurrentHostName: ", curHostName)
			thisHostAddrwithPort := fmt.Sprintf("%s:%d", curHostName, globob.SDFS_PORT)
			conn, err := net.Dial("udp", thisHostAddrwithPort)
			if err != nil {
				log.Fatal("Connection error", err)
			}

			_, err = conn.Write(updateMsg)

			if err != nil {
				log.Println("update States 2nd error: ", err)
			}

			conn.Close()
		
		}
	}

	}



func (s *Server) getFile(sdfsFileName string, localFileName string, verNo string) {
	fileTabVerMan := s.fileTab.versionMan
	fileTabRepMan := s.fileTab.replicaMan

	_, ok := fileTabVerMan[sdfsFileName];
        if !ok {
            fmt.Printf("Error: No such sdfs file: %s" , sdfsFileName)
            return
        }

        fromId := fileTabRepMan[sdfsFileName][0]
        v := fileTabVerMan[sdfsFileName]
        // to get last updated version by default (command get)
        if len(verNo) == 0 {
            version := v
            vFileName := sdfsFileName + "," + strconv.Itoa(version)
            fmt.Println(vFileName)
            prefix := "uk3" + "@" + s.getHostFromId(fromId)
            fmt.Printf("Get file %s from chosen replica %d" , vFileName, fromId)
	    	//fmt.Println(prefix + ":" + globob.SDFS_PATH+localFileName)
	    	cmd := exec.Command("scp", prefix + ":" + globob.SDFS_PATH+vFileName, localFileName)
	    	err := cmd.Run()
            if err != nil {
				fmt.Println("Copy Error: ", err)
            }
        } else {
        	verNoInt, _ := strconv.Atoi(verNo)

            if verNoInt > v {
                fmt.Printf("Error: Only %d versions available, request %d.", v, verNo)
                return
            }

            fo, err := os.Create(localFileName)
            if err != nil {
            	fmt.Println("Error in Creating File: ", err)
            }

            for i := v; i >= v - verNoInt +1; i-- {
            	prefix := "uk3" + "@" + s.getHostFromId(fromId)
				tempVFileName := "temp.txt" + "," + strconv.Itoa(i)
            	vFileName := sdfsFileName + "," + strconv.Itoa(i)
            	fmt.Println(prefix + ":" + globob.SDFS_PATH+localFileName)
            	cmd := exec.Command("scp", prefix + ":" + filepath.Join(globob.SDFS_PATH, vFileName), tempVFileName)
            	error := cmd.Run()
            	if error != nil {
					fmt.Println("Copy Error: ", error)
           	 	}

           	 	fi, er := os.Open(tempVFileName)

           	 	if er != nil {
        			panic(er)
    			}

           	 	w := bufio.NewWriter(fo)
           	 	scanner := bufio.NewScanner(fi)

           	 	var line string

           	 	for scanner.Scan() {
           	 		line = scanner.Text()
           	 		fmt.Fprintln(w, line)
           	 	}

           	 	w.Flush()

           	 	fo.Close()
           	 	fi.Close()

           	 	os.Remove(tempVFileName)

            }
        }
}

func (s *Server) deleteFile (sdfsFileName string) {
	fileTabVerMan := s.fileTab.versionMan
	// fileTabRepMan := &s.fileTab.replicaMan

	_, ok := fileTabVerMan[sdfsFileName];
    if !ok {
        fmt.Printf("Error: No such sdfs file: %s" , sdfsFileName)
        return
    }

    s.fileTab.deleteFile(sdfsFileName)

    var messageEncode uint8

	messageEncode = 4
	replyBuf := []byte{byte(messageEncode)}
	replyBuf = append(replyBuf, ':')
		
	replyBuf = append(replyBuf, []byte(sdfsFileName)...)

	for _, curHostName := range globob.ALL_HOSTS {
		currAddrwithPort := fmt.Sprintf("%s:%d", curHostName, globob.SDFS_PORT)
		conn, err := net.Dial("udp", currAddrwithPort)

		if err != nil {
			log.Fatal("Connection error", err)
		}

		defer conn.Close()
		

		_, err = conn.Write(replyBuf)

		if err != nil {
			log.Fatal("Connection write error", err)
		}

	}

}

func (s *Server) listSdfsFiles (sdfsFileName string) {
		fileTabVerMan := s.fileTab.versionMan
		fileTabRepMan := s.fileTab.replicaMan

		_, ok := fileTabVerMan[sdfsFileName]
        if !ok {
            fmt.Printf("Error: No such sdfs file: %s" , sdfsFileName)
            return
        }

        fmt.Printf("All the machines where %s is stored: \n", sdfsFileName)
        fmt.Println(fileTabRepMan[sdfsFileName])

}

func (s *Server) showStore () {
	hostId := getIdFromHost(s.host)
	//get host id and display the idm for this host id. Write the function for host id at the beginning as a helper function
	fmt.Println("All the files stored on this machine:")
	for key, _ := range s.fileTab.idm[hostId] {
		fmt.Print(key, " ")
	}
	fmt.Print("\n")
}


// func (s Server) monitor() {
// 	fmt.Println("Monitor")

// 	for true {
// 		fmt.Print("-->")
// 		var arg string
// 		fmt.Scanln(&arg)

// 		commandSplit := strings.Split(arg, " ")

// 		if strings.HasPrefix(arg, "get-versions") {
// 			if len(commandSplit) != 4 {
// 				fmt.Println("Error in the arguments, Correct Format: get-versions sdfs_file_name num-versions local_file_name")
// 				continue
// 			}
			
// 			s.getFile(commandSplit[1], commandSplit[3], commandSplit[2])
// 		} else if strings.HasPrefix(arg, "get"){
// 			if len(commandSplit) != 3 {
// 				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
// 				continue
// 			}
// 			s.getFile(commandSplit[1], commandSplit[3], "")
// 		} else if strings.HasPrefix(arg, "put"){
// 			if len(commandSplit) != 4 {
// 				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
// 				continue
// 			}
// 			//s.putFile()
// 		} else if strings.HasPrefix(arg, "delete"){
// 			if len(commandSplit) != 4 {
// 				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
// 				continue
// 			}
// 			s.deleteFile(commandSplit[1])
// 		} else if strings.HasPrefix(arg, "ls"){
// 			if len(commandSplit) != 4 {
// 				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
// 				continue
// 			}
// 			s.listSdfsFiles(commandSplit[1])
// 		} else if strings.HasPrefix(arg, "store"){
// 			if len(commandSplit) != 4 {
// 				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
// 				continue
// 			}
// 			s.showStore()
// 		} else if arg == "fm" {
//                 fmt.Println(s.fileTab.replicaMan)
//                 fmt.Println(s.fileTab.versionMan)
//         } else if arg == "idm" {
//                 fmt.Println(s.fileTab.idm)
//         } else if arg == "join" {
//                 //s.failure_detector.join()
//         } else if arg == "leave" {
//                 //s.failure_detector.leave()
//         } else if arg == "ml" {
//                 //s.failure_detector.print_ml()
//         } else if arg == "lives" {
//                 fmt.Println(s.lives)
//         } else {
//             fmt.Println("[ERROR] Invalid input arg %s", arg)
//         }


// 	}

// }

func (s *Server) Run() {
	if _, err := os.Stat(globob.SDFS_PATH); !os.IsNotExist(err) {
		os.RemoveAll(globob.SDFS_PATH)
	}
	os.Mkdir(globob.SDFS_PATH, 0700)

	//go s.monitor()

	go s.receiver()

}

func main() {
	port := flag.Int("port", 50054, "The failure detector process port")
	flag.Parse()
	log.Printf("port: %v", *port)
	wg := new(sync.WaitGroup)
	wg.Add(4)
	if logtofile {
		// write logs of the service process to process.log file
		f, err := os.OpenFile(fmt.Sprintf("process-%v.log", *port), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	if *devmode {
		introducerAddress = "localhost"
	}

	introAddr := fmt.Sprintf("%s:%d", introducerAddress, introducerPort)

	// Start the process
	process.Run(*port, *udpserverport, *log_process_port, wg, introAddr, *devmode)

	ser := &Server{}
	ser.fileTab = SDFSFileManager{}
	ser.fileTab.versionMan = make(map[string]int)
	ser.fileTab.replicaMan = make(map[string][]int)
	ser.fileTab.idm = make(map[int]map[string]struct{})
	// fmt.Println("New File Manager")
	for i:=1; i<=10; i++ {
		ser.fileTab.idm[i] = make(map[string]struct{})
	}
	thisHostName, err := os.Hostname()
	//fmt.Println(thisHostName)

	if err != nil {
		fmt.Println("Main Err: ", err)
	}
	ser.host = thisHostName
	ser.port = globob.SDFS_PORT
	ser.id = getIdFromHost(ser.host)
	ser.addr.host = ser.host
	ser.addr.port = ser.port
	ser.lives = []int{}
	ser.lives = append(ser.lives, ser.id)
	//fmt.Println("Server lives: ", ser.lives)
	//ser.init()
	ser.Run()

	for {

		helperString := "\n======  Command List  ====== \n\t" +
        "- get [sdfs_file_name] [local_file_name] \n\t" +
        "- get-versions [sdfs_file_name] [num-versions] [local_file_name] \n\t" +
        "- put [local_file_name] [sdfs_file_name] \n\t" +
        "- delete [sdfs_file_name] \n\t" +
        "- ls [sdfs_file_name] \n\t" +
        "- store \n\t" +
		"- lives \n\t" +
        "- printmembershiplist \n\t" +
        "- printtopology \n\t" +
		"- leave \n\t" +
		"- exit (To exit) \n\t" +
        "============================ \n\n\t: "
		
		fmt.Printf(helperString)

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
    		command := scanner.Text()
    		//fmt.Printf("Input was: %q\n", command)


		// Taking input from user
		//fmt.Scanln(&command)
		//fmt.Println(command)

		commandSplit := strings.Split(command, " ")
		//fmt.Println(commandSplit)

		if strings.HasPrefix(command, "get-versions") {
			if len(commandSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get-versions sdfs_file_name num-versions local_file_name")
				continue
			}
			
			ser.getFile(commandSplit[1], commandSplit[3], commandSplit[2])
		} else if strings.HasPrefix(command, "get"){
			if len(commandSplit) != 3 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			ser.getFile(commandSplit[1], commandSplit[2], "")
		} else if strings.HasPrefix(command, "put"){
			if len(commandSplit) != 3 {
				fmt.Println("Error in the arguments, Correct Format: put sdfs_file_name local_file_name")
				continue
			}
			ser.putFile(commandSplit[1], commandSplit[2])
		} else if strings.HasPrefix(command, "delete"){
			if len(commandSplit) != 2 {
				fmt.Println("Error in the arguments, Correct Format: delete sdfs_file_name")
				continue
			}
			ser.deleteFile(commandSplit[1])
		} else if strings.HasPrefix(command, "ls"){
			if len(commandSplit) != 2 {
				fmt.Println("Error in the arguments, Correct Format: ls sdfs_file_name")
				continue
			}
			ser.listSdfsFiles(commandSplit[1])
		} else if strings.HasPrefix(command, "store"){
			if len(commandSplit) != 1 {
				fmt.Println("Error in the arguments, Correct Format: store")
				continue
			}
			ser.showStore()
		} else if command == "fm" {
                fmt.Println(ser.fileTab.replicaMan)
                fmt.Println(ser.fileTab.versionMan)
        } else if command == "idm" {
                fmt.Println(ser.fileTab.idm)
        } else if command == "printmembershiplist" {
				fmt.Println(process.GetMemberList().GetList())
        } else if command == "lives" {
                fmt.Println(ser.lives)
        } else if command == "printtopology" {
				fmt.Println(process.GetNetworkTopology())
		} else if command == "leave" {
				process.LeaveNetwork()
		} else if command == "exit" {
			os.Exit(3)
		} else {
            fmt.Println("[ERROR] Invalid input arg %s", command)
        }

		}
	}

	// Wait for the wait group to be done
	// wg.Wait()
}
