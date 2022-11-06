package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.service.ProductsHtmlService;

import javax.servlet.http.HttpServletRequest;
import java.io.PrintWriter;

public class QueryServlet extends AbstractBaseProductServlet {

    public QueryServlet(ProductsHtmlService service) {
        super(service);
    }

    @Override
    protected void doRequest(HttpServletRequest request, PrintWriter respWriter) {
        var command = request.getParameter("command");
        switch (command) {
            case "max":
                htmlService.getMaxPriceProduct(respWriter);
                break;
            case "min":
                htmlService.getMinPriceProduct(respWriter);
                break;
            case "sum":
                htmlService.getSumPrice(respWriter);
                break;
            case "count":
                htmlService.getCount(respWriter);
                break;
            default:
                throw new IllegalArgumentException("Unknown command: " + command);
        }
    }
}
