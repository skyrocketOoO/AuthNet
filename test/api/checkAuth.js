import http from 'k6/http';
import { check, group, sleep } from 'k6';


export function TestCheckAuthAPI(relationUrl, headers){
    let resp;
    let edge = {
        "obj_ns": "role",
        "obj_name": "rd",
        "obj_rel": "parent",
        "sbj_ns": "role",
        "sbj_name": "rd-director",
    };

    resp = http.post(relationUrl + "/check", JSON.stringify(
        {
            "sbj": {
                "ns": edge['sbj_ns'],
                "name": edge['sbj_name']
            },
            "obj": {
                "ns": edge['obj_ns'],
                "name": edge['obj_name'],
                "rel": edge['obj_rel']
            }
        }),
        {headers:headers}
    )
    check(resp, { 'check': (r) => r.status == 403 });

    resp = http.post(relationUrl, JSON.stringify(
        {
            "edge":edge
        },
        {headers:headers}
    ))
    check(resp, { 'Create': (r) => r.status == 200 });

    resp = http.post(relationUrl + "/check", JSON.stringify(
        {
            "sbj": {
                "ns": edge['sbj_ns'],
                "name": edge['sbj_name']
            },
            "obj": {
                "ns": edge['obj_ns'],
                "name": edge['obj_name'],
                "rel": edge['obj_rel']
            }
        }),
        {headers:headers}
    )
    check(resp, { 'check2': (r) => r.status == 200 });
}