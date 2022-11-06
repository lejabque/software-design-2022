package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.service.ProductsHtmlService;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class QueryServlet extends HttpServlet {
    private final ProductsHtmlService htmlService;

    public QueryServlet(ProductsHtmlService service) {
        htmlService = service;
    }

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws IOException {
        String command = request.getParameter("command");

        var writer = response.getWriter();
        switch (command) {
            case "max":
                htmlService.getMaxPriceProduct(writer);
                break;
            case "min":
                htmlService.getMinPriceProduct(writer);
                break;
            case "sum":
                htmlService.getSumPrice(writer);
                break;
            case "count":
                htmlService.getCount(writer);
                break;
            default:
                writer.println("Unknown command: " + command);
                break;
        }

        response.setContentType("text/html");
        response.setStatus(HttpServletResponse.SC_OK);
    }

}
