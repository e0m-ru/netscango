# **NetScanGo** 🌐  
*A lightweight local network scanner written in Go*  

![Go](https://img.shields.io/badge/Go-1.23+-blue?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)  

**NetScanGo** is a CLI tool for discovering active hosts and open ports in a local network. It uses TCP scanning to map devices and services. Perfect for home labs, security checks, and network diagnostics.  

---

## **Features**  
✅ **Subnet detection** – Automatically finds your local IP range.  
✅ **Fast host discovery** – Uses TCP.  
✅ **Port scanning** – Checks common/open ports on live hosts.  
✅ **Concurrent scanning** – Goroutines for speed.  
✅ **Simple & lightweight** – No bloat, just raw network scanning.  

---

## **Installation**  
### **From Source**  
```bash
git clone https://github.com/e0m-ru/netscango.git
cd netscango
go build -o netscango
./netscango --help
```

### **Using Go**  
```bash
go install github.com/e0m-ru/netscango@latest
```

---

## **Usage**  
### **Scan Local Network**  
```bash
netscango scan -p 22,80,443 -t 1000 -w 1000
```
**Output:**  
```
1) 192.168.1.4/24
2) 172.26.128.1/20
3) 172.17.64.1/20
Select network number: 1
192.168.1.1:80
192.168.1.40:80
```

### **Flags**  
| Flag       | Description                          | Example                     |
|------------|--------------------------------------|-----------------------------|
| `-p`       | TCP ports to scan (comma separated)  | `--ports 1-80,443,8080`       |
| `--timeout`| Timeout per request (ms)             | `--timeout 500`             |

---

## **How It Works**  
1. **Subnet Detection**  
   - Gets your local IP (e.g., `192.168.1.100`) and subnet mask (`/24`).  
   - Calculates scan range (`192.168.1.1`–`192.168.1.254`).  

2. **Host Discovery**  
   - **TCP**: Checks if ports (e.g., 80/443) respond.  

3. **Concurrency**  
   - Uses Go routines to scan multiple IPs/ports simultaneously.  

---

## **Roadmap**  
- [ ] **ICMP**: Ping sweep.  
- [ ] **ARP**: Reliable for LAN (needs admin rights).  
- [ ] MAC vendor lookup (OUI database).  
- [ ] Export results to JSON/CSV.  
- [ ] Detect OS via TTL/TCP fingerprints.  

---

## **Contributing**  
PRs welcome! For major changes, open an issue first.  

---

## **License**  
MIT © [e0m](https://github.com/e0m-ru)  

---
