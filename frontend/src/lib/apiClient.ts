// src/lib/apiClient.ts
// Central Axios instance for all Go backend calls.
// Go wraps every response as: { success: bool, message: string, data: T }
// The response interceptor unwraps `.data` so callers receive T directly.

import axios from "axios";

const GO_API =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

const apiClient = axios.create({
  baseURL: GO_API,
  headers: {
    "Content-Type": "application/json",
  },
});

// Unwrap Go's envelope: { success, message, data } → data
apiClient.interceptors.response.use(
  (response) => {
    const body = response.data;
    // Go always wraps with a `data` field when there is a payload
    if (body && typeof body === "object" && "data" in body) {
      response.data = body.data;
    }
    return response;
  },
  (error) => {
    // Preserve the Go error message when available
    const goMessage =
      error.response?.data?.message ?? error.response?.data?.error;
    if (goMessage) {
      error.message = goMessage;
    }
    return Promise.reject(error);
  }
);

export default apiClient;
