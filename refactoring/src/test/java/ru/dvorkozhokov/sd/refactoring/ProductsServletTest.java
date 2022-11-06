package ru.dvorkozhokov.sd.refactoring;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import ru.dvorkozhokov.sd.refactoring.database.Products;
import ru.dvorkozhokov.sd.refactoring.service.ProductsHtmlService;
import ru.dvorkozhokov.sd.refactoring.servlet.AddProductServlet;
import ru.dvorkozhokov.sd.refactoring.servlet.GetProductsServlet;
import ru.dvorkozhokov.sd.refactoring.servlet.QueryServlet;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.util.Collections;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.*;


public class ProductsServletTest {
    private final Products products;
    private final AddProductServlet addServlet;
    private final GetProductsServlet getServlet;
    private final QueryServlet queryServlet;

    @BeforeEach
    public void resetTables() throws SQLException {
        products.resetTables();
    }

    public ProductsServletTest() {
        ProductsHtmlService service;
        try {
            var connection = DriverManager.getConnection("jdbc:sqlite:test.db");
            products = new Products(connection);
            service = new ProductsHtmlService(products);
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
        this.addServlet = new AddProductServlet(service);
        this.getServlet = new GetProductsServlet(service);
        this.queryServlet = new QueryServlet(service);
    }

    @Test
    public void emptyGetTest() throws ServletException, IOException {
        var resp = servletGet(getServlet, Collections.emptyMap(), 200);
        assertEquals(renderHtmlMap("", Map.of()), resp);
    }

    @Test
    public void simpleAddGetTest() throws ServletException, IOException {
        addProduct("iphone", 1000);

        var resp = servletGet(getServlet, Collections.emptyMap(), 200);
        assertEquals(renderHtmlMap("", Map.of("iphone", "1000")), resp);
    }

    @Test
    public void queryMaxTest() throws ServletException, IOException {
        addProduct("iphone", 1000);
        addProduct("samsung", 500);

        var resp = servletGet(queryServlet, Map.of("command", "max"), 200);
        assertEquals(renderHtmlMap("<h1>Product with max price: </h1>\n", Map.of("iphone", "1000")), resp);
    }

    @Test
    public void queryMinTest() throws ServletException, IOException {
        addProduct("iphone", 1000);
        addProduct("samsung", 500);

        var resp = servletGet(queryServlet, Map.of("command", "min"), 200);
        assertEquals(renderHtmlMap("<h1>Product with min price: </h1>\n", Map.of("samsung", "500")), resp);
    }

    @Test
    public void querySumTest() throws ServletException, IOException {
        addProduct("iphone", 1000);
        addProduct("samsung", 500);

        var resp = servletGet(queryServlet, Map.of("command", "sum"), 200);
        assertEquals("<html><body>\n" + "Summary price: \n1500\n" + "</body></html>\n", resp);
    }


    @Test
    public void queryCountTest() throws ServletException, IOException {
        addProduct("iphone", 1000);
        addProduct("samsung", 500);

        var resp = servletGet(queryServlet, Map.of("command", "count"), 200);
        assertEquals("<html><body>\n" + "Number of products: \n2\n" + "</body></html>\n", resp);
    }

    private String renderHtmlMap(String header, Map<String, String> map) {
        var sb = new StringBuilder();
        sb.append("<html><body>\n").append(header);
        for (var entry : map.entrySet()) {
            sb.append(entry.getKey()).append("\t").append(entry.getValue()).append("</br>\n");
        }
        sb.append("</body></html>\n");
        return sb.toString();
    }

    private void addProduct(String name, int price) throws ServletException, IOException {
        var resp = servletGet(addServlet, Map.of("name", name, "price", Integer.toString(price)), 200);
        assertEquals("OK\n", resp);
    }

    private String servletGet(HttpServlet servlet, Map<String, String> params, int expectedStatus) throws ServletException, IOException {
        var request = mock(HttpServletRequest.class);
        var response = mock(HttpServletResponse.class);
        when(request.getMethod()).thenReturn("GET");
        for (var entry : params.entrySet()) {
            when(request.getParameter(entry.getKey())).thenReturn(entry.getValue());
        }

        var respWriter = new StringWriter();
        when(response.getWriter()).thenReturn(new PrintWriter(respWriter));

        servlet.service(request, response);
        verify(response, times(1)).setStatus(expectedStatus);

        return respWriter.toString();
    }
}
