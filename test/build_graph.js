import { group, check } from "k6"
import { BuildGraph } from "./benchmark/build_graph.js"
import http from 'k6/http';


export const options = {
  vus: 1,
}

export default function() {
  const SERVER_URL = "http://localhost:8080"
  const Headers = {
    'Content-Type': 'application/json',
  }
  const layer = 6, base = 6;

  let res = http.del(`${SERVER_URL}/relation/`, null, {headers:Headers});
  check(res, { 'ClearAllRelations: status == 200': (r) => r.status == 200 });

  group("build graph", () => {
    BuildGraph(SERVER_URL, Headers, layer, base);
  })
}
