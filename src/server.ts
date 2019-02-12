import * as http from "http";
import * as url from "url";
import * as nodeStatic from "node-static";
import { pubsubWs } from "./ws";

const PORT = process.env.PORT || 3000;
const samples = new nodeStatic.Server("./sample");
const server = http.createServer((req, res) => {
  if (process.env.SAMPLE) {
    req.addListener("end", () => {
      samples.serve(req, res);
    }).resume();
  }
});

// websocket endpoint のハンドリング
server.on("upgrade", (request, socket, head) => {
  const pathname = url.parse(request.url).pathname;
  if (pathname === "/ws") {
    // pubsub
    pubsubWs.handleUpgrade(request, socket, head, (ws) =>  {
      pubsubWs.emit("connection", ws, request);
    });
  } else {
    socket.destroy();
  }
});

server.listen(PORT, (error: Error) => {
    if (error) {
      console.error("Unable to listen on port", PORT, error);
      return;
    }
  console.log("server is now running port=", PORT);
});

