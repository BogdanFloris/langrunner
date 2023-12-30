type Response = {
  type: "output" | "error";
  payload: string;
};

export const run = async (
  language: string,
  code: string,
): Promise<Response> => {
  const response = await fetch("http://localhost:3000/run", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ language, code }),
  });
  if (!response.ok) {
    throw new Error(`Server error: ${response.status}`);
  }
  const json = await response.json();
  return json;
};
