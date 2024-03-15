import weaviate from "weaviate-ts-client";

const vectorClient = weaviate.client({
  scheme: process.env.WEAVIATE_SCHEME || "http",
  host: process.env.WEAVIATE_HOST || "localhost:8080",
});

const RepositoryPushVectorClass = {
  class: "RepositoryPush",
  vectorIndexConfig: {
    distance: "cosine",
  },
  properties: [
    {
      name: "content",
      dataType: ["text"],
    },
  ],
};

async function getBatchWithCursor(
  collectionName: string,
  batchSize: number,
  cursor: string
): Promise<any[]> {
  // First prepare the query to run through data
  const query = vectorClient.graphql
    .get()
    .withClassName(collectionName)
    .withFields("content _additional { id vector }")
    .withLimit(batchSize);

  if (cursor) {
    // Fetch the next set of results
    let result = await query.withAfter(cursor).do();
    return result.data.Get[collectionName];
  } else {
    // Fetch the first set of results
    let result = await query.do();
    return result.data.Get[collectionName];
  }
}

const viewAll = async () => {
  // STEP 2 - Iterate through the data
  let cursor = null;

  // Batch import all objects to the target instance
  while (true) {
    // Get Request next batch of objects
    let nextBatch = await getBatchWithCursor(
      RepositoryPushVectorClass.class,
      100,
      cursor
    );

    // Break the loop if empty – we are done
    if (nextBatch.length === 0) break;

    // Here is your next batch of objects
    console.log(JSON.stringify(nextBatch));

    // Move the cursor to the last returned uuid
    cursor = nextBatch.at(-1)["_additional"]["id"];
  }
};

viewAll();
