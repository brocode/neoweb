server {
    listen_addr = "127.0.0.1:8080"
}

log {
    format = "text"
    level = "debug"
}

nvim {
    cmd = "docker"
    args = ["run", "--network", "none", "--memory=40m", "--memory-swap=40m", "--rm", "-i", "nvim", "--embed"]
    forwardEnvVars = ["DOCKER_HOST", "PATH"]
}
