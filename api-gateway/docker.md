# Building and Pushing docker images to GCP registry
- Build docker image using `docker build -t <tag> .`
- Tag the image to gcr `docker tag <image> eu.gcr.io/<project>/<image>:<version>`
- Pushing the image - `docker push eu.gcr.io/<project>/<image>:<version>`

# Required permissions and pre-push
# Steps to push docker image to google cloud registry
1. Authentication and Project Setup:
- `gcloud auth print-access-token | docker login -u oauth2accesstoken --password-stdin https://<region>-docker.pkg.dev` This command will configure docker to push to gcloud registry.
- Additionally you can also use `https://eu.gcr.io`.

2. Configure Docker for GCP
- `gcloud auth configure-docker`

3. IAM Roles:
    - Artifact Registry Administrator
    - roles/artifactregistry.createOnPushRepoAdmin
    - Storage Admin

4. Export the Credentials:
 - `export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/credentials.json"`

5. Tag the Docker Image:
 - `docker tag image <pkg>/project_id/repository/image:v1`.
 - Here `pkg` uses `<region>-docker.pkg.dev` or `gcr.io` from docker configuration.
 - You must run step 2 for this to work.

6. Push the Docker Image:
 - `docker push image`

references: 
- https://cloud.google.com/artifact-registry/docs/docker/pushing-and-pulling
- https://cloud.google.com/build/docs/build-push-docker-image