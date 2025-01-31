# Self-Updating Application
This program demonstrates an application that can update itself across different platforms (Windows, Mac, Linux). It accomplishes this through Docker containerization.

Pre-Requisites to Run
- Git
- Docker installed and running

Commands to Run
- git clone https://github.com/starfishpanda/Self-Updating-Application.git
- cd Self-Updating-Application
- docker build -t self-updating-app .
- docker run -it self-updating-app

Expected Output
You should expect to see something like this:

2025/01/30 01:28:57 Server starting on :8080
2025/01/30 01:28:57 Self-Updating Application is running version: 1.1.1
2025/01/30 01:28:57 Checking for updates...
2025/01/30 01:28:57 New version available: 1.1.2
2025/01/30 01:28:57 Successfully downloaded update.
2025/01/30 01:28:57 Successfully verified checksum of update.
2025/01/30 01:28:57 Successfully started new version: 1.1.2