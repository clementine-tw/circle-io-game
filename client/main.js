import { getRandomInt } from "./random.js";
import { Circle } from "./circle.js";
import { Player } from "./player.js";

const canvas = document.getElementById("game");
const ctx = canvas.getContext("2d");

// tune pixel ratio
const screenRatio = [16, 9];
const height = 720;
const width = (height * screenRatio[0]) / screenRatio[1];

// css presentation
canvas.style.width = `${width}px`;
canvas.style.height = `${height}px`;

// canvas size
const dpr = window.devicePixelRatio || 1;
canvas.width = width * dpr;
canvas.height = height * dpr;

// draw size
ctx.scale(dpr, dpr);

const controlledPlayer = new Player(new Circle(20, 20, 10));
const otherPlayers = new Map();

function renderScreen() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  controlledPlayer.render(ctx);
  for (const [_, player] of otherPlayers) {
    player.render(ctx);
  }
}

function updateState() {
  controll(controlledPlayer.circle);
}

class Input {
  contructor() {
    this.w = false;
    this.s = false;
    this.a = false;
    this.d = false;
  }

  isMoving() {
    return this.w || this.s || this.a || this.d;
  }
}

const input = new Input();

function controll(target) {
  if (input.w) target.coord.y--;
  if (input.s) target.coord.y++;
  if (input.a) target.coord.x--;
  if (input.d) target.coord.x++;
}

document.addEventListener("keydown", function (event) {
  if (event.key == "w") input.w = true;
  if (event.key == "s") input.s = true;
  if (event.key == "a") input.a = true;
  if (event.key == "d") input.d = true;
});

document.addEventListener("keyup", function (event) {
  if (event.key == "w") input.w = false;
  if (event.key == "s") input.s = false;
  if (event.key == "a") input.a = false;
  if (event.key == "d") input.d = false;
});

function update() {
  updateState();
  renderScreen();
  requestAnimationFrame(update);
}

requestAnimationFrame(update);

function handleClickStart() {
  startButton.disabled = true;

  const uuid = crypto.randomUUID();
  const socket = new WebSocket("ws://localhost:8080?token=" + uuid);

  socket.onopen = function () {
    console.log("websocket opened");
    // send the player position to server while player first joins the game
    socket.send(JSON.stringify(controlledPlayer.circle));
  };
  socket.onclose = function (e) {
    console.log("websocket closed: " + e.code);
    clearInterval(sendMsgId);
    startButton.disabled = false;
  };
  socket.onmessage = function (e) {
    const data = JSON.parse(e.data);

    console.log("recieved message: " + e.data);
    if (data.id == uuid) {
      controlledPlayer.circle.radius = data.radius;
      return;
    }

    if (!otherPlayers.has(data.id)) {
      otherPlayers.set(
        data.id,
        new Player(
          new Circle(data.coord.x, data.coord.y, data.radius),
          `rgb(${getRandomInt(0, 255)},${getRandomInt(0, 255)},${getRandomInt(0, 255)})`,
        ),
      );
      return;
    }

    let c = otherPlayers.get(data.id).circle;
    c.coord = data.coord;
    c.radius = data.radius;
  };

  let { x: lx, y: ly } = controlledPlayer.circle.coord;
  const sendMsgId = setInterval(() => {
    let { x: cx, y: cy } = controlledPlayer.circle.coord;
    if (
      socket &&
      socket.readyState === WebSocket.OPEN &&
      (cx != lx || cy != ly)
    ) {
      lx = controlledPlayer.circle.coord.x;
      ly = controlledPlayer.circle.coord.y;
      socket.send(JSON.stringify(controlledPlayer.circle));
    }
  }, 50);
}

const startButton = document.getElementById("button-start");
startButton.addEventListener("click", handleClickStart);
