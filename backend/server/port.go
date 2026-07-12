package server

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// freePort stops any process currently listening on the given TCP port,
// matching the behavior of start.bat so restarts do not hit bind errors.
func freePort(port string) {
	pids := listeningPIDs(port)
	self := os.Getpid()
	killed := 0

	for _, pid := range pids {
		if pid == 0 || pid == self {
			continue
		}
		log.Printf("Port %s is in use by PID %d — stopping it...", port, pid)
		if err := killPID(pid); err != nil {
			log.Printf("Warning: could not stop PID %d: %v", pid, err)
			continue
		}
		killed++
	}

	if killed > 0 {
		// Give the OS a moment to release the socket.
		time.Sleep(400 * time.Millisecond)
	}
}

func listeningPIDs(port string) []int {
	needle := ":" + port
	seen := map[int]struct{}{}
	var pids []int

	switch runtime.GOOS {
	case "windows":
		out, err := exec.Command("netstat", "-ano").Output()
		if err != nil {
			log.Printf("Warning: netstat failed while checking port %s: %v", port, err)
			return nil
		}
		for _, line := range strings.Split(string(out), "\n") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}
			// TCP  0.0.0.0:20128  0.0.0.0:0  LISTENING  12345
			if !strings.EqualFold(fields[0], "TCP") {
				continue
			}
			if !strings.Contains(fields[1], needle) {
				continue
			}
			if !strings.EqualFold(fields[3], "LISTENING") {
				continue
			}
			pid, err := strconv.Atoi(fields[len(fields)-1])
			if err != nil || pid == 0 {
				continue
			}
			if _, ok := seen[pid]; ok {
				continue
			}
			seen[pid] = struct{}{}
			pids = append(pids, pid)
		}
	default:
		out, err := exec.Command("lsof", "-ti", fmt.Sprintf("tcp:%s", port)).Output()
		if err != nil {
			return nil
		}
		for _, field := range strings.Fields(string(out)) {
			pid, err := strconv.Atoi(field)
			if err != nil || pid == 0 {
				continue
			}
			if _, ok := seen[pid]; ok {
				continue
			}
			seen[pid] = struct{}{}
			pids = append(pids, pid)
		}
	}

	return pids
}

func killPID(pid int) error {
	if runtime.GOOS == "windows" {
		return exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid)).Run()
	}
	return exec.Command("kill", "-9", strconv.Itoa(pid)).Run()
}
