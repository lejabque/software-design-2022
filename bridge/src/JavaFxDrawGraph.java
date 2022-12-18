import javafx.application.Application;
import javafx.stage.Stage;

import java.util.function.Function;

public class JavaFxDrawGraph implements GraphDrawerApplication {
    private final Function<DrawingApi, Graph> GraphFactory;

    public static class JavaFxApp extends Application {
        static Function<DrawingApi, Graph> GraphFactory;

        public static void run() {
            launch();
        }

        @Override
        public void start(Stage primaryStage) {
            DrawingApi drawingApi = new JavaFxDrawingApi(primaryStage, 1280, 720);
            var graph = GraphFactory.apply(drawingApi);
            graph.drawGraph();
        }
    }


    public JavaFxDrawGraph(Function<DrawingApi, Graph> graphFactory) {
        this.GraphFactory = graphFactory;
    }

    @Override
    public void drawGraph() {
        // todo: move factory to globals?
        JavaFxApp.GraphFactory = GraphFactory;
        JavaFxApp.run();
    }
}
