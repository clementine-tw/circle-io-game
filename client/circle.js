export class Circle {
  constructor(x, y) {
    this.coord = {};
    this.coord.x = x;
    this.coord.y = y;
  }

  render(ctx, strokeStyle) {
    ctx.beginPath();
    ctx.strokeStyle = strokeStyle;
    ctx.arc(this.coord.x, this.coord.y, 10, 0, 2 * Math.PI);
    ctx.stroke();
  }
}
