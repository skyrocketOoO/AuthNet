import http from 'k6/http';
import { check, group, sleep } from 'k6';

export function TestGetObjAuthAPI(relationUrl, headers){
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

    resp = http.post(relationUrl + "/obj-auths", JSON.stringify(
        {
            "sbj": {
                "ns": edge["sbj_ns"],
                "name": edge["sbj_name"]
            }
        }),
        {headers:headers}
    )
    check(resp, { 'findObj': (r) => r.status == 200 });
    check(resp, {  'findobj body': (r) => {
        return JSON.parse(r.body)["vertices"].length == 1
    }})
}