# Portscan

Scan for open ports. Developed using golang.

## Usage

    portscan -h yoursite.com -min 1 -max 1024 -c 15
	
	Above command will scan for open ports on yoursite.com starting from port 1 and stopping at port 1024. 15 concurrent connections will be used for the scanning.
	
	## Installation
	
    go install github.com/nirandas/portscan
	