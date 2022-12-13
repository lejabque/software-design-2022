import java.util.List;

public class MatrixGraph extends Graph {
    private final List<List<Boolean>> matrix;
    private final int vertexesCount;

    public MatrixGraph(List<List<Boolean>> matrix, DrawingApi drawingApi) {
        super(drawingApi);
        this.matrix = matrix;

        this.vertexesCount = matrix.size();
    }

    @Override
    protected void drawEdges() {
        for (int i = 0; i < matrix.size(); i++) {
            for (int j = 0; j < matrix.get(i).size(); j++) {
                if (matrix.get(i).get(j)) {
                    drawEdge(i, j);
                }
            }
        }
    }

    @Override
    protected int vertexesCount() {
        return vertexesCount;
    }
}
