# ☕ Decafein  
**Lightweight Time-Management & User Session Control — written in Go**

![Go](https://img.shields.io/badge/Go-1.23-blue?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-Apache--2.0-green?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Alpha-orange?style=for-the-badge)
![Platforms](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-purple?style=for-the-badge)

Decafein is a small, efficient, and modular **time-management system** that tracks application usage, manages per-user sessions, and optionally blocks websites for productivity or parental control.

---

## 🚀 Features

### 🖥️ Application Usage Tracking
- Monitors running applications  
- Records active time per window or process  
- Logging designed for analytics and auditing  

### 👥 Per-User Sessions
- Perfect for shared computers  
- Separate rules per user  
- Good for schools, cafes, and families  

### ⛔ Website Blocking *(optional)*
- Productivity mode  
- Parental control mode  
- Custom blacklist / whitelist  

### ⚙️ Lightweight by Design
- Written in **Go**  
- Low resource usage  
- Minimal external dependencies  

---

## 🧭 Use Cases

- 🏫 Schools & universities  
- 🎮 Cybercafés / LAN houses  
- 🏠 Shared family computers  
- 👨‍👩‍👧 Parental screen-time control  
- ⏱️ Personal productivity setups  

---

## 📦 Installation

```bash
git clone https://github.com/davezant/decafein
cd decafein
go build ./cmd/decafein
