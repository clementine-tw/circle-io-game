export class Circle {
  constructor(x, y, radius) {
    this.coord = {};
    this.coord.x = x;
    this.coord.y = y;
    this.radius = radius;
  }

  render(ctx, strokeStyle) {
    ctx.beginPath();
    ctx.strokeStyle = strokeStyle;
    ctx.arc(this.coord.x, this.coord.y, this.radius, 0, 2 * Math.PI);
    ctx.stroke();
  }
}
