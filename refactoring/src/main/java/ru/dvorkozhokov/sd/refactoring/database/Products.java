package ru.dvorkozhokov.sd.refactoring.database;

import ru.dvorkozhokov.sd.refactoring.models.Product;

import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.ArrayList;
import java.util.List;

public class Products {
    private static final String PRODUCTS_TABLE = "PRODUCT";
    private final Connection connection;

    public Products(Connection connection) throws SQLException {
        this.connection = connection;
        createTables();
    }

    public void addProduct(Product product) throws SQLException {
        String sql = String.format("INSERT INTO %s (NAME, PRICE) VALUES (\"%s\", %d)",
                PRODUCTS_TABLE, product.getName(), product.getPrice());
        Statement stmt = connection.createStatement();
        stmt.executeUpdate(sql);
        stmt.close();
    }


    public List<Product> getProducts() throws SQLException {
        Statement stmt = connection.createStatement();
        ResultSet rs = stmt.executeQuery(String.format("SELECT * FROM %s", PRODUCTS_TABLE));
        List<Product> products = new ArrayList<>();
        while (rs.next()) {
            String name = rs.getString("name");
            int price = rs.getInt("price");
            products.add(new Product(name, price));
        }
        rs.close();
        stmt.close();
        return products;
    }

    public Product getMaxPriceProduct() throws SQLException {
        Statement stmt = connection.createStatement();
        ResultSet rs = stmt.executeQuery(String.format("SELECT * FROM %s ORDER BY PRICE DESC LIMIT 1", PRODUCTS_TABLE));
        Product product = null;
        if (rs.next()) {
            String name = rs.getString("name");
            int price = rs.getInt("price");
            product = new Product(name, price);
        }
        rs.close();
        stmt.close();
        return product;
    }

    public Product getMinPriceProduct() throws SQLException {
        Statement stmt = connection.createStatement();
        ResultSet rs = stmt.executeQuery(String.format("SELECT * FROM %s ORDER BY PRICE LIMIT 1", PRODUCTS_TABLE));
        Product product = null;
        if (rs.next()) {
            String name = rs.getString("name");
            int price = rs.getInt("price");
            product = new Product(name, price);
        }
        rs.close();
        stmt.close();
        return product;
    }

    public int getSumPrice() throws SQLException {
        Statement stmt = connection.createStatement();
        ResultSet rs = stmt.executeQuery(String.format("SELECT SUM(price) FROM %s", PRODUCTS_TABLE));
        int sum = 0;
        if (rs.next()) {
            sum = rs.getInt(1);
        }
        rs.close();
        stmt.close();
        return sum;
    }

    public int getCount() throws SQLException {
        Statement stmt = connection.createStatement();
        ResultSet rs = stmt.executeQuery(String.format("SELECT COUNT(*) FROM %s", PRODUCTS_TABLE));
        int count = 0;
        if (rs.next()) {
            count = rs.getInt(1);
        }
        rs.close();
        stmt.close();
        return count;
    }

    public void createTables() throws SQLException {
        String sql = String.format("CREATE TABLE IF NOT EXISTS %s", PRODUCTS_TABLE) +
                "(ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
                " NAME           TEXT    NOT NULL, " +
                " PRICE          INT     NOT NULL)";
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
}
