# Run docker-compose locally:

docker-compose -f docker-compose.dev.yaml up --build

# How to push to artifact-registry

1. Configure docker registry

2. Create a tag of the desired image:
   docker tag image:tag LOCATION-docker.pkg.dev/PROJECT_ID/ARTIFACT_REGISTRY_REPOSITORY/IMAGE:TAG

   LOCATION = asia-east1
   ARTIFACT_REGISTRY_REPOSITORY = ai

3. Push to registry
   docker push LOCATION.pkg.dev/PROJECT_ID/ARTIFACT_REGISTRY_REPOSITORY/IMAGE:TAG

# How to run DB migrations [https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md]

- migrate create -ext sql -dir db/migrations -seq create_users_table
- migrate -database DATABASE_URL -path ./db/migrations up
- migrate -database DATABASE_URL -path ./db/migrations down
