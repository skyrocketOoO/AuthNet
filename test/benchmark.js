import { group } from "k6"
import { BuildGraph } from "./benchmark/build_graph.js"
import { Check } from "./benchmark/check.js"


export const options = {
  vus: 1,
}

export default function() {
  const SERVER_URL = "http://localhost:8080"
  const Headers = {
    'Content-Type': 'application/json',
  }
  const layer = 6, base = 5;

  group("check", () => {
    Check(SERVER_URL, Headers, layer, base);
  })
}
