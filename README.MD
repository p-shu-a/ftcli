### A file transfer utility

#### Command syntax:
`ftcli -h`
`ftcli -rip=192.168.60.70 -role=send -file=/loc/to/file -pass=passw0rd` to send
`ftcli -role=recv -pass=passw0rd` to recieve

#### Immediate goals
- be able to send a file and validate file on recieving end (hashes must match) ✅
- at no point should the entire file be loaded in memory for the sender
    - using io.copy() so ✅
- files must be encrypted with password during transfer
- offer the reciever option to decline download

#### Next steps
- for large files, offer compression
- if the download is interupted, you should be able to resume from there you left off.
- improve the syntax of the command:
    - `ftcli -send -rip 192.168.60.70 -f /loc/to/file` (prompted for pass)