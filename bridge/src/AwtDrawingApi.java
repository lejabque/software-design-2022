import java.awt.*;

public class AwtDrawingApi implements DrawingApi {
    private final int width;
    private final int height;

    private final Graphics2D graphics;

    public AwtDrawingApi(Graphics2D graphics, int width, int height) {
        this.width = width;
        this.height = height;
        this.graphics = graphics;
    }


    @Override
    public long getDrawingAreaWidth() {
        return width;
    }

    @Override
    public long getDrawingAreaHeight() {
        return height;
    }

    @Override
    public void drawCircle(Circle circle) {
        var center = circle.center();
        var radius = circle.radius();
        graphics.setColor(Color.GREEN);
        graphics.fillOval((int) center.x(), (int) center.y(), (int) radius, (int) radius);
    }

    @Override
    public void drawLine(Point from, Point to) {
        graphics.drawLine((int) from.x(), (int) from.y(), (int) to.x(), (int) to.y());
    }

    @Override
    public void show() {
        // do nothing
    }
}
