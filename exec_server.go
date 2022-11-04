package main

import (
	//"bufio"
	//"bytes"
	//"encoding/json"
	"fmt"
	//"io"
	//"log"
	//"math/rand"
	//"net"
	"os"
	"os/exec"
	//"sort"
	"strconv"
	"strings"
	//"time"
)

var SDFS_PATH = "/home/uk3/CS425_MP3/files/"

//func getReplicas(fileId string) string []{
//	return {(file_id + i) % 10 + 1 for i in range(4)}
//}


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
	port string
}

type Server struct {
	fileTab SDFSFileManager
	host string
	port string
	id string
	addr Address
	lives []string

	//failureDetector FailureDetector
}

func (s Server) newFileManager () SDFSFileManager {
	fileMan := SDFSFileManager{}
	fileMan.versionMan = make(map[string]int)
	fileMan.replicaMan = make(map[string][]int)
	fileMan.idm = make(map[int]map[string]struct{})
	for i:=1; i<=10; i++ {
		fileMan.idm[i] = make(map[string]struct{})
	}
	return fileMan
}

// func init() {
// 	fileMan := SDFSFileManager()
// 	fileMan.fm = make(map[string]map[string][]interface{})
// 	fileMan.fm[sdfsFileName]["replicas"][1] = "fa22-cs425-5405.cs.illinois.edu"
// 	fileMan.fm[sdfsFileName]["version"] = 2
// }


func (s Server) getIdFromHost(host string) int {
        // """
        // e.g. fa18-cs425-g33-01.cs.illinois.edu -> 1
        // :param host: host name
        // :return: an integer id
        // """
        hostSplit := (strings.Split(host, "."))[0]
        hostIdReturn := (strings.Split(hostSplit, "-"))
        return strconv.Atoi(hostIdReturn[len(hostIdReturn) - 1])
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


func (s Server) getFile(fileManager SDFSFileManager, sdfsFileName string, localFileName string, verNo string) {
		fileTabVerMan := &s.fileTab.versionMan
		fileTabRepMan := &s.fileTab.replicaMan

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
            prefix := "uk3" + "@" + "fa22-cs425-5401.cs.illinois.edu"
            //fmt.Printf("Get file %s from chosen replica %d" % (v_file_name, from_id))
	    	fmt.Println(prefix + ":" + SDFS_PATH+localFileName)
	    	cmd := exec.Command("scp", prefix + ":" + SDFS_PATH+localFileName, "copied.txt")
	    	err := cmd.Run()
            if err != nil {
				fmt.Println("Copy Error: ", err)
            }
        } else {

            if verNo > v {
                fmt.Printf("Error: Only %d versions available, request %d.", v, verNo)
                return
            }

            fo, err := os.Create(localFileName)
            if err != nil {
            	fmt.Println("Error in Creating File: ", err)
            }

            for i := v; i >= v - verNo +1; i-- {
            	prefix := "uk3" + "@" + "fa22-cs425-5401.cs.illinois.edu"
            	vFileName := "temp.txt" + ',' + strconv.Itoa(i)
            	fmt.Println(prefix + ":" + SDFS_PATH+localFileName)
            	cmd := exec.Command("scp", prefix + ":" + SDFS_PATH+localFileName, vFileName)
            	error := cmd.Run()
            	if error != nil {
					fmt.Println("Copy Error: ", error)
           	 	}

           	 	fi, er := os.Open(vFileName)

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
	fileTabRepMan := &s.fileTab.replicaMan

	_, ok := *fileTabVerMan[sdfsFileName];
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

        _, ok := *fileTabVerMan[sdfsFileName];
        if !ok {
            fmt.Printf("Error: No such sdfs file: %s" , sdfsFileName)
            return
        }

        fmt.Printf("All the machines where %s is stored: \n", sdfsFileName)
        fmt.Println(*fileTabRepMan[sdfsFileName])

}

func (s Server) showStore () {
	hostId := s.getIdFromHost(s.host)
	//get host id and display the idm for this host id. Write the function for host id at the beginning as a helper function
	fmt.Println("All the files stored on this machine:")
	fmt.Println(s.fileTab.idm[hostId])
}


func (s Server) monitor() {
	fmt.Println("Monitor")

	for true {
		fmt.Print("-->")
		var arg string
		argSplit = strings.Split(arg, " ")
		fmt.Scanln(&arg)

		if strings.HasPrefix(arg, "get-versions") {
			if len(argSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get-versions sdfs_file_name num-versions local_file_name")
				continue
			}
			getFile()
		} else if strings.HasPrefix(arg, "get"){
			if len(argSplit) != 3 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			getFile()
		} else if strings.HasPrefix(arg, "put"){
			if len(argSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			getFile()
		} else if strings.HasPrefix(arg, "delete"){
			if len(argSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			getFile()
		} else if strings.HasPrefix(arg, "ls"){
			if len(argSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			getFile()
		} else if strings.HasPrefix(arg, "store"){
			if len(argSplit) != 4 {
				fmt.Println("Error in the arguments, Correct Format: get sdfs_file_name local_file_name")
				continue
			}
			getFile()
		} else if arg == "fm" {
            pprint(self.ft.fm)
        } else if arg == "idm" {
                fmt.Println(s.fileTab.idm)
        } else if arg == "join" {
                s.failure_detector.join()
        } else if arg == "leave" {
                s.failure_detector.leave()
        } else if arg == "ml" {
                s.failure_detector.print_ml()
        } else if arg == "lives" {
                fmt.Println(s.lives)
        } else {
            fmt.Println("[ERROR] Invalid input arg %s", arg)
        }


	}

}



func main() {
	// getFile("foo.sdfs", "copy_this.txt", "1")
	//monitor()

	fileMan := newFileManager()
	idList := []int{3,4,5,6}
	idList2  := []int{1,2,3,4}
	idList3 := []int{7,8,9,10}
	fileMan.insertFile("foo.sdfs", idList)
	fileMan.insertFile("goo.sdfs", idList2)
	fileMan.insertFile("hoo.sdfs", idList3)

	fileMan.insertFile("foo.sdfs", idList)

	fileMan.insertFile("foo.sdfs", idList)
	fileMan.deleteFile("foo.sdfs")

	fmt.Printf("foo.fdfs from chosen replica %d", fileMan.versionMan["foo.sdfs"])

}
