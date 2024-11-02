server {
    listen_addr = "127.0.0.1:8080"
}

log {
    format = "text"
    level = "info"
}

nvim {
    #cmd = "docker"
    #args = ["run", "-p", "6666:6666", "--memory=100m", "--memory-swap=150m", "--rm", "-i", "nvim", "--embed", "--listen", "127.0.0.1:6666"]
    forwardEnvVars = ["DOCKER_HOST", "PATH"]
    cmd = "nvim"
    args = ["--embed", "--listen", "127.0.0.1:6666"]
}
