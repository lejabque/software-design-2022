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
        var radius = 10;
        drawingApi.drawCircle(new Circle(center, radius));
    }

    private Point getVertexPoint(int vertex) {
        var x = (vertex + 1) * drawingApi.getDrawingAreaWidth() / (vertexesCount() + 1);
        var y = drawingApi.getDrawingAreaHeight() / 2;
        return new Point(x, y);
    }
}