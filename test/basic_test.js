import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { TestGetAPI } from './api/get.js';
import { TestCreateAPI } from './api/create.js';
import { TestDeleteAPI } from './api/delete.js';
import { TestClearAllAPI } from './api/clearAll.js';
import { TestCheckAuthAPI } from './api/checkAuth.js';
import { TestGetObjAuthAPI } from './api/getObjAuths.js';
import { TestGetSbjAuthAPI } from './api/getSbjsWhoHasAuth.js';

export const options = {
  vus: 1,
}

export default function() {
  const SERVER_URL = "http://localhost:8081"
  const Headers = {
    'Content-Type': 'application/json',
  }

  let res = http.get(`${SERVER_URL}/ping`);
  check(res, { 'Server can ping': (r) => r.status == 200 });

  res = http.get(`${SERVER_URL}/healthy`);
  check(res, { 'Server is healthy': (r) => r.status == 200 });

  res = http.del(`${SERVER_URL}/relation/all`, null, {headers:Headers});
  check(res, { 'ClearAllRelations': (r) => r.status == 200 });

  group("api", () => {
    const RELATION_URL = SERVER_URL + "/relation"
    group("get", () => {
      TestGetAPI(RELATION_URL, Headers);
    })
    group("create", () => {
      TestCreateAPI(RELATION_URL, Headers);
    })
    group("delete", () => {
      TestDeleteAPI(RELATION_URL, Headers);
    })
    group("clearAll", () => {
      TestClearAllAPI(RELATION_URL, Headers);
    })
    group("checkAuth", () => {
      TestCheckAuthAPI(RELATION_URL, Headers);
    })
    group("getObjAuths", () => {
      TestGetObjAuthAPI(RELATION_URL, Headers);
    })
    group("getSbjsWhoHasAuth", () => {
      TestGetSbjAuthAPI(RELATION_URL, Headers);
    })
    // TestClearAllAPI(RELATION_URL, Headers);
  });

  // group("scenario", () => {
  //   group("cycle", () => {
  //     TestCycle(SERVER_URL, Headers);
  //   });
  //   group("require_attr", () => {
  //     TestRequiredAttr(SERVER_URL, Headers);
  //   });
  //   group("reserved_word", () => {
  //     TestReservedWord(SERVER_URL, Headers);
  //   });
  // });
}
