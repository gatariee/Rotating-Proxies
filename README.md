
# Rotating Proxies 
A shitty script that bypasses IP bans that are commonly done by security software (e.g Fail2Ban) utilizing rotating proxies.

## Proxy API
If you're using unreliable proxies, it's probably better for you to store a list of working proxies in a file and feeding it into the script rather than calling them from an API

However, free Proxy APIs should work fine, the one I used for this project is [here](https://geonode.com/free-proxy-list).

## Bypassing Firewalls
Security software tend to set rules in the firewall to restrict access based on IP addresses, which can make it difficult to bruteforce some services.

Rotating proxies can be used in your red team engagements to bypass these IP address restrictions and conduct bruteforce attacks more effectively.

## Considerations 
The script is designed to demonstrate the concept of using rotating proxies to bypass IP bans
However, it should be noted that the script is still relatively unstable and made for a specific use-case.




## Installation & Usage
Note that you should supply your wordlists and targets in the scripts before compiling them 
```bash
git clone https://github.com/gatariee/Rotating-Proxies
cd /Rotating-Proxies/scripts
```
GO 1.20.1
```
go build main.go -o /main
chmod +x main
./main
```
Python 3.9+
```
python -m venv env
source env/bin/activate
pip install -r requirements.txt
chmod +x main.py
./main.py
```

