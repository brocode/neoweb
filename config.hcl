
server {
  listen_addr = "127.0.0.1:8080"
}

log {
    format = "text"
    level = "debug"
}

nvim {
    cmd = "docker"
    args = ["run", "--rm", "-i", "nvim", "--embed"]
	forwardEnvVars = ["DOCKER_HOST", "PATH"]
}
