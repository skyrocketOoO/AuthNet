import http from 'k6/http';
import { check } from 'k6';

export function Check(serverUrl, headers, layer, base){
    const relationUrl = `${serverUrl}/relation`
    let res, payload;
    const namespace = "role", relation = "parent";
    const start = "0_0";
    const end = (layer).toString() + "_" + (Math.pow(base, layer)-1).toString();

    payload = {
        object: {
            namespace: namespace,
            name: start,
            relation: relation,
        },
        subject: {
            namespace: namespace,
            name: end,
            relation: relation,
        },
    };
    res = http.post(`${relationUrl}/check`, JSON.stringify(payload), {
        headers: headers, 
        timeout: '900s',
    });
    check(res, { 'Check': (r) => r.status ==  200 });
};