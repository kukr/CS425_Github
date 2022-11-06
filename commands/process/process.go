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

	func (fileMan SDFSFileManager) insertFile(sdfsFileName string, ids []int ) {
		for _, id := range ids {
			fileMan.idm[id][sdfsFileName] = struct{}{}
		}
		_, isPresent := fileMan.versionMan[sdfsFileName]
		if !isPresent {
			fileMan.versionMan[sdfsFileName] = 0
		}
		fileMan.versionMan[sdfsFileName] += 1
		for _, id := range ids {
			fileMan.replicaMan[sdfsFileName] = append(fileMan.replicaMan[sdfsFileName], id)
		}

	}

	func (fileMan SDFSFileManager) deleteFile(sdfsFileName string) {
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

func (s Server) newFileManager () {
	s.fileTab = SDFSFileManager{}
	s.fileTab.versionMan = make(map[string]int)
	s.fileTab.replicaMan = make(map[string][]int)
	s.fileTab.idm = make(map[int]map[string]struct{})
	for i:=1; i<=10; i++ {
		s.fileTab.idm[i] = make(map[string]struct{})
	}
	
}

func (s *Server) init() {
	s.newFileManager()
	thisHostName, err := os.Hostname()

	intro.CheckErr(err)
	s.host = thisHostName
	s.port = globob.SDFS_PORT
	s.id = s.getIdFromHost(s.host)
	s.addr.host = s.host
	s.addr.port = s.port
	s.lives = []int{}

	// fileMan.fm = make(map[string]map[string][]interface{})
	// fileMan.fm[sdfsFileName]["replicas"][1] = "fa22-cs425-5405.cs.illinois.edu"
	// fileMan.fm[sdfsFileName]["version"] = 2
}

func NewServer() *Server {
	server := &Server{}
	//server.loadConfigFromJSON(jsonFile)
	server.init()
	return server
}

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

func (s Server) receiver () {
	// fileTabVerMan := &s.fileTab.versionMan
	// fileTabRepMan := &s.fileTab.replicaMan

	// indexMan := &s.fileTab.idm

	// serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.port))

	// if err != nil {
	// 	fmt.Println("err: ", err)
	// }

	// conn, err := net.ListenUDP("udp", serverAddr)

	// if err != nil {
	// 	fmt.Println("Listen err: ", err)
	// }

	// defer conn.Close()

	// var recBuf []byte
	// n, remote_addr, err := conn.ReadFromUDP(recBuf)

	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	//continue
	// }

	// buf := recBuf[:n]

	// if len(buf) == 0 {
	// 	//continue
	// }




}

func (s Server) getIdFromHost(host string) int {
        // """
        // e.g. fa18-cs425-g33-01.cs.illinois.edu -> 1
        // :param host: host name
        // :return: an integer id
        // """
        hostSplit := (strings.Split(host, "."))[0]
        hostIdReturn := (strings.Split(hostSplit, "-"))
        hostIdInt,_ := strconv.Atoi(hostIdReturn[len(hostIdReturn) - 1])
        return hostIdInt
}
  
func (s Server) getHostFromId (hostId string) string {
        // """
        // e.g. 1 -> fa18-cs425-g33-01.cs.illinois.edu
        // :param host_id: int id
        // :return: host str
        // """
        return fmt.Sprintf("fa18-cs425-g33-%02d.cs.illinois.edu", hostId)
}
// func putFile()


func (s Server) getFile(sdfsFileName string, localFileName string, verNo string) {
		fileTabVerMan := &s.fileTab.versionMan
		// fileTabRepMan := &s.fileTab.replicaMan

        _, ok := (*fileTabVerMan)[sdfsFileName];
        if !ok {
            fmt.Printf("Error: No such sdfs file: %s" , sdfsFileName)
            return
        }

        // Use it later : fromId := (*fileTabRepMan)[sdfsFileName][0]
        v := (*fileTabVerMan)[sdfsFileName]
        // to get last updated version by default (command get)
        if len(verNo) == 0 {
            version := v
            vFileName := sdfsFileName + "," + strconv.Itoa(version)
            fmt.Println(vFileName)
            prefix := "uk3" + "@" + "fa22-cs425-5401.cs.illinois.edu"
            //fmt.Printf("Get file %s from chosen replica %d" % (v_file_name, from_id))
	    	fmt.Println(prefix + ":" + globob.SDFS_PATH+localFileName)
	    	cmd := exec.Command("scp", prefix + ":" + globob.SDFS_PATH+localFileName, "copied.txt")
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
            	prefix := "uk3" + "@" + "fa22-cs425-5401.cs.illinois.edu"
            	vFileName := "temp.txt" + "," + strconv.Itoa(i)
            	fmt.Println(prefix + ":" + globob.SDFS_PATH+localFileName)
            	cmd := exec.Command("scp", prefix + ":" + globob.SDFS_PATH+localFileName, vFileName)
            	error := cmd.Run()
            	if error != nil {
					fmt.Println("Copy Error: ", error)
           	 	}

           	 	fi, er := os.Open(vFileName)

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

           	 	os.Remove(vFileName)

            }
        }
}

func (s Server) deleteFile (sdfsFileName string) {
	fileTabVerMan := &s.fileTab.versionMan
	// fileTabRepMan := &s.fileTab.replicaMan

	_, ok := (*fileTabVerMan)[sdfsFileName];
    if !ok {
        fmt.Printf("Error: No such sdfs file: %s" , sdfsFileName)
        return
    }

    s.fileTab.deleteFile(sdfsFileName)

    //TODO: Send delete message to all the servers
}

func (s Server) listSdfsFiles (sdfsFileName string) {
        fileTabVerMan := &s.fileTab.versionMan
        fileTabRepMan := &s.fileTab.replicaMan

        _, ok := (*fileTabVerMan)[sdfsFileName];
        if !ok {
            fmt.Printf("Error: No such sdfs file: %s" , sdfsFileName)
            return
        }

        fmt.Printf("All the machines where %s is stored: \n", sdfsFileName)
        fmt.Println((*fileTabRepMan)[sdfsFileName])

}

func (s Server) showStore () {
	hostId := s.getIdFromHost(s.host)
	//get host id and display the idm for this host id. Write the function for host id at the beginning as a helper function
	fmt.Println("All the files stored on this machine:")
	fmt.Println(s.fileTab.idm[hostId])
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

	ser := NewServer()
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
		var command string

		// Taking input from user
		fmt.Scanln(&command)

		commandSplit := strings.Split(command, " ")

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
			ser.getFile(commandSplit[1], commandSplit[3], "")
		} else if strings.HasPrefix(command, "put"){
			if len(commandSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			//s.putFile()
		} else if strings.HasPrefix(command, "delete"){
			if len(commandSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			ser.deleteFile(commandSplit[1])
		} else if strings.HasPrefix(command, "ls"){
			if len(commandSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			ser.listSdfsFiles(commandSplit[1])
		} else if strings.HasPrefix(command, "store"){
			if len(commandSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
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

	// Wait for the wait group to be done
	// wg.Wait()
}
