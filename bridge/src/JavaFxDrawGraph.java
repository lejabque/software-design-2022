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
            Graph graph = new ListsGraph(List.of(
                    new Edge(0, 1),
                    new Edge(0, 2),
                    new Edge(0, 3),
                    new Edge(1, 2),
                    new Edge(1, 3),
                    new Edge(2, 3)
            ), drawingApi);
            graph.drawGraph();

            // Graph graph = new MatrixGraph(List.of(
            //         List.of(false, true, true, true),
            //         List.of(true, false, true, true),
            //         List.of(true, true, false, true),
            //         List.of(true, true, true, false)
            // ), drawingApi);
            // graph.drawGraph();
        }
    }

    // for some reasons it doesn't work when launched directly from Idea
    // but works when launched via such proxy class
    public static void main(final String[] args) {
        Start.run(args);
    }
}
