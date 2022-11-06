package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.database.Products;
import ru.dvorkozhokov.sd.refactoring.models.Product;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.sql.*;
import java.util.List;

public class GetProductsServlet extends HttpServlet {
    private final Products products;

    public GetProductsServlet(Products db) {
        products = db;
    }

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws IOException {
        List<Product> products;
        try {
            products = this.products.getProducts();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        var writer = response.getWriter();
        writer.println("<html><body>");
        for (var product : products) {
            writer.println(product.getName() + "\t" + product.getPrice() + "</br>");
        }
        writer.println("</body></html>");
        response.setContentType("text/html");
        response.setStatus(HttpServletResponse.SC_OK);
    }
}
