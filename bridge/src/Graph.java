public abstract class Graph {
    /**
     * Bridge to drawing api
     */
    protected DrawingApi drawingApi;

    public Graph(DrawingApi drawingApi) {
        this.drawingApi = drawingApi;
    }

    public void drawGraph() {
        for (int i = 0; i < vertexesCount(); i++) {
            drawVertex(i);
        }
        drawEdges();
        drawingApi.show();
    }

    protected abstract void drawEdges();

    protected abstract int vertexesCount();

    protected void drawEdge(int from, int to) {
        var fromPoint = getVertexPoint(from);
        var toPoint = getVertexPoint(to);
        drawingApi.drawLine(fromPoint, toPoint);
    }

    private void drawVertex(int vertex) {
        var center = getVertexPoint(vertex);
        var radius = 20;
        drawingApi.drawCircle(new Circle(center, radius));
    }

    private Point getVertexPoint(int vertex) {
        var meshSize = (long) Math.ceil(Math.sqrt(vertexesCount()));
        var x = vertex % meshSize + 1; // [1, .., meshSize]
        var y = vertex / meshSize + 1;
        return new Point(x * drawingApi.getDrawingAreaWidth() / (meshSize + 1),
                y * drawingApi.getDrawingAreaHeight() / (meshSize + 1));
    }
}