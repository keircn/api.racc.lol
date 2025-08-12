import { Elysia } from "elysia";
import { facts } from "../data/facts";
import { respond } from "../lib/respond";
import { fileService } from "../lib/fileService";
import { photoRateLimiter } from "../lib/rateLimit";

function getRateLimitHeaders(rateLimitResult: {
  remaining: number;
  resetTime: number;
}) {
  return {
    "X-RateLimit-Limit": "5",
    "X-RateLimit-Remaining": rateLimitResult.remaining.toString(),
    "X-RateLimit-Reset": rateLimitResult.resetTime.toString(),
  };
}

function getRateLimitErrorHeaders(
  rateLimitResult: { resetTime: number },
  retryAfter: number
) {
  return {
    "Content-Type": "application/json",
    "X-RateLimit-Limit": "5",
    "X-RateLimit-Remaining": "0",
    "X-RateLimit-Reset": rateLimitResult.resetTime.toString(),
    "Retry-After": retryAfter.toString(),
  };
}

function getTimeBasedIndex(total: number, period: string): number {
  const now = new Date();
  let seed = now.getFullYear() + now.getMonth();

  if (period === "hourly") seed += now.getDate() + now.getHours();
  else if (period === "daily") seed += now.getDate();
  else if (period === "weekly") seed += Math.floor(now.getDate() / 7);
  else return Math.floor(Math.random() * total);

  return seed % total;
}

export const raccsController = new Elysia()
  .get("/v1/raccoon", ({ request }) => {
    const url = new URL(request.url);
    const newUrl =
      url.origin + url.pathname.replace("/v1/raccoon", "/raccoon") + url.search;
    return Response.redirect(newUrl, 302);
  })
  .get(
    "/raccoon",
    async ({
      query,
      request,
    }: {
      query: Record<string, string | undefined>;
      request: Request;
    }) => {
      const rateLimitResult = photoRateLimiter.checkRateLimit(request);

      if (!rateLimitResult.allowed) {
        const retryAfter = Math.ceil(
          (rateLimitResult.resetTime - Date.now()) / 1000
        );
        return new Response(
          JSON.stringify({
            success: false,
            error: "Rate limit exceeded",
            message:
              "Too many photo requests. You can make 5 requests every 10 seconds.",
            retryAfter,
          }),
          {
            status: 429,
            headers: getRateLimitErrorHeaders(rateLimitResult, retryAfter),
          }
        );
      }
      try {
        const files = await fileService.listFiles("", ".jpg");

        if (files.length === 0) {
          return respond(404, { success: false, error: "No JPG images found" });
        }

        const period = query.daily
          ? "daily"
          : query.hourly
          ? "hourly"
          : query.weekly
          ? "weekly"
          : null;
        const imageIndex = period
          ? getTimeBasedIndex(files.length, period)
          : Math.floor(Math.random() * files.length);

        const selectedFile = files[imageIndex];
        const fileBuffer = await fileService.getFile(selectedFile.path);

        if (!fileBuffer) {
          return respond(404, { success: false, error: "Image not found" });
        }

        if (query.json === "true") {
          const url = new URL(request.url);
          const baseUrl = `${url.protocol}//${url.host}`;

          return new Response(
            JSON.stringify({
              success: true,
              data: {
                url: `${baseUrl}/raccoon`,
                size: selectedFile.size,
                type: period || "random",
                contentType: "image/jpeg",
              },
            }),
            {
              status: 200,
              headers: {
                "Content-Type": "application/json",
                ...getRateLimitHeaders(rateLimitResult),
              },
            }
          );
        }

        return new Response(new Uint8Array(fileBuffer), {
          headers: {
            "Content-Type": "image/jpeg",
            "Cache-Control": "no-cache",
            ...getRateLimitHeaders(rateLimitResult),
          },
        });
      } catch (error: any) {
        console.error("Error in /raccoon endpoint:", error);
        return respond(500, {
          success: false,
          error: "Failed to retrieve raccoon image",
          message: error.message || "Unknown error",
        });
      }
    }
  )
  .get(
    "/raccoon/:id",
    async ({
      params,
      query,
      request,
    }: {
      params: { id: string };
      query: Record<string, string | undefined>;
      request: Request;
    }) => {
      const rateLimitResult = photoRateLimiter.checkRateLimit(request);

      if (!rateLimitResult.allowed) {
        const retryAfter = Math.ceil(
          (rateLimitResult.resetTime - Date.now()) / 1000
        );
        return new Response(
          JSON.stringify({
            success: false,
            error: "Rate limit exceeded",
            message:
              "Too many photo requests. You can make 5 requests every 10 seconds.",
            retryAfter,
          }),
          {
            status: 429,
            headers: getRateLimitErrorHeaders(rateLimitResult, retryAfter),
          }
        );
      }

      try {
        const files = await fileService.listFiles("", ".jpg");
        const id = parseInt(params.id);

        if (isNaN(id) || id < 1 || id > files.length) {
          return respond(404, {
            success: false,
            error: "Raccoon not found",
            message: `Raccoon ${params.id} does not exist. Available range: 1-${files.length}`,
          });
        }

        const selectedFile = files[id - 1];
        const fileBuffer = await fileService.getFile(selectedFile.path);

        if (!fileBuffer) {
          return respond(404, {
            success: false,
            error: "Raccoon file not found",
          });
        }

        if (query.json === "true") {
          const url = new URL(request.url);
          const baseUrl = `${url.protocol}//${url.host}`;

          return new Response(
            JSON.stringify({
              success: true,
              data: {
                url: `${baseUrl}/raccoon/${id}`,
                id: id,
                size: selectedFile.size,
                contentType: "image/jpeg",
              },
            }),
            {
              status: 200,
              headers: {
                "Content-Type": "application/json",
                ...getRateLimitHeaders(rateLimitResult),
              },
            }
          );
        }

        return new Response(new Uint8Array(fileBuffer), {
          headers: {
            "Content-Type": "image/jpeg",
            "Cache-Control": "no-cache",
            ...getRateLimitHeaders(rateLimitResult),
          },
        });
      } catch (error: any) {
        console.error("Error in /raccoon/:id endpoint:", error);
        return respond(500, {
          success: false,
          error: "Failed to retrieve raccoon image",
          message: error.message || "Unknown error",
        });
      }
    }
  )
  .get("/raccoons", async ({ request }) => {
    const files = await fileService.listFiles("", ".jpg");
    const url = new URL(request.url);
    const baseUrl = `${url.protocol}//${url.host}`;

    const imageData = files.map((file, i: number) => {
      return {
        url: `${baseUrl}/raccoon/${i + 1}`,
        index: i + 1,
        size: file.size,
      };
    });

    return { success: true, data: imageData };
  })
  .get("/memes", async ({ request }) => {
    const files = await fileService.listFiles("memes", ".jpg");
    const url = new URL(request.url);
    const baseUrl = `${url.protocol}//${url.host}`;

    const imageData = files.map((file, i: number) => {
      return {
        url: `${baseUrl}/meme/${i + 1}`,
        index: i + 1,
        size: file.size,
      };
    });

    return { success: true, data: imageData };
  })
  .get(
    "/video",
    async ({
      query,
      request,
    }: {
      query: Record<string, string | undefined>;
      request: Request;
    }) => {
      const rateLimitResult = photoRateLimiter.checkRateLimit(request);

      if (!rateLimitResult.allowed) {
        const retryAfter = Math.ceil(
          (rateLimitResult.resetTime - Date.now()) / 1000
        );
        return new Response(
          JSON.stringify({
            success: false,
            error: "Rate limit exceeded",
            message:
              "Too many photo requests. You can make 5 requests every 10 seconds.",
            retryAfter,
          }),
          {
            status: 429,
            headers: getRateLimitErrorHeaders(rateLimitResult, retryAfter),
          }
        );
      }
      const files = await fileService.listFiles("videos", ".mp4");

      if (files.length === 0) {
        return { success: false, error: "No videos found" };
      }

      const randomIndex = Math.floor(Math.random() * files.length);
      const selectedFile = files[randomIndex];
      const fileBuffer = await fileService.getFile(selectedFile.path);

      if (!fileBuffer) {
        return { success: false, error: "Video not found" };
      }

      const filename = selectedFile.name;

      if (query.json === "true") {
        const url = new URL(request.url);
        const baseUrl = `${url.protocol}//${url.host}`;

        return new Response(
          JSON.stringify({
            success: true,
            data: {
              url: `${baseUrl}/video`,
              filename,
              size: selectedFile.size,
              format: "mp4",
            },
          }),
          {
            status: 200,
            headers: {
              "Content-Type": "application/json",
              ...getRateLimitHeaders(rateLimitResult),
            },
          }
        );
      }

      return new Response(new Uint8Array(fileBuffer), {
        headers: {
          "Content-Type": "video/mp4",
          "Cache-Control": "no-cache",
          ...getRateLimitHeaders(rateLimitResult),
        },
      });
    }
  )
  .get("/fact", ({ query }: { query: Record<string, string | undefined> }) => {
    const randomFact = facts[Math.floor(Math.random() * facts.length)];
    return query.json === "true"
      ? { success: true, data: { fact: randomFact, total: facts.length } }
      : new Response(randomFact, { headers: { "Content-Type": "text/plain" } });
  })
  .get(
    "/meme",
    async ({
      query,
      request,
    }: {
      query: Record<string, string | undefined>;
      request: Request;
    }) => {
      const rateLimitResult = photoRateLimiter.checkRateLimit(request);

      if (!rateLimitResult.allowed) {
        const retryAfter = Math.ceil(
          (rateLimitResult.resetTime - Date.now()) / 1000
        );
        return new Response(
          JSON.stringify({
            success: false,
            error: "Rate limit exceeded",
            message:
              "Too many photo requests. You can make 5 requests every 10 seconds.",
            retryAfter,
          }),
          {
            status: 429,
            headers: getRateLimitErrorHeaders(rateLimitResult, retryAfter),
          }
        );
      }

      try {
        const files = await fileService.listFiles("memes", ".jpg");

        if (files.length === 0) {
          return respond(404, {
            success: false,
            error: "No meme images found",
          });
        }

        const period = query.daily
          ? "daily"
          : query.hourly
          ? "hourly"
          : query.weekly
          ? "weekly"
          : null;
        const imageIndex = period
          ? getTimeBasedIndex(files.length, period)
          : Math.floor(Math.random() * files.length);

        const selectedFile = files[imageIndex];
        const fileBuffer = await fileService.getFile(selectedFile.path);

        if (!fileBuffer) {
          return respond(404, {
            success: false,
            error: "Meme image not found",
          });
        }

        if (query.json === "true") {
          const url = new URL(request.url);
          const baseUrl = `${url.protocol}//${url.host}`;

          return new Response(
            JSON.stringify({
              success: true,
              data: {
                url: `${baseUrl}/meme`,
                size: selectedFile.size,
                type: period || "random",
                contentType: "image/jpeg",
              },
            }),
            {
              status: 200,
              headers: {
                "Content-Type": "application/json",
                ...getRateLimitHeaders(rateLimitResult),
              },
            }
          );
        }

        return new Response(new Uint8Array(fileBuffer), {
          headers: {
            "Content-Type": "image/jpeg",
            "Cache-Control": "no-cache",
            ...getRateLimitHeaders(rateLimitResult),
          },
        });
      } catch (error: any) {
        console.error("Error in /meme endpoint:", error);
        return respond(500, {
          success: false,
          error: "Failed to retrieve meme image",
          message: error.message || "Unknown error",
        });
      }
    }
  )
  .get(
    "/meme/:id",
    async ({
      params,
      query,
      request,
    }: {
      params: { id: string };
      query: Record<string, string | undefined>;
      request: Request;
    }) => {
      const rateLimitResult = photoRateLimiter.checkRateLimit(request);

      if (!rateLimitResult.allowed) {
        const retryAfter = Math.ceil(
          (rateLimitResult.resetTime - Date.now()) / 1000
        );
        return new Response(
          JSON.stringify({
            success: false,
            error: "Rate limit exceeded",
            message:
              "Too many photo requests. You can make 5 requests every 10 seconds.",
            retryAfter,
          }),
          {
            status: 429,
            headers: getRateLimitErrorHeaders(rateLimitResult, retryAfter),
          }
        );
      }

      try {
        const files = await fileService.listFiles("memes", ".jpg");
        const id = parseInt(params.id);

        if (isNaN(id) || id < 1 || id > files.length) {
          return respond(404, {
            success: false,
            error: "Meme not found",
            message: `Meme ${params.id} does not exist. Available range: 1-${files.length}`,
          });
        }

        const selectedFile = files[id - 1];
        const fileBuffer = await fileService.getFile(selectedFile.path);

        if (!fileBuffer) {
          return respond(404, {
            success: false,
            error: "Meme file not found",
          });
        }

        if (query.json === "true") {
          const url = new URL(request.url);
          const baseUrl = `${url.protocol}//${url.host}`;

          return new Response(
            JSON.stringify({
              success: true,
              data: {
                url: `${baseUrl}/meme/${id}`,
                id: id,
                size: selectedFile.size,
                contentType: "image/jpeg",
              },
            }),
            {
              status: 200,
              headers: {
                "Content-Type": "application/json",
                ...getRateLimitHeaders(rateLimitResult),
              },
            }
          );
        }

        return new Response(new Uint8Array(fileBuffer), {
          headers: {
            "Content-Type": "image/jpeg",
            "Cache-Control": "no-cache",
            ...getRateLimitHeaders(rateLimitResult),
          },
        });
      } catch (error: any) {
        console.error("Error in /meme/:id endpoint:", error);
        return respond(500, {
          success: false,
          error: "Failed to retrieve meme image",
          message: error.message || "Unknown error",
        });
      }
    }
  )
  .get("/stats", async () => {
    const photoFiles = await fileService.listFiles("", ".jpg");
    const videoFiles = await fileService.listFiles("videos", ".mp4");
    const memeFiles = await fileService.listFiles("memes", ".jpg");

    const photos = photoFiles.length;
    const videos = videoFiles.length;
    const memes = memeFiles.length;

    return { success: true, data: { photos, videos, memes } };
  })
  .all("*", () => {
    return respond(404, {
      success: false,
      error: "Endpoint not found",
      message:
        "The requested endpoint does not exist. Check out our documentation for available endpoints.",
      documentation: "https://racc.lol/documentation",
    });
  });
