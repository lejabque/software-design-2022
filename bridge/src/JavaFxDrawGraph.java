import javafx.application.Application;
import javafx.stage.Stage;

import java.util.List;

public class JavaFxDrawGraph {
    public static class Start extends Application {
        public static void run(String[] args) {
            launch(args);
        }

        @Override
        public void start(Stage primaryStage) {
            DrawingApi drawingApi = new JavaFxDrawingApi(primaryStage, 1280, 720);

            Examples.drawList(drawingApi);
            // Examples.drawMatrix(drawingApi);
        }
    }

    // for some reasons it doesn't work when launched directly from Idea
    // but works when launched via such proxy class
    public static void main(final String[] args) {
        Start.run(args);
    }
}
