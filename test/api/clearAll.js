import http from 'k6/http';
import { check, group, sleep } from 'k6';


export function TestClearAllAPI(relationUrl, headers){
    let resp;
    let edge = {
        "obj_ns": "role",
        "obj_name": "rd",
        "obj_rel": "parent",
        "sbj_ns": "role",
        "sbj_name": "rd-director",
    };

    resp = http.post(relationUrl, JSON.stringify(
        {
            "edge":edge
        },
        {headers:headers}
    ))
    check(resp, { 'Create': (r) => r.status == 200 });

    edge['obj_ns'] = "role2"
    resp = http.post(relationUrl, JSON.stringify(
        {
            "edge":edge
        },
        {headers:headers}
    ))
    check(resp, { 'Create2': (r) => r.status == 200 });


    resp = http.del(`${relationUrl}/all`, null, {headers:headers});
    check(resp, { 'ClearAllRelations': (r) => r.status == 200 });

    resp = http.get(relationUrl + "?query_mode=true", null, {headers:headers});
    const data = JSON.parse(resp.body)
    check(data, { 'check length': (d) => d["edges"].length == 0 });
}