package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.service.ProductsHtmlService;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.io.PrintWriter;

public abstract class AbstractBaseProductServlet extends HttpServlet {
    protected final ProductsHtmlService htmlService;

    public AbstractBaseProductServlet(ProductsHtmlService service) {
        htmlService = service;
    }

    protected ProductsHtmlService getHtmlService() {
        return htmlService;
    }

    protected abstract void doRequest(HttpServletRequest request, PrintWriter respWriter);

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws IOException {
        var writer = response.getWriter();
        try {
            doRequest(request, writer);
            response.setStatus(HttpServletResponse.SC_OK);
        } catch (IllegalArgumentException e) {
            response.setStatus(HttpServletResponse.SC_BAD_REQUEST);
            writer.println(e.getMessage());
        } catch (Exception e) {
            response.setStatus(HttpServletResponse.SC_INTERNAL_SERVER_ERROR);
            writer.println(e.getMessage());
        }
        response.setContentType("text/html");
    }
}
