import java.io.File;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

record Edge(int from, int to) {
}

public class ListsGraph extends Graph {
    private final List<Edge> edges;
    private final int vertexesCount;

    public ListsGraph(List<Edge> edges, DrawingApi drawingApi) {
        super(drawingApi);
        this.edges = edges;

        int vertexesCount = 0;
        for (var edge : edges) {
            vertexesCount = Math.max(vertexesCount, edge.from());
            vertexesCount = Math.max(vertexesCount, edge.to());
        }
        this.vertexesCount = vertexesCount + 1;
    }


    public ListsGraph(DrawingApi drawingApi) {
        this(readMatrixFromStdin(), drawingApi);
    }

    private static List<Edge> readMatrixFromStdin() {
        try (var scanner = new Scanner(System.in)) {
            int vertexesCount = scanner.nextInt();
            int edgesCount = scanner.nextInt();
            var edges = new ArrayList<Edge>(edgesCount);
            for (int i = 0; i < edgesCount; i++) {
                edges.add(new Edge(scanner.nextInt(), scanner.nextInt()));
            }
            return edges;
        }
    }

    @Override
    public void drawEdges() {
        for (var edge : edges) {
            drawEdge(edge.from(), edge.to());
        }
    }

    @Override
    protected int vertexesCount() {
        return vertexesCount;
    }
}
