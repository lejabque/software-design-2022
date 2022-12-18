import javafx.scene.Group;
import javafx.scene.Scene;
import javafx.scene.canvas.Canvas;
import javafx.scene.canvas.GraphicsContext;
import javafx.scene.paint.Color;
import javafx.stage.Stage;

public class JavaFxDrawingApi implements DrawingApi {
    private final int width;
    private final int height;
    private final Stage primaryStage;
    private final Canvas canvas;
    private final GraphicsContext gc;


    public JavaFxDrawingApi(Stage primaryStage, int width, int height) {
        this.width = width;
        this.height = height;
        this.primaryStage = primaryStage;
        this.canvas = new Canvas(width, height);
        this.gc = canvas.getGraphicsContext2D();
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
        gc.setFill(Color.BLUE);
        gc.fillOval(center.x() - radius, center.y() - radius, radius * 2, radius * 2);
    }

    @Override
    public void drawLine(Point from, Point to) {
        gc.strokeLine(from.x(), from.y(), to.x(), to.y());
    }

    @Override
    public void show() {
        var root = new Group();
        root.getChildren().add(canvas);
        primaryStage.setScene(new Scene(root));
        primaryStage.show();
    }
}
