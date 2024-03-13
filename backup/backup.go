package backup

import (
	"fmt"
	"net"
	"os/exec"
	"time"
	"path/filepath"
	"os"
)


func Backmyshiup(Addr1 string, Addr2 string, id string ) {
	fmt.Println("initilaizing backup")
	udpAddr, err := net.ResolveUDPAddr("udp", Addr1)
    if err != nil {
        fmt.Println("Error resolving UDP address:", err)
        return
    }
    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        fmt.Println("Error listening on UDP:", err)
        return
    }
	buffer := make([]byte, 1024)
	defer conn.Close()

	for {
		conn.SetReadDeadline(time.Now().Add(4 * time.Second))
		_, _, err := conn.ReadFromUDP(buffer)
		i := buffer[0]
		if err != nil {
			go Process(i, Addr1, Addr2, id)
			break
		}
	}
}


func Process(Count byte, Addr1 string, Addr2 string, id string) {
	fmt.Println("I am the captain now")

	dir, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    // Construct the path to the file
    filename := filepath.Join(dir, "main")

    // Run the command
	//cmdStr := `tell app "Terminal" to do script "cd ` + dir + ` && go run ` + filename + `.go -port=` + Addr2 + ` -id=` + id + `"`
    //cmd := exec.Command("osascript", "-e", cmdStr)
    //err = cmd.Run()
    //if err != nil {
    //    fmt.Println("Command finished with error: ", err)
    //}

	cmdStr := "cd " + dir + " && go run " + filename + ".go -port=" + Addr2 + " -id=" + id
	cmd := exec.Command("gnome-terminal", "-x", "sh", "-c", cmdStr)
	err = cmd.Run()
	if err != nil {
    fmt.Println("Command finished with error: ", err)
	}
	
	udpAddr, _ := net.ResolveUDPAddr("udp", Addr1)
	conn, _ := net.DialUDP("udp",nil, udpAddr)
	defer conn.Close()
	
	for {
	time.Sleep(1 * time.Second)
		Count++
		fmt.Println("Count:", Count)
		conn.Write([]byte{Count}) //write to server
	}
	
}
