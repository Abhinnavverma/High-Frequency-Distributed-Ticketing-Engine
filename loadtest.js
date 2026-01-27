import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  stages: [
    { duration: '10s', target: 300 },
    { duration: '30s', target: 1500 }, 
    { duration: '10s', target: 0 },
  ],
};

export default function () {
  const BASE_URL = 'http://localhost'; 
  const params = { headers: { 'Content-Type': 'application/json' } };

  // 🆔 1. Generate Identity
  const uniqueId = `${__VU}-${__ITER}-${Date.now()}`; 
  const email = `user${uniqueId}@test.com`;
  const password = 'password123';

  // 📝 2. Register
  const registerPayload = JSON.stringify({ email: email, password: password });
  http.post(`${BASE_URL}/register`, registerPayload, params);

  // 🔑 3. Login
  const loginPayload = JSON.stringify({ email: email, password: password });
  const loginRes = http.post(`${BASE_URL}/login`, loginPayload, params);

  if (loginRes.status !== 200) return; // Skip if login failed
  const token = loginRes.json('token');
  
  // 🎟️ 4. The Booking Attack (Using ID now!)
  // Since we inserted 100 seats, the IDs are guaranteed to be 1 to 100.
  const seatId = randomIntBetween(1, 100);

  const bookingPayload = JSON.stringify({
    seat_id: seatId, // 🚨 Matches your new DB Schema/Go Struct
  });

  const authParams = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
  };

  const bookRes = http.post(`${BASE_URL}/bookings`, bookingPayload, authParams);

  check(bookRes, {
    'System Integrity': (r) => r.status === 201 || r.status === 409,
    'Booked': (r) => r.status === 201,
  });

  sleep(Math.random() * 0.5);
}