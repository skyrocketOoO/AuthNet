import http from 'k6/http';
import { check } from 'k6';
  
export function TestReservedWord(serverUrl, headers){
    const relationUrl = `${serverUrl}/relation`
    let res, payload;

    payload = {
        object_namespace: "test_f%ile",
        object_name: "1",
        relation: "write",
        subject_namespace: "teat",
        subject_name: "1",
        subject_relation: "read",
    };
    res = http.post(`${relationUrl}`, JSON.stringify(payload), {headers:headers});
    check(res, { '%': (r) => r.status != 200 });

    res = http.post(`${relationUrl}/clear-all-relations`, null, {headers:headers});
    check(res, { 'ClearAllRelations': (r) => r.status == 200 });
}