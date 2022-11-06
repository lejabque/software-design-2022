package ru.dvorkozhokov.sd.refactoring.service;

import ru.dvorkozhokov.sd.refactoring.database.Products;
import ru.dvorkozhokov.sd.refactoring.models.Product;

import java.io.PrintWriter;
import java.sql.SQLException;

public class ProductsHtmlService {
    private final Products products;

    public ProductsHtmlService(Products db) {
        products = db;
    }

    public void addProduct(Product product) {
        try {
            products.addProduct(product);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public void getProducts(PrintWriter writer) {
        try {
            var res = products.getProducts();
            var sb = new StringBuilder();
            for (var product : res) {
                sb.append(product.getName()).append("\t").append(product.getPrice()).append("</br>");
            }
            writeHtmlBody(writer, sb);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public void getMaxPriceProduct(PrintWriter writer) {
        Product res;
        try {
            res = products.getMaxPriceProduct();
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
        writeHtmlBody(writer, writeProduct(new StringBuilder().append("<h1>Product with max price: </h1>\n"),
                res).append("</br>"));
    }

    public void getMinPriceProduct(PrintWriter writer) {
        Product res;
        try {
            res = products.getMinPriceProduct();
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
        writeHtmlBody(writer, writeProduct(new StringBuilder().append("<h1>Product with min price: </h1>\n"),
                res).append("</br>"));
    }

    public void getSumPrice(PrintWriter writer) {
        int sum;
        try {
            sum = products.getSumPrice();
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
        writeHtmlBody(writer, new StringBuilder().append("Summary price: \n").append(sum));
    }

    public void getCount(PrintWriter writer) {
        int cnt;
        try {
            cnt = products.getCount();
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
        writeHtmlBody(writer, new StringBuilder().append("Number of products: \n").append(cnt));
    }

    private StringBuilder writeProduct(StringBuilder sb, Product product) {
        sb.append(product.getName()).append("\t").append(product.getPrice());
        return sb;
    }

    private void writeHtmlBody(PrintWriter writer, StringBuilder body) {
        writer.println("<html><body>");
        if (body.length() > 0) {
            writer.println(body);
        }
        writer.println("</body></html>");
    }
}
