import http from "k6/http";

export const options = {
  vus: 150, // concurrent user
  //duration: "1s",
  iterations: 500, // total items attempt to buy
};

// TEST
// export default function () {
//   http.get("http://localhost:8080/ping");
// }

// different vus attempt
// 10 ok
// 50 ok
// 100 ok
// 250 fail

// PHASE 1
// export default function () {
//   const payload = JSON.stringify({
//     user_id: 69,
//     item_id: 1,
//   });
//   const headers = { "Content-Type": "application/json" };
//   http.post("http://localhost:8080/buy-item", payload, { headers });
// }
// PHASE 2
const urls = [
  "http://localhost:8080",
  "http://localhost:8081",
  "http://localhost:8082",
];

export default function () {
  http.batch(
    urls.map((url) => ({
      method: "POST",
      url: `${url}/buy-item-with-dis-lock`,
      body: JSON.stringify({
        user_id: 69,
        item_id: 1,
      }),
      params: {
        headers: { "Content-Type": "application/json" },
      },
    })),
  );
}
