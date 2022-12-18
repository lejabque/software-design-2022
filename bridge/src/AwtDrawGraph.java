import java.awt.*;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;
import java.util.function.Function;


public class AwtDrawGraph extends Frame implements GraphDrawerApplication {
    private final Function<DrawingApi, Graph> GraphFactory;

    AwtDrawGraph(Function<DrawingApi, Graph> graphFactory) {
        this.GraphFactory = graphFactory;
    }

    @Override
    public void paint(Graphics g) {
        super.paint(g);
        var graphics2D = (Graphics2D) g;
        graphics2D.clearRect(0, 0, 1280, 720);

        DrawingApi drawingApi = new AwtDrawingApi((Graphics2D) g, 1280, 720);
        var graph = GraphFactory.apply(drawingApi);
        graph.drawGraph();
    }

    public void drawGraph() {
        this.addWindowListener(new WindowAdapter() {
            public void windowClosing(WindowEvent we) {
                System.exit(0);
            }
        });
        this.setSize(1280, 720);
        this.setVisible(true);
    }
}