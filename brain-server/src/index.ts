import Express from "express";
import { ChatOllama } from "@langchain/community/chat_models/ollama";
import RepositoryPushStore from "./store/RepositoryPush";

const chatModel = new ChatOllama({
  baseUrl: "http://localhost:11434", // Default value
  model: "gemma",
});

const app = Express();
const rsStore = new RepositoryPushStore();

app.use(Express.json());

app.get("/", async (req, res) => {
  const question = req.query.question as string;
  const repositoryPushResult = await rsStore.search(question);
  const answer = await chatModel.invoke(repositoryPushResult + question);

  return res.send(answer.content);
});

app.post("/repository-push/batch", async (req, res) => {
  const batchContent: string[] = req.body.batchContent;
  await rsStore.addBatch(batchContent);

  res.send();
});

app.listen(7000, () => {
  console.log("Listening on port 7000");
});
