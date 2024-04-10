import http from 'k6/http';
import { check, group, sleep } from 'k6';


export function TestGetAPI(relationUrl, headers){
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

    resp = http.get(relationUrl +
                    "?obj_ns=" + edge["obj_ns"] +
                    "&obj_name=" + edge["obj_name"] +
                    "&obj_rel=" + edge["obj_rel"] +
                    "&sbj_ns=" + edge["sbj_ns"] +
                    "&sbj_name=" + edge["sbj_name"] +
                    "&query_mode=false",
                    null, {headers:headers});

    check(resp, { 'match status': (r) => r.status == 200 });
    check(resp, { 'match body': (r) => {
        return JSON.parse(r.body)['edges'][0]['obj_rel'] === 'parent';
    }});

    resp = http.get(relationUrl +
                    "?sbj_name=" + edge["sbj_name"] +
                    "&query_mode=true",
                    null, {headers:headers});

    check(resp, { 'query status': (r) => r.status == 200 });
    check(resp, { 'query body': (r) => {
        const data = JSON.parse(r.body);
        return data['edges'][0]['obj_rel'] = 'parent';
    }});
    
    resp = http.get(relationUrl +
                    "?query_mode=true",
                    null, {headers:headers});

    check(resp, { 'query all status': (r) => r.status == 200 });
    check(resp, { 'query all body': (r) => {
        const data = JSON.parse(r.body);
        return data['edges'][0]['obj_rel'] = 'parent';
    }});

    resp = http.get(relationUrl +
                    "?obj_ns=" + edge["obj_ns"] +
                    "&obj_name=" + edge["obj_name"] +
                    "&obj_rel=" + edge["obj_rel"] +
                    "&sbj_ns=" + edge["sbj_ns"] +
                    "&query_mode=false",
                    null, {headers:headers});
    check(resp, { 'query not find': (r) => r.status == 404 });
}