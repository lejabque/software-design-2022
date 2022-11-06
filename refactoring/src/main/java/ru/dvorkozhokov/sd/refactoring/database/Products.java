package ru.dvorkozhokov.sd.refactoring.database;

import ru.dvorkozhokov.sd.refactoring.models.Product;

import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Function;

public class Products {
    private static final String PRODUCTS_TABLE = "PRODUCT";
    private final Connection connection;

    public Products(Connection connection) throws SQLException {
        this.connection = connection;
        createTables();
    }

    public void addProduct(Product product) {
        String query = String.format("INSERT INTO %s (NAME, PRICE) VALUES (\"%s\", %d)", PRODUCTS_TABLE, product.getName(), product.getPrice());
        executeUpdateQuery(query);
    }

    public List<Product> getProducts() {
        String query = String.format("SELECT * FROM %s", PRODUCTS_TABLE);
        return executeReadQuery(query, this::parseProducts);
    }

    public Product getMinPriceProduct() {
        String query = String.format("SELECT * FROM %s ORDER BY PRICE LIMIT 1", PRODUCTS_TABLE);
        return executeReadQuery(query, this::fetchProduct);
    }

    public Product getMaxPriceProduct() {
        String query = String.format("SELECT * FROM %s ORDER BY PRICE DESC LIMIT 1", PRODUCTS_TABLE);
        return executeReadQuery(query, this::fetchProduct);
    }

    public int getSumPrice() throws SQLException {
        String query = String.format("SELECT SUM(price) FROM %s", PRODUCTS_TABLE);
        return executeReadQuery(query, this::parseInt);
    }

    public int getCount() throws SQLException {
        String query = String.format("SELECT COUNT(*) FROM %s", PRODUCTS_TABLE);
        return executeReadQuery(query, this::parseInt);
    }

    public void createTables() throws SQLException {
        String sql = String.format("CREATE TABLE IF NOT EXISTS %s", PRODUCTS_TABLE) + "(ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," + " NAME           TEXT    NOT NULL, " + " PRICE          INT     NOT NULL)";
        Statement stmt = connection.createStatement();
        stmt.executeUpdate(sql);
        stmt.close();
    }

    public void resetTables() throws SQLException {
        String sql = String.format("DROP TABLE IF EXISTS %s", PRODUCTS_TABLE);
        Statement stmt = connection.createStatement();
        stmt.executeUpdate(sql);
        stmt.close();
        createTables();
    }

    private Product fetchProduct(ResultSet rs) {
        var products = parseProducts(rs);
        assert products.size() <= 1;
        return products.size() == 1 ? products.get(0) : null;
    }

    private List<Product> parseProducts(ResultSet rs) {
        List<Product> products = new ArrayList<>();
        try {
            while (rs.next()) {
                String name = rs.getString("name");
                int price = rs.getInt("price");
                products.add(new Product(name, price));
            }
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
        return products;
    }

    private int parseInt(ResultSet rs) {
        try {
            if (rs.next()) {
                return rs.getInt(1);
            }
            throw new RuntimeException("No result");
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }

    private <R> R executeReadQuery(String query, Function<ResultSet, R> resultParser) {
        try {
            Statement stmt = connection.createStatement();
            ResultSet rs = stmt.executeQuery(query);
            var res = resultParser.apply(rs);
            rs.close();
            stmt.close();
            return res;
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }

    private void executeUpdateQuery(String query) {
        try {
            Statement stmt = connection.createStatement();
            stmt.executeUpdate(query);
            stmt.close();
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }
}
