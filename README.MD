### A file transfer utility
- Used to transfer a file over a network using symetric key for encryption
- Using ChaCha20-Poly1305 for encryption. KDF is Argon2.
- Keen focus on memory efficiency
    - Files are loaded into memory in 4096-bytes chunks. This helps reduce the memory load from large files
    - MACs are calculated for each individual chunk
- Current focus: resumption of an interupted file transfer

#### Command syntax:
- `ftcli -h`
- Send file:     `ftcli -rip=192.168.60.70 -role=send -file=/loc/to/file -pass=passw0rd`
- Recieve file:  `ftcli -role=recv -pass=passw0rd`

#### Technical criteria
- files must be encrypted with password during transfer ✅
- offer the reciever option to decline download ✅
- be able to send a file and validate file on recieving end (hashes must match) ✅
- at no point should the entire file be loaded in memory: ✅
    - No loading in memory needed if using a stream cipher, with AES or ChaCha20 (alone, no MAC)
    - Not possible when adding a Poly1305 MAC, since the ciphertext must be computed before sending
    - However, we can send the file in chunks, computing the MAC for all of them
- Only use the Go standard lib ❌
    - Possible if using AES
    - Have to use an extended lib while implementing ChaCha20-Poly1305

#### Next steps
- if the download is interupted, you should be able to resume from there you left off.
- for large files, offer compression
- improve the syntax of the command:
    - `ftcli -send -rip 192.168.60.70 -f /loc/to/file` (prompted for pass)
- parallel processing for super-large files? files over X chunks in size
