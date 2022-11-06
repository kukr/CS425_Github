package main
import ("os/exec"
		"fmt"
		"hash/crc32"
		"strconv")
var membership_list = [3]string{"sandhu5@fa22-cs425-5402.cs.illinois.edu:/home/sandhu5/", "sandhu5@fa22-cs425-5402.cs.illinois.edu:/home/sandhu5/", "sandhu5@fa22-cs425-5402.cs.illinois.edu:/home/sandhu5/"}
var local_cache = make(map[string]int)
func check_in_cache(filename string) int {
	version_no, ok := local_cache[filename]
	if !ok{
		local_cache[filename] = 1
		return 1}	else{
		local_cache[filename] = version_no+1
		return version_no + 1}
}
func hashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	return v%len(membership_list)
}

func putfile(filename string){
	hash_val := hashcode(filename)
	for i:=hash_val;i<hash_val+3;i++{
		version_no := check_in_cache(filename)
		cmd := exec.Command("scp", filename , membership_list[i%(len(membership_list))] + filename + strconv.Itoa(version_no))
		err := cmd.Run()
		if err != nil{
		fmt.Println("copy error", err)}
	}

} 
func main(){
	putfile("abc.txt")
	fmt.Println(hashcode("b"))
}
