## Sailor 🌊 ![alpha](https://img.shields.io/badge/version-alpha-blue?style=plastic)

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
./sailor scan [name image] --save-result=/path/to/results.json
```

Further development plans:
- Connecting the database
- Improving the CLI interface
- The ability to scan multiple containers simultaneously
- Adding various checks
- Tips for correcting identified issues
- Implementation in CI/CD
- Add more flags 
