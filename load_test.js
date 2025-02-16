import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 100 },
        { duration: '30s', target: 1000 },
        { duration: '10s', target: 0 }
    ],
    thresholds: {
        http_req_duration: ['p(99)<50'],
        http_req_failed: ['rate<0.0001'],
    }
};

const BASE_URL = 'http://localhost:8080';

function getToken() {
    let payload = JSON.stringify({ username: 'user1', password: 'user1' });
    let headers = { 'Content-Type': 'application/json' };
    let res = http.post(`${BASE_URL}/api/auth`, payload, { headers });

    check(res, { 'auth success': (r) => r.status === 200 });

    let cookies = res.cookies['accessToken'];
    return cookies ? cookies[0].value : null;
}
export default function () {
    let token = getToken();

    if (!token) {
        console.error("Auth failed, no token received!");
        return;
    }

    let headers = {
        'Content-Type': 'application/json',
        'Cookie': `accessToken=${token}`
    };

    let res = http.post(`${BASE_URL}/api/buy/socks`, null, { headers });

    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 50ms': (r) => r.timings.duration < 50
    });

    sleep(1);
}
