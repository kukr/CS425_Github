package main
import ("os/exec"
		"fmt"
		"hash/crc32"
		"strconv")
func hashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	return v%len(membership_list)
}
func inform_nodes(sdfsFileName string, nodeid int){
	for(i=0;i<len(membership_list);i++){
		conn, err := net.Dial("udp", getHostFromId(i))
		if err != nil {
			fmt.Printf("Some error %v", err)
		}
		fmt.Fprintf(conn, "written "+ sdfsFileName + "to" + getHostFromId(nodeid))
	}
}
func putfile(sdsFilename string, localfilename string){
	if 
	hash_val := hashcode(filename)
	version_no, ok := *fileTabVerMan[sdfsFileName];
    if !ok {
        version_no = 0;
	   }
	version_no += 1
	cmd := exec.Command("scp", localfilename , getHostFromId(hash_val%(len(membership_list))) + "primary/"+ sdfsFilename + strconv.Itoa(version_no))//get_hostname_from_id
	err := cmd.Run()
	if err != nil{
	fmt.Println("copy error", err)
}
else{
insertfile(sdfsFileName,[hash_val])
inform_nodes(sdfsFileName, hash_val)
}
cmd = exec.Command("scp", localfilename , getHostFromId((hash_val+1)%(len(membership_list))) +"replica_1/" + sdsFilename + strconv.Itoa(version_no))//get_hostname_from_id
err = cmd.Run()
if err != nil{
fmt.Println("copy error", err)
}
else{
	insertfile(sdfsFileName,[hash_val+1])
	inform_nodes(sdfsFileName, hash_val+1)
}
cmd = exec.Command("scp", localfilename , getHostFromId((hash_val+2)%(len(membership_list))) +"replica_2/" + sdfsFilename + strconv.Itoa(version_no))//get_hostname_from_id
err = cmd.Run()
if err != nil{
fmt.Println("copy error", err)
}
else{
insertfile(sdfsFileName,[hash_val+2])
inform_nodes(sdfsFileName, hash_val+2)
}
}
func main(){
	putfile("abc.txt","abc.txt")
	fmt.Println(hashcode("b"))
}
