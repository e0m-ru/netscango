# **NetScanGo** 🌐  
*A lightweight local network scanner written in Go*  

![Go](https://img.shields.io/badge/Go-1.21+-blue?logo=go)  
![License](https://img.shields.io/badge/License-MIT-green)  

**NetScanGo** is a CLI tool for discovering active hosts and open ports in a local network. It uses ICMP ping, ARP, and TCP scanning to map devices and services. Perfect for home labs, security checks, and network diagnostics.  

---

## **Features**  
✅ **Subnet detection** – Automatically finds your local IP range.  
✅ **Fast host discovery** – Uses ICMP, TCP, and ARP (optional).  
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
netscango scan --arp --ports 22,80,443
```
**Output:**  
```
192.168.1.1   [ACTIVE]  Ports: 22 (SSH), 80 (HTTP)  
192.168.1.5   [ACTIVE]  Ports: 443 (HTTPS)  
192.168.1.102 [ACTIVE]  (ARP)  
```

### **Flags**  
| Flag       | Description                          | Example                     |
|------------|--------------------------------------|-----------------------------|
| `--icmp`   | Use ICMP ping (default)              | `--icmp`                    |
| `--arp`    | Use ARP requests (requires sudo)     | `--arp`                     |
| `--ports`  | TCP ports to scan (comma-separated)  | `--ports 80,443,8080`       |
| `--timeout`| Timeout per request (ms)             | `--timeout 500`             |

---

## **How It Works**  
1. **Subnet Detection**  
   - Gets your local IP (e.g., `192.168.1.100`) and subnet mask (`/24`).  
   - Calculates scan range (`192.168.1.1`–`192.168.1.254`).  

2. **Host Discovery**  
   - **ICMP**: Ping sweep (optional).  
   - **ARP**: Reliable for LAN (needs admin rights).  
   - **TCP**: Checks if ports (e.g., 80/443) respond.  

3. **Concurrency**  
   - Uses Go routines to scan multiple IPs/ports simultaneously.  

---

## **Roadmap**  
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
