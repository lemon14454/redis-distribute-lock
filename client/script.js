import http from "k6/http";

export const options = {
  vus: 1000, // concurrent user
  // duration: "1s",
  iterations: 1000, // total items attempt to buy
};

export default function () {
  const payload = JSON.stringify({
    user_id: 69,
    item_id: 1,
  });
  const headers = { "Content-Type": "application/json" };

  http.post("http://localhost:8080/buy-item-lock", payload, { headers });
  // single instance -> OK
  // multiple instance -> NO OK
}
