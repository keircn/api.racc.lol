import { Elysia } from "elysia";
import { respond } from "../lib/respond";

const mainController = new Elysia().get("/", async (context) => {
  return respond(200, {
    success: true,
    message:
      "welcome! this is the racc api, view my docs here https://racc.lol/documentation",
  });
});

export default mainController;
