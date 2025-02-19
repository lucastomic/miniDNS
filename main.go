package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/miekg/dns"
)

var forbidden = map[string]any{
	"example.com.": 1,
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	for _, q := range r.Question {
		host := q.Name
		fmt.Println("[DNS] Received query for:", host)
		domainAndTLD := strings.Join(strings.Split(host, ".")[1:], ".")
		if _, found := forbidden[domainAndTLD]; found {
			fmt.Println("[DNS] Refusing connection")
			m.Rcode = dns.RcodeRefused
			w.WriteMsg(m)
			return
		} else {
			fmt.Println("[DNS] Forwarding query to Google DNS")
			resp, err := forwardDNSQuery(r)
			if err != nil {
				log.Println("[DNS] Error forwarding query:", err)
			} else {
				m = resp
			}
		}
	}

	w.WriteMsg(m)
}

func forwardDNSQuery(r *dns.Msg) (*dns.Msg, error) {
	c := new(dns.Client)
	resp, _, err := c.Exchange(r, "8.8.8.8:53")
	return resp, err
}

func startDNSServer(done chan struct{}) {
	dns.HandleFunc(".", handleDNSRequest)

	server := &dns.Server{
		Addr: ":53",
		Net:  "udp",
	}

	fmt.Println("[DNS] Server starting on port 53...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("[DNS] Failed to start: %v", err)
	}
	close(done)
}

func backupAndModifyDNSSettings() (string, error) {
	out, err := exec.Command("networksetup", "-getdnsservers", "Wi-Fi").Output()
	if err != nil {
		return "", err
	}
	backup := string(out)

	err = exec.Command("networksetup", "-setdnsservers", "Wi-Fi", "127.0.0.1").Run()
	if err != nil {
		return "", err
	}

	fmt.Println("[CONFIG] DNS settings modified to use local DNS (127.0.0.1)")
	return backup, nil
}

func restoreDNSSettings(backup string) {
	err := exec.Command("networksetup", "-setdnsservers", "Wi-Fi", backup).Run()
	if err != nil {
		log.Println("[CONFIG] Error restoring DNS settings:", err)
	} else {
		fmt.Println("[CONFIG] DNS settings restored to original settings.")
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("[CONFIG] No sites provided. Usage: %s <site1> <site2> ...", "miniDNS")
	}

	sites := os.Args[1:]
	for _, site := range sites {
		forbidden[site+"."] = 1
	}

	originalResolConf, err := backupAndModifyDNSSettings()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to modify DNS settings: %v", err)
	}
	defer restoreDNSSettings(originalResolConf)

	dnsDone := make(chan struct{})

	go startDNSServer(dnsDone)
	time.Sleep(100 * time.Millisecond)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Printf("\n[MAIN] Received signal: %v. Shutting down...\n", sig)

	time.Sleep(200 * time.Millisecond)
}
