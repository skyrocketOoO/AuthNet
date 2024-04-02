import http from 'k6/http';
import { check } from 'k6';

export function BuildGraph(serverUrl, headers, layer, base){
    const relationUrl = `${serverUrl}/relation`
    let res, payload;
    const namespace = "role", relation = "parent";
    
    res = http.post(`${relationUrl}/clear-all-relations`, null, {headers:headers});
    check(res, { 'ClearAllRelations': (r) => r.status == 200 });

    let curLayer = 1;
    while (curLayer <= layer){
        const count = Math.pow(base, curLayer);

        for (let i = 0; i < count; i++){
            payload = {
                object_namespace: namespace,
                object_name: curLayer.toString() + "_" + i.toString(),
                relation: relation,
                subject_namespace: namespace,
                subject_name: (curLayer-1).toString() + "_" + Math.floor(i / base).toString(),
                subject_relation: relation,
            };
            res = http.post(`${relationUrl}`, JSON.stringify(payload), {headers: headers});
            check(res, { 'Create: status == 200': (r) => r.status == 200 });
        };

        curLayer += 1;
    };
};