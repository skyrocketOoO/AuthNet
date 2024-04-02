import http from 'k6/http';
import { check, group, sleep } from 'k6';

export function TestGetSbjAuthAPI(relationUrl, headers){
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

    resp = http.post(relationUrl + "/sbj-who-has-auth", JSON.stringify(
        {
            "obj": {
                "ns": edge["obj_ns"],
                "name": edge["obj_name"],
                "rel": "parent",
            }
        }),
        {headers:headers}
    )
    check(resp, { 'find': (r) => r.status == 200 });
    check(resp, {  'find body': (r) => {
        return JSON.parse(r.body)["vertices"].length == 1
    }})
}