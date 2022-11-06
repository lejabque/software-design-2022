package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.database.Products;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.sql.*;

public class QueryServlet extends HttpServlet {
    private final Products products;

    public QueryServlet(Products db) {
        products = db;
    }

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws IOException {
        String command = request.getParameter("command");

        var writer = response.getWriter();
        if ("max".equals(command)) {
            try {
                var max = products.getMaxProduct();
                writer.println("<html><body>");
                writer.println("<h1>Product with max price: </h1>");
                writer.println(max.getName() + "\t" + max.getPrice() + "</br>");
                writer.println("</body></html>");
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        } else if ("min".equals(command)) {
            try {
                var min = products.getMinProduct();
                writer.println("<html><body>");
                writer.println("<h1>Product with min price: </h1>");
                writer.println(min.getName() + "\t" + min.getPrice() + "</br>");
                writer.println("</body></html>");
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        } else if ("sum".equals(command)) {
            try {
                var sum = products.getSum();
                writer.println("<html><body>");
                writer.println("Summary price: ");
                writer.println(sum);
                writer.println("</body></html>");
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        } else if ("count".equals(command)) {
            try {
                var cnt = products.getCount();
                writer.println("<html><body>");
                writer.println("Number of products: ");
                writer.println(cnt);
                writer.println("</body></html>");
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        } else {
            writer.println("Unknown command: " + command);
        }

        response.setContentType("text/html");
        response.setStatus(HttpServletResponse.SC_OK);
    }

}
