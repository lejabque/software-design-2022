import java.util.List;

public class MatrixGraph extends Graph {
    private final List<List<Boolean>> matrix;
    private final int vertexesCount;

    public MatrixGraph(List<List<Boolean>> matrix, DrawingApi drawingApi) {
        super(drawingApi);
        this.matrix = matrix;

        int vertexesCount = 0;
        for (var row : matrix) {
            for (var cell : row) {
                vertexesCount += cell ? 1 : 0;
            }
        }
        this.vertexesCount = vertexesCount;
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
