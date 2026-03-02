# Sailor 🌊 ![realese](https://img.shields.io/badge/version-v1.0.0-blue?style=plastic)
![image](https://github.com/EWinterhalter/sailor/blob/beta/src/DOCKER.png)
## About
A simple Go CLI app for runtime checks of docker containers. With the ability to save the results to a database. It is supposed to be used in CI/CI pipeline.

Сhecks:
- Network Connections (warn\meduim)
- Image Version (warn\meduim)
- Writable FS (warn\meduim)
- Open Ports (warn\meduim\high)
- Image History (high)
- Env (high)
- Netcat Listener (critical)
- Netcat Process (critical)
- Open Port 4444 (critical)
- Root  User (critical)
- Remote Shell (critical)

Usage:
```sh
./sailor scan [name image] 
```
With database:
```sh
./sailor scan [name image] --save-db --db-host [host] --db-port [port]
```
Example:
```sh
[INFO] Starting security scan for image: test-app-bad:id
[INFO] Starting container...
[INFO] Container started: id
[🐳] Container ID: id
[⏰] Scan started: timestamp

  ⠋ Open Ports Analysis
✓ PASS Open Ports [61ms] 0 open ports
  ⠋ Network Connections
✓ PASS Connections [61ms] 0 established, 0 listening
  ⠋ Environment Variables
✗ FAIL Environment [54ms] Found: DB_PASSWORD, DB_PASSWORD, DB_PASSWORD, API_KEY
  ⠋ Root User Check
✗ FAIL Root User [49ms] UID=0
  ⠋ Image History Check
✗ FAIL Image History [31ms] Secrets detected: PASSWORD, SECRET, API_KEY
  ⠋ Writable Filesystem Check
⚠ WARN Writable FS [52ms] Filesystem is writable
  ⠋ Netcat Listener Check
✓ PASS Netcat Listener [52ms] not detected
  ⠋ Netcat Process Check
✓ PASS Netcat Process [51ms] not running
  ⠋ Port 4444 Check
✓ PASS Port 4444 [51ms] closed
  ⠋ Remote Shell Check
✓ PASS Remote Shell [51ms] no shell

Total Checks:    10
✓ Passed:       6
⚠ Warnings:     1
✗ Failed:       3

Severity Breakdown:
● High:      3
● Medium:    1
● Low:       6
⏱ Total Time:   855ms

⚠ SECURITY ISSUES DETECTED ⚠
```
Use in CI/CD pipline:
![image](https://github.com/EWinterhalter/sailor/blob/beta/src/PIPELINE.png)

Verification criteria:

IThub-college Graduation Project 2026 by Elina Bulanova 4R2.22