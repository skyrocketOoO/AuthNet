import http from 'k6/http';
import { check, group, sleep } from 'k6';


export function TestDeleteAPI(relationUrl, headers){
    let resp;
    let edge = {
        "obj_ns": "role",
        "obj_name": "rd",
        "obj_rel": "parent",
        "sbj_ns": "role",
        "sbj_name": "rd-director",
    };

    resp = http.del(`${relationUrl}/all`, null, {headers:headers});
    check(resp, { 'ClearAllRelations': (r) => r.status == 200 });

    resp = http.del(relationUrl, JSON.stringify(
        {
            "edge":edge,
            "query_mode": false
        },
        {headers:headers}
    ))
    check(resp, { 'not found': (r) => r.status == 404 });

    resp = http.post(relationUrl, JSON.stringify(
        {
            "edge":edge
        },
        {headers:headers}
    ))
    check(resp, { 'Create': (r) => r.status == 200 });

    resp = http.del(relationUrl, JSON.stringify(
        {
            "edge":edge,
            "query_mode": false
        },
        {headers:headers}
    ))
    check(resp, { 'ok': (r) => r.status == 200 });

    resp = http.post(relationUrl, JSON.stringify(
        {
            "edge":edge
        },
        {headers:headers}
    ))
    check(resp, { 'Create2': (r) => r.status == 200 });

    resp = http.del(relationUrl, JSON.stringify(
        {
            "edge":{
                obj_rel: edge['obj_rel']
            },
            "query_mode": true
        },
        {headers:headers}
    ))
    check(resp, { 'query delete': (r) => r.status == 200 });

    resp = http.get(relationUrl +
                    "?obj_ns=" + edge["obj_ns"] +
                    "&obj_name=" + edge["obj_name"] +
                    "&obj_rel=" + edge["obj_rel"] +
                    "&sbj_ns=" + edge["sbj_ns"] +
                    "&sbj_name=" + edge["sbj_name"] +
                    "&query_mode=false",
                    null, {headers:headers});

    check(resp, { 'check delete': (r) => r.status == 404 });

}
