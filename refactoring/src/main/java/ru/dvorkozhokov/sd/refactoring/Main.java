package ru.dvorkozhokov.sd.refactoring;

import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.servlet.ServletContextHandler;
import org.eclipse.jetty.servlet.ServletHolder;
import ru.dvorkozhokov.sd.refactoring.database.Products;
import ru.dvorkozhokov.sd.refactoring.service.ProductsHtmlService;
import ru.dvorkozhokov.sd.refactoring.servlet.AddProductServlet;
import ru.dvorkozhokov.sd.refactoring.servlet.GetProductsServlet;
import ru.dvorkozhokov.sd.refactoring.servlet.QueryServlet;

import java.sql.DriverManager;
import java.sql.SQLException;

public class Main {
    private static final String DB_URL = "jdbc:sqlite:test.db"; // TODO: move to config

    public static void main(String[] args) throws Exception {
        ProductsHtmlService service;
        try {
            var connection = DriverManager.getConnection(DB_URL);
            service = new ProductsHtmlService(new Products(connection));
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }

        Server server = new Server(8081);

        ServletContextHandler context = new ServletContextHandler(ServletContextHandler.SESSIONS);
        context.setContextPath("/");
        server.setHandler(context);

        context.addServlet(new ServletHolder(new AddProductServlet(service)), "/add-product");
        context.addServlet(new ServletHolder(new GetProductsServlet(service)), "/get-products");
        context.addServlet(new ServletHolder(new QueryServlet(service)), "/query");

        server.start();
        server.join();
    }
}
