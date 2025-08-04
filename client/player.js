export class Player {
  constructor(circle, strokeStyle = "black") {
    this.circle = circle;
    this.strokeStyle = strokeStyle;
  }

  render(ctx) {
    this.circle.render(ctx, this.strokeStyle);
  }
}
