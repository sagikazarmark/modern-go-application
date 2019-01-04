workflow "Docker build" {
  on = "push"
  resolves = [
    "Push Docker image",
    "Only push Docker image when on master",
  ]
}

action "Build Docker image" {
  uses = "actions/docker/cli@76ff57a"
  args = "build -t $GITHUB_REPOSITORY ."
}

action "Log in to Docker registry" {
  uses = "actions/docker/login@76ff57a"
  needs = ["Only push Docker image when on master"]
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Push Docker image" {
  uses = "actions/docker/cli@76ff57a"
  needs = ["Build Docker image", "Log in to Docker registry", "Only push Docker image when on master"]
  args = "push $GITHUB_REPOSITORY"
}

action "Only push Docker image when on master" {
  uses = "actions/bin/filter@b2bea07"
  args = "branch master"
}
