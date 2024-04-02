import http from 'k6/http';
import { check, group, sleep } from 'k6';


export function TestCreateAPI(relationUrl, headers){
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


    resp = http.post(relationUrl, JSON.stringify(
        {
            "edge":edge
        },
        {headers:headers}
    ))
    check(resp, { 'Create': (r) => r.status == 200 });

    resp = http.post(relationUrl, JSON.stringify(
        {
            "edge":edge
        },
        {headers:headers}
    ))
    check(resp, { 'Create duplicate': (r) => r.status == 400 });
}
