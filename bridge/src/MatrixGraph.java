import java.io.File;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

public class MatrixGraph extends Graph {
    private final List<List<Boolean>> matrix;
    private final int vertexesCount;

    public MatrixGraph(List<List<Boolean>> matrix, DrawingApi drawingApi) {
        super(drawingApi);
        this.matrix = matrix;
        this.vertexesCount = matrix.size();
    }

    public MatrixGraph(DrawingApi drawingApi) {
        this(readMatrixFromStdin(), drawingApi);
    }

    private static List<List<Boolean>> readMatrixFromStdin() {
        try (var scanner = new Scanner(System.in)) {
            int vertexesCount = scanner.nextInt();
            var matrix = new ArrayList<List<Boolean>>(vertexesCount);
            for (int i = 0; i < vertexesCount; i++) {
                var row = new ArrayList<Boolean>(vertexesCount);
                for (int j = 0; j < vertexesCount; j++) {
                    row.add(scanner.nextInt() == 1);
                }
                matrix.add(row);
            }
            return matrix;
        }
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
