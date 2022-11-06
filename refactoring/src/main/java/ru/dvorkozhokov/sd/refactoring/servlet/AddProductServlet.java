package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.database.Products;
import ru.dvorkozhokov.sd.refactoring.models.Product;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.sql.*;

public class AddProductServlet extends HttpServlet {
    private final Products products;

    public AddProductServlet(Products db) {
        products = db;
    }

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws IOException {
        String name = request.getParameter("name");
        long price = Long.parseLong(request.getParameter("price"));

        try {
            products.addProduct(new Product(name, price));
        } catch (Exception e) {
            throw new RuntimeException(e);
        }


        response.setContentType("text/html");
        response.setStatus(HttpServletResponse.SC_OK);
        response.getWriter().println("OK");
    }
}
