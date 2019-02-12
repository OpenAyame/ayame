import * as WebSocket from "ws";
export const pubsubWs = new WebSocket.Server({ noServer: true });

pubsubWs.on("connection", (ws: WebSocket) => {
  console.log("established websocket connection");
  const currentClientCount = pubsubWs.clients.size;
  // 3 人目は入れない処理をする
  if (currentClientCount > 2) {
    console.log("over member count", currentClientCount);
    ws.send(JSON.stringify({"type": "close"}));
    ws.close();
  }
  ws.on("open", () =>  {
    console.log("connected");
  });
  ws.on("message", (data) => {
    console.log("onmessage----", data);
    // クライアント全員にdataをbroadcast
    pubsubWs.clients.forEach((client) => {
      // 送信主の場合は skip
      if (ws !== client) {
        client.send(data);
      }
    });
  });
  ws.on("ping", (bytes) => {
    console.log("onping----", bytes);
    ws.pong();
  });
  ws.on("pong", (bytes) => {
    // pong はスルーする
  });
  ws.on("close", (bytes) => {
    console.log("onclose----", bytes);
    // close させる
    pubsubWs.clients.forEach((client) => {
      client.close();
    });
  });
});
