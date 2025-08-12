import { app } from "./app";
import { raccsController } from "./routes/raccoons";

app.use(raccsController);

export default {
  async fetch(request: Request): Promise<Response> {
    return await app.fetch(request);
  },
};
