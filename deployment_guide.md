# Deployment Guide
## Prerequisites
* Docker installed on the system
* Docker Compose installed on the system

## Deployment Steps
1. Clone the repository: `git clone https://github.com/axentx/social-assembly.git`
2. Navigate to the repository directory: `cd social-assembly`
3. Build the Docker image: `docker build -t social-assembly .`
4. Start the container: `docker run -p 8080:8080 social-assembly`