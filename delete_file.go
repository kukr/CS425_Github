package main
import ("os/exec"
		"fmt")
	
func deleteFilefromitself(filename string){
	cmd := exec.Command("rm", "replica_2/" + filename)
	err := cmd.Run()
	if err != nil{
		fmt.Println("deleteFilefromitself error or file not present in replica_2", err)
		} else{
			//deleteFile(filename)
		}	
	cmd = exec.Command("rm", "replica_1/" + filename)
	err = cmd.Run()
	if err != nil{
	fmt.Println("deleteFilefromitself error or file not present in replica_1", err)
	} else{
		//deleteFile(filename)
	}
	cmd = exec.Command("rm", "primary/" + filename)
	err = cmd.Run()
		if err != nil{
		fmt.Println("deleteFilefromitself error or file not present in primary", err)
		} else{
			//deleteFile(filename)
		}

} 
func main(){
	deleteFilefromitself("vm1.log")
}
