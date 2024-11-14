import http from "k6/http";

export const options = {
  vus: 2000, // concurrent user
  // duration: "1s",
  iterations: 2000, // total items attempt to buy
};

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
      // url: `${url}/buy-item-with-lock`,
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
