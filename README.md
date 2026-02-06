## Sailor 🌊 ![beta](https://img.shields.io/badge/version-beta-blue?style=plastic)

## Update
+Connecting the database.

+Implementation in CI/CD

+New checks

## About
A cli utility designed for dynamic analysis of the security of Docker containers in CI/CD pipelines.
The utility starts the image, performs checks, and stops executing the image. The results are output in the terminal and can also be saved in JSON.

Сhecks:
- Open Ports Analysis
- Network Connections
- Environment Variables
- Root User Check
- Image History Check
- Writable Filesystem Check

Usage:
```sh
./sailor scan [name image] 
```
Flags:
```sh
./sailor scan [name image] --save-db --db-host [host] --db-port [port]
./sailor scan [name image] --save-result=/path/to/results.json #delete in realese!
```

Further development plans:
- The ability to scan multiple containers simultaneously
- Adding various checks
- Tips for correcting identified issues
