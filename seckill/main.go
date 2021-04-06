package main

import (
	"net/http"
	"os"
	"strings"
)

/*
	秒杀服务
*/

func main() {
	http.HandleFunc("/buy/ticket", handleReq)
	http.ListenAndServe(":3004", nil)
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	failedMsg := "handle in port:"
	writeLog(failedMsg, "./stat.log")
}

func writeLog(msg string, logPath string) {
	fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	content := strings.Join([]string{msg, "\r\n"}, "3001")
	buf := []byte(content)
	fd.Write(buf)
}
