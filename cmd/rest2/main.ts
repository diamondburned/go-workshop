#!/usr/bin/env -S deno run -A
import { serve } from "https://deno.land/std/http/server.ts";

type echoRequest = {
	message: string;
};

function handleEcho(request: Request): Response {
	console.log(request.body);
}

const s = serve({ port: 8000 });
