# How to push to artifact-registry

1. Configure docker registry

2. Create a tag of the desired image:
   docker tag image:tag LOCATION-docker.pkg.dev/PROJECT_ID/ARTIFACT_REGISTRY_REPOSITORY/IMAGE:TAG

   LOCATION = asia-east1
   ARTIFACT_REGISTRY_REPOSITORY = ai

3. Push to registry
   docker push LOCATION.pkg.dev/PROJECT_ID/ARTIFACT_REGISTRY_REPOSITORY/IMAGE:TAG

# How to run DB migrations [Migration Docs](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md)

- migrate create -ext sql -dir db/migrations -seq create_users_table
- migrate -database DATABASE_URL -path ./db/migrations up
- migrate -database DATABASE_URL -path ./db/migrations down

# TO DOs

[] explore self-hosting option

- Embedding generation: [nomic-embed-text](https://ollama.com/library/nomic-embed-text)
- Vector database: [Weaviate](https://weaviate.io/developers/weaviate)
- LLM model: [gemma](https://ollama.com/library/gemma)
- Processor: [Langchain](https://js.langchain.com/docs/get_started/quickstart)

[] add payload validation to webhook requests

[] add validation for trigger endpoint

[] setup scheduled job in production to trigger publishing of events
