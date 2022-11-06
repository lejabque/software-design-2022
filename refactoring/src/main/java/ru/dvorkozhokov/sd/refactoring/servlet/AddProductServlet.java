package ru.dvorkozhokov.sd.refactoring.servlet;

import ru.dvorkozhokov.sd.refactoring.models.Product;
import ru.dvorkozhokov.sd.refactoring.service.ProductsHtmlService;

import javax.servlet.http.HttpServletRequest;
import java.io.PrintWriter;

public class AddProductServlet extends AbstractBaseProductServlet {

    public AddProductServlet(ProductsHtmlService service) {
        super(service);
    }

    @Override
    protected void doRequest(HttpServletRequest request, PrintWriter respWriter) {
        var name = request.getParameter("name");
        var price = request.getParameter("price");
        if (name == null || price == null) {
            throw new IllegalArgumentException("Name and price are required");
        }
        long parsedPrice;
        try {
            parsedPrice = Long.parseLong(price);
        } catch (NumberFormatException e) {
            throw new IllegalArgumentException("Invalid price: " + price);
        }
        getHtmlService().addProduct(new Product(name, parsedPrice), respWriter);
    }
}
