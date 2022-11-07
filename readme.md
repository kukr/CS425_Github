# SDFS File Manager #
Our SDFS design follows a peer-to-peer architecture, where each node can be queried for a file and requested for storing the file. The main idea of the SDFS is that each node contains file which contains meta-information about storage related to each file being stored in the SDFS:-  filename, versions, primary node and replica nodes. The consistency of this file is ensured by ONE-to-ALL message sharing about each update over the network.

# Usage #

On each of the VMs run the following command
`go run commands/process/process.go`

To read a file to SDFS

`get sdfsfilename filename`

To write a file to SDFS

`put filename sdfsfilename`

To delete file

`delete filename`

To list all machine (VM) addresses where this file is currently being stored

`ls filename`

To list all files currently being stored at this machine

`store`
