import { ChatOllama } from "@langchain/community/chat_models/ollama";
import { OllamaEmbeddings } from "@langchain/community/embeddings/ollama";

export const chatModel = new ChatOllama({
  baseUrl: "http://localhost:11434", // Default value
  model: "gemma",
});

export const embeddings = new OllamaEmbeddings({
  model: "nomic-embed-text",
  maxConcurrency: 5,
});
