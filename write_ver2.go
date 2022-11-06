package main
import ("os/exec"
		"fmt"
		"hash/crc32"
		"strconv")
membership_list := [] //Need a genuine memebershiplist
num_replicas := 0
func hashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	return v%len(membership_list)
}
func inform_nodes_insert_into_sdfds(sdfsFileName string, nodeid int){
	for i=0;i<len(membership_list);i++ {
		conn, err := net.Dial("tcp", get_hostname_from_id(nodeid)+":8080")
		if err != nil {
			log.Fatal("Connection error", err)
		}
		encoder := gob.NewEncoder(conn)
		p := &message{0,sdfsFileName,nodeid}
		encoder.Encode(p)
		conn.Close()
	}
}
func putfile(sdsFilename string, localfilename string){
	hash_val := hashcode(filename)
	version_no, ok := *fileTabVerMan[sdfsFileName];
    if !ok {
        version_no = 0;
	   }
	version_no += 1
	for (i=0; i<num_replicas;i++){ //set num_replicas
	cmd := exec.Command("scp", localfilename , getHostFromId((hash_val+i)%(len(membership_list))) + sdfsFilename + strconv.Itoa(version_no))//get_hostname_from_id
	err := cmd.Run()
	if err != nil{
	fmt.Println("putfile copy error" + getHostFromId((hash_val+i)%(len(membership_list))), err)
}
else{
insertfile(sdfsFileName,[hash_val+i])
inform_nodes_insert_into_sdfds(sdfsFileName, hash_val+i)
}
}
}
func main(){
	putfile("abc.txt","abc.txt")
	fmt.Println(hashcode("b"))
}

package main
import ("os/exec"
		"fmt"
		"hash/crc32"
		"strconv")
membership_list := [] //Need a genuine memebershiplist
num_replicas := 0
func hashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	return v%len(membership_list)
}
func inform_nodes_insert_into_sdfds(sdfsFileName string, nodeid int){
	for i=0;i<len(membership_list);i++ {
		conn, err := net.Dial("tcp", get_hostname_from_id(nodeid)+":8080")
		if err != nil {
			log.Fatal("Connection error", err)
		}
		encoder := gob.NewEncoder(conn)
		p := &message{0,sdfsFileName,nodeid}
		encoder.Encode(p)
		conn.Close()
	}
}
func putfile(sdsFilename string, localfilename string){
	hash_val := hashcode(filename)
	version_no, ok := *fileTabVerMan[sdfsFileName];
    if !ok {
        version_no = 0;
	   }
	version_no += 1
	for (i=0; i<num_replicas;i++){ //set num_replicas
	cmd := exec.Command("scp", localfilename , getHostFromId((hash_val+i)%(len(membership_list))) + sdfsFilename + strconv.Itoa(version_no))//get_hostname_from_id
	err := cmd.Run()
	if err != nil{
	fmt.Println("putfile copy error" + getHostFromId((hash_val+i)%(len(membership_list))), err)
}
else{
insertfile(sdfsFileName,[hash_val+i])
inform_nodes_insert_into_sdfds(sdfsFileName, hash_val+i)
}
}
}
func main(){
	putfile("abc.txt","abc.txt")
	fmt.Println(hashcode("b"))
}
