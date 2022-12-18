record Point(double x, double y) {
}

record Circle(Point center, double radius) {
}

public interface DrawingApi {
    long getDrawingAreaWidth();

    long getDrawingAreaHeight();

    void drawCircle(Circle circle);

    void drawLine(Point from, Point to);

    void show();
}