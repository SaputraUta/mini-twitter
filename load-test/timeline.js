import http from 'k6/http';
import { check } from 'k6';

const BASE = __ENV.BASE || 'http://localhost:8002';
const TLPATH = __ENV.TLPATH || '/timeline';

export const options = {
    stages: [
        {
            duration: '10s', target: 300
        },
        {
            duration: '30s', target: 300
        },
        {
            duration: '5s', target: 0
        },
    ]
}

export default function () {
    const user = Math.floor(Math.random() * 10000) + 1;
    const res = http.get(`${BASE}${TLPATH}/${user}`);
    check(res, { 'status 200': (r) => r.status === 200})
}