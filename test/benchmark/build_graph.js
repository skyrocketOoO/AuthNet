import http from 'k6/http';
import { check } from 'k6';

export function BuildGraph(serverUrl, headers, layer, base){
    const relationUrl = `${serverUrl}/relation`
    let res, payload;
    const namespace = "role", relation = "parent";
    
    res = http.del(`${relationUrl}/all`, null, {headers:headers});
    check(res, { 'ClearAllRelations': (r) => r.status == 200 });

    let curLayer = 1;
    while (curLayer <= layer){
        const count = Math.pow(base, curLayer);

        for (let i = 0; i < count; i++){
            payload = {
                edge:{
                    obj_ns: namespace,
                    obj_name: curLayer.toString() + "_" + i.toString(),
                    obj_rel: relation,
                    sbj_ns: namespace,
                    sbj_name: (curLayer-1).toString() + "_" + Math.floor(i / base).toString(),
                    sbj_rel: relation,
                }
            };
            res = http.post(`${relationUrl}`, JSON.stringify(payload), {headers: headers});
            check(res, { 'Create: status == 200': (r) => r.status == 200 });
        };

        curLayer += 1;
    };
};