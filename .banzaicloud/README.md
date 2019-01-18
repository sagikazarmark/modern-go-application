# Spotguide for Modern Go Application

[Banzai Cloud](https://banzaicloud.com/) [Spotguides](https://banzaicloud.com/tags/spotguides/) are 
**CI/CD driven cloud-native application templates**
which makes it possible to painlessly transition from code to a production ready Kubernetes deployment.

This spotguide contains a [Modern Go Application](https://github.com/sagikazarmark/modern-go-application) boilerplate
integrated into the Banzai Cloud Pipeline CI/CD flow,
serving as a reference implementation for private Spotguide templates.

You can easily try it on the **[Banzai Cloud Pipeline](https://banzaicloud.com)** Public Beta platform:

<p align="center">
    <a href="http://beta.banzaicloud.io/">
        <img src="https://banzaicloud.com/img/try_pipeline_button.svg" />
      </a>
</p>


## Usage

This spotguide is not part of the official spotguide catalog for a reason: it demonstrates how easily you can create
private spotguide templates. (Don't get confused by the fact that this is a public repository: it works with private repos too).

Here are the steps for launching this spotguide:

1. Fork the repository (or create a copy of it manually)
2. Add `spotguide` to the list of repository topics
3. Create a Github release for the latest stable tag at https://github.com/YOURORG/modern-go-application/releases/new
4. Go to https://beta.banzaicloud.io/ui/YOURORG/spotguide/create
5. Click the synchronization icon in the right upper corner
6. Modern Go Application should appear, use it as any other spotguide in the catalog
7. *Recommended:* After completing the spotguide flow, visit the generated repository for further instructions.
