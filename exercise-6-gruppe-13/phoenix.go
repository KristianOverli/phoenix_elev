package main

import (
	"encoding/json"
	"net"
	"os/exec"
	"time"
	"fmt"
)


const bcPeriod = 300
const MsgMissedThereshold = 10


// Interface

//monitoringPhoenix
func main() { 
	isBackup := true

	var timeout = make(chan int)
	var lastBroadcast time.Time
	if isBackup {
		go phoenixListen(timeout)
	}
	<-timeout
	isBackup = false

	time.Sleep(time.Millisecond * 500)

	if !(isBackup) {
		spawnBackup()
		fmt.Println("spawnedBackup")
		go func() {
			for {
				if time.Now().Sub(lastBroadcast) > time.Millisecond*bcPeriod {
					lastBroadcast = time.Now()
					primaryBroadcast()
				}
			}
		}()
	}
}



func phoenixListen(timeout chan int) {
	missedMsg := 0
	LocalUdpAdder, err := net.ResolveUDPAddr("udp", ":20014")
	if err != nil {
		println(err.Error())
	}
	conn, err := net.ListenUDP("udp", LocalUdpAdder)
	if err != nil {
		println(err.Error())
	}
	defer conn.Close()

	var buf [1024]byte
	for {

		conn.SetDeadline(time.Now().Add(time.Millisecond * bcPeriod))

		_, _, err := conn.ReadFromUDP(buf[:])

		if err == nil {
			missedMsg = 0
		} else {
			missedMsg += 1
			if missedMsg >= MsgMissedThereshold{
				fmt.Println("timeout")
				timeout <- 1
				return
			}
		}

	}
}


func primaryBroadcast() {

	isAlive := true

	remoteUdpAdder, err := net.ResolveUDPAddr("udp", ":20014")
	if err != nil {
		println(err.Error())
	}
	conn, err := net.DialUDP("udp", nil, remoteUdpAdder)
	if err != nil {
		println(err.Error())
	}
	defer conn.Close()

	jsonBuf, _ := json.Marshal(isAlive)
	conn.Write(jsonBuf)
}


func spawnBackup() {

	cmd := exec.Command("gnome-terminal", "-x", "sh", "-c", "./phoenix")
	cmd.Start()
}	