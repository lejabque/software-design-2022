import javafx.application.Application;
import javafx.stage.Stage;

import java.awt.*;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;
import java.awt.geom.Ellipse2D;
import java.util.List;


public class AwtDrawGraph extends Frame {

    @Override
    public void paint(Graphics g) {
        super.paint(g);
        var graphics2D = (Graphics2D) g;
        graphics2D.clearRect(0, 0, 1280, 720);

        DrawingApi drawingApi = new AwtDrawingApi((Graphics2D) g, 1280, 720);
        Graph graph = new ListsGraph(List.of(
                new Edge(0, 1)
        ), drawingApi);
        graph.drawGraph();

        // Graph graph = new MatrixGraph(List.of(
        //         List.of(false, true, true, true),
        //         List.of(true, false, true, true),
        //         List.of(true, true, false, true),
        //         List.of(true, true, true, false)
        // ), drawingApi);
        // graph.drawGraph();
    }

    public static void main(String[] args) {
        Frame frame = new AwtDrawGraph();
        frame.addWindowListener(new WindowAdapter() {
            public void windowClosing(WindowEvent we) {
                System.exit(0);
            }
        });
        frame.setSize(1280, 720);
        frame.setVisible(true);
    }
}