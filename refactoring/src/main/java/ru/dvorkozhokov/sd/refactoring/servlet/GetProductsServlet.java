package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.service.ProductsHtmlService;

import javax.servlet.http.HttpServletRequest;
import java.io.PrintWriter;

public class GetProductsServlet extends AbstractBaseProductServlet {

    public GetProductsServlet(ProductsHtmlService service) {
        super(service);
    }

    @Override
    protected void doRequest(HttpServletRequest request, PrintWriter respWriter) {
        getHtmlService().getProducts(respWriter);
    }
}
