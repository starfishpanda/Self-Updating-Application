# 🚀 Self-Updating Application  

This project demonstrates a self-updating application that works across multiple platforms (Windows, Mac, Linux). It achieves this through Docker containerization and an update mechanism that verifies and applies new versions seamlessly and authentically.  


## ✅ Features

- **Digital Signature Verification:** Ensures the authenticity and integrity of the update by using public and private keys, protecting against man-in-the-middle (MITM) attacks.
- **Checksum Verification:** Validates that the update binary has not been tampered with by checking its SHA-256 hash checksum.
- **Docker Containerization:** Enables the program to update itself across different operating systems, as it is self-contained with its runtime and dependencies.


## 📌 Prerequisites  

Before running the application, ensure you have the following installed:  

- [Git](https://git-scm.com/downloads)  
- [Docker](https://www.docker.com/get-started) (installed and running)  


## 🛠️ How to Run  

Clone the repository and build the Docker container:  

```sh
git clone https://github.com/starfishpanda/Self-Updating-Application.git
cd Self-Updating-Application
docker build -t self-updating-app .
docker run -it self-updating-app
```


## 📜 Expected Output  

Upon running the application, you should see something like this:  

- `2025/01/31 06:33:21 Successfully loaded public key from /app/client/public.pem`
- `2025/01/31 06:33:21 Self-Updating Application is running version: 1.1.1`
- `2025/01/31 06:33:21 Checking for updates...`
- `2025/01/31 06:33:21 Successfully loaded private key from /app/server/private.pem`
- `2025/01/31 06:33:21 Server starting on :8080`
- `2025/01/31 06:33:21 New version available: 1.1.2`
- `2025/01/31 06:33:21 Successfully downloaded update.`
- `2025/01/31 06:33:21 Successfully verified binary signature`
- `2025/01/31 06:33:21 Successfully verified checksum of update.`
- `2025/01/31 06:33:21 Successfully started new version: 1.1.2`
