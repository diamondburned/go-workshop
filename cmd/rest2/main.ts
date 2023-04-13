#!/usr/bin/env -S deno run -A cmd/rest2/main.ts
import { serve } from "https://deno.land/std@0.182.0/http/server.ts";

async function sh(cmd: string) {
  const r = Deno.run({
    cmd: ["sh", "-c", cmd],
    stdout: "inherit",
    stderr: "inherit",
  });
  await r.status();
}

serve((req: Request) => {
  // No built-in router in Deno :(
  switch ((new URL(req.url)).pathname) {
    case "/echo":
      return handleEcho(req);
    default:
      return new Response("Not Found", { status: 404 });
  }
}, { port: 12345 });

// POST to the endpoint with a null body.
await sh(`httpie -p hb POST http://localhost:12345/echo <<< null`);
Deno.exit(0);

type echoRequest = {
  message: string;
};

async function handleEcho(request: Request): Promise<Response> {
  const resp = await request.json() as echoRequest; // performs no validation
  console.log("JS says", resp.message); // 💥
  return new Response();
}
