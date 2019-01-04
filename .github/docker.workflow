workflow "Docker build" {
  on = "push"
  resolves = ["Log in to Docker registry", "Push Docker image"]
}

action "Build Docker image" {
  uses = "actions/docker/cli@76ff57a"
  args = "build -t $GITHUB_REPOSITORY ."
}

action "Log in to Docker registry" {
  uses = "actions/docker/login@76ff57a"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Push Docker image" {
  uses = "actions/docker/cli@76ff57a"
  needs = ["Build Docker image", "Log in to Docker registry"]
  args = "push $GITHUB_REPOSITORY"
}
