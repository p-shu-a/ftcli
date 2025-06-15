### A file transfer utility
- Used to transfer a file over a network using symetric key for encryption

#### Command syntax:
- `ftcli -h`
- Send file:     `ftcli -rip=192.168.60.70 -role=send -file=/loc/to/file -pass=passw0rd`
- Recieve file:  `ftcli -role=recv -pass=passw0rd`

#### Immediate goals
- files must be encrypted with password during transfer ✅
- offer the reciever option to decline download ✅
- be able to send a file and validate file on recieving end (hashes must match) ✅
- at no point should the entire file be loaded in memory: 
    - No loading in memeory needed if using a stream cipher, with AES or ChaCha20 ✅
    - Not possible when adding a Poly1305 MAC, since the ciphertext must be computed before sending ❌
- Only use the Go standard lib
    - If using the AES ✅
    - Had to use an extended lib while implementing ChaCha20-Poly1305 ❌

#### Next steps
- for large files, offer compression
- if the download is interupted, you should be able to resume from there you left off.
- improve the syntax of the command:
    - `ftcli -send -rip 192.168.60.70 -f /loc/to/file` (prompted for pass)
