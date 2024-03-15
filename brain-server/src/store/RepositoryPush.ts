import { ObjectsBatcher } from "weaviate-ts-client";
import { embeddings } from "../models";
import { vectorClient } from "./vector";

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
export default class RepositoryPushStore {
  public async initialize() {
    const definition = await vectorClient.schema
      .classCreator()
      .withClass(RepositoryPushVectorClass)
      .do();

    console.log("Created vector class:", definition);
  }

  public async add(content: string) {
    const resultEmbedding = await embeddings.embedQuery(content);
    const res = await vectorClient.data
      .creator()
      .withClassName(RepositoryPushVectorClass.class)
      .withProperties({
        content,
      })
      .withVector(resultEmbedding)
      .do();

    console.log("INSERTED REPOSITORY PUSH EVENT", res);
  }

  public async addBatch(batchContent: string[]) {
    const resultEmbedding = await embeddings.embedDocuments(batchContent);

    let batcher: ObjectsBatcher = vectorClient.batch.objectsBatcher();
    let counter: number = 0;
    let batchSize: number = 50;

    for (const content of batchContent) {
      const obj = {
        class: RepositoryPushVectorClass.class,
        properties: {
          content,
        },
        vector: resultEmbedding[counter],
      };

      batcher = batcher.withObject(obj);
      if (counter++ % batchSize === 0) {
        await batcher.do();
        batcher = vectorClient.batch.objectsBatcher();
      }
    }

    await batcher.do();
    console.log("BATCH INSERT OF REPOSITORY PUSH EVENTS DONE");
  }

  public async search(
    searchText: string,
    filters: {
      limit: number;
      distance: number;
    } = {
      limit: 10,
      distance: 0.15,
    }
  ) {
    const resultEmbedding = await embeddings.embedQuery(searchText);
    const result = await vectorClient.graphql
      .get()
      .withClassName(RepositoryPushVectorClass.class)
      .withNearVector({
        vector: resultEmbedding,
      })
      .withLimit(filters.limit)
      .withFields("content")
      .do();

    const resultString = (result.data.Get.RepositoryPush as any[])
      .map((entry) => entry.content)
      .join("\n");

    return resultString;
  }
}
