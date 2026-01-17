import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
  scenarios: {
    steady_users: {
      executor: 'constant-vus',
      vus: 50,         
      duration: '30s',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<2000'],
  },
};

const BASE_URL = 'http://localhost:8080';

// ---------- HELPERS ----------
function randomDate() {
  return `2025-${randomIntBetween(1, 12)
    .toString().padStart(2, '0')}-${randomIntBetween(1, 28)
    .toString().padStart(2, '0')}T${randomIntBetween(8, 17)}:00:00Z`;
}

export default function () {
  // ===============================
  // 1. SIGNUP
  // ===============================
  const email = `patient_${uuidv4()}@test.com`;
  const password = 'password123';

  const signupPayload = JSON.stringify({
    firstname: 'Test',
    lastname: 'User',
    email: email,
    password: password,
    phone: '9999999999',
  });

  const signupRes = http.post(
    `${BASE_URL}/auth/signup`,
    signupPayload,
    { headers: { 'Content-Type': 'application/json' } }
  );

  const signupOk = check(signupRes, {
    'signup succeeded': (r) => r.status === 200 || r.status === 201,
  });

  if (!signupOk) {
    return; 
  }

  const token = signupRes.json('token');

  if (!token) {
    return;
  }

  const authHeaders = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };

  // ===============================
  // 2. CREATE APPOINTMENTS
  // ===============================
  const appointmentPayload = JSON.stringify({
    appointmentType: 'CONSULTATION',
    appointmentDate: randomDate(),
    startTime: '10:00',
    endTime: '11:00',
    mode: 'Online',
    notes: 'k6 signup → appointment test',
    doctorId: '11111111-1111-1111-1111-111111111111',
  });

  const apptRes = http.post(
    `${BASE_URL}/appointments`,
    appointmentPayload,
    { headers: authHeaders }
  );

  check(apptRes, {
    'appointment created': (r) => r.status === 200 || r.status === 201,
  });

  sleep(1);
}
