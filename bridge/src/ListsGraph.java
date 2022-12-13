import java.util.List;

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
