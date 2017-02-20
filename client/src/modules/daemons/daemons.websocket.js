var ws;

export const open = () => {
  if (ws) {
    ws.close();
  }
  ws = new WebSocket('ws://localhost:8080/ws/daemons');
  return ws;
};

export const close = () => {
  if (ws) {
    ws.close();
  }
};

export default {
  open,
  close
};
