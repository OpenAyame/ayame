import * as WebSocket from "ws";
export const pubsubWs = new WebSocket.Server({ noServer: true });
interface Client extends WebSocket {
   roomId: string | null;
}

const clients = new Map<String, Client[]>();

pubsubWs.on("connection", (ws: Client) => {
  console.log("established websocket connection");
  ws.on("open", () =>  {
    console.log("connected");
  });
  ws.on("message", (data: string) => {
    console.log("onmessage----", data);
    const message = JSON.parse(data);
    const roomId: string = message.room_id;
    if (message.type === "register" && message.room_id) {
      const roomClients: Client[] = clients.get(roomId) || [];
      console.log(roomClients.length);
      const currentClientCount = roomClients.length;
      // 3 人目は入れない処理をする
      if (currentClientCount > 1) {
        console.log("over member count", currentClientCount);
        ws.send(JSON.stringify({"type": "reject"}));
        ws.close();
      } else {
        ws.roomId = roomId;
        roomClients.push(ws);
        clients.set(roomId, roomClients);
        ws.send(JSON.stringify({"type": "accept"}));
      }
    }
    else {
      if (ws.roomId) {
        const roomClients: Client[] = clients.get(ws.roomId);
        // クライアント全員にdataをbroadcast
        roomClients.forEach((client) => {
          // 送信主の場合は skip
          if (ws !== client) {
            client.send(data);
          }
        });
      }
    }
  });
  ws.on("ping", (bytes) => {
    console.log("onping----", bytes);
    ws.pong();
  });
  ws.on("pong", (bytes) => {
    // pong はスルーする
  });
});
