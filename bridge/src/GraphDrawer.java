import java.util.function.Function;

public class GraphDrawer {
    public static void main(String[] args) {
        if (args.length != 2) {
            System.out.println("Usage: java GraphDrawer <drawing_api> <graph_type>");
            System.out.println("drawing_api: awt or javafx");
            System.out.println("graph_type: list or matrix");
            return;
        }

        // 2 arg - graph type
        Function<DrawingApi, Graph> graphFactory;
        switch (args[1]) {
            case "list" -> graphFactory = ListsGraph::new;
            case "matrix" -> graphFactory = MatrixGraph::new;
            default -> {
                System.out.println("Unknown graph type: " + args[1]);
                return;
            }
        }

        GraphDrawerApplication application;
        switch (args[0]) {
            case "awt" ->
                    application = new AwtDrawGraph(graphFactory);
            case "javafx" ->
                    application = new JavaFxDrawGraph(graphFactory);
            default -> {
                System.out.println("Unknown drawing api: " + args[0]);
                return;
            }
        }
        application.drawGraph();
    }
}
