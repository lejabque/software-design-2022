package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.service.ProductsHtmlService;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class GetProductsServlet extends HttpServlet {
    private final ProductsHtmlService htmlService;

    public GetProductsServlet(ProductsHtmlService service) {
        htmlService = service;
    }

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws IOException {
        htmlService.getProducts(response.getWriter());
        response.setContentType("text/html");
        response.setStatus(HttpServletResponse.SC_OK);
    }
}
