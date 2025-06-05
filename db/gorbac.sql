-- 🔓 Desactivar chequeo de claves foráneas
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS roles_apis;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS apis;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;

-- (resto del script aquí...)

-- 🔒 Reactivar chequeo de claves foráneas
SET FOREIGN_KEY_CHECKS = 1;



-- Tabla de usuarios
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user VARCHAR(120) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(120),
  last_name VARCHAR(120),
  password VARCHAR(255),
  provider VARCHAR(32) DEFAULT 'local',
  provider_id VARCHAR(255),
  last_login DATETIME DEFAULT NULL,
  active BOOLEAN DEFAULT TRUE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT NULL,
  deleted_at DATETIME DEFAULT NULL
);

-- Tabla de roles
CREATE TABLE roles (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT NULL,
  deleted_at DATETIME DEFAULT NULL
);

-- Tabla de APIs (servicios/permisos)
CREATE TABLE apis (
  id INT AUTO_INCREMENT PRIMARY KEY,
  endpoint VARCHAR(100) NOT NULL UNIQUE,
  description VARCHAR(100),
  hidden BOOLEAN DEFAULT FALSE,
  public BOOLEAN DEFAULT FALSE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT NULL,
  deleted_at DATETIME DEFAULT NULL
);

-- Relación muchos a muchos: usuarios ↔ roles
CREATE TABLE user_roles (
  user_id INT NOT NULL,
  role_id INT NOT NULL,
  PRIMARY KEY (user_id, role_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Relación muchos a muchos: roles ↔ apis
CREATE TABLE roles_apis (
  role_id INT NOT NULL,
  api_id INT NOT NULL,
  PRIMARY KEY (role_id, api_id),
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (api_id) REFERENCES apis(id) ON DELETE CASCADE
);

-- Inserción de roles
INSERT INTO roles (name) VALUES ('admin'), ('user');

-- Inserción de apis
INSERT INTO apis (endpoint, description, hidden, public) VALUES
('/users/index', 'Gestión de usuarios', FALSE, FALSE),
('/roles/index', 'Gestión de roles', FALSE, FALSE),
('/products/index', 'Gestión de productos', FALSE, TRUE);

-- Inserción de usuarios
INSERT INTO users (user, email, name, last_name, password, provider) VALUES
('seba', 'seba@email.com', 'Sebas', 'Apellido', '$2a$14$tqBEgMKSgxRJE7.SSd7Nwe5r5sbhqZd09/HcDBxICOI53SHyHxrSm', 'local'),
('juan', 'juan@email.com', 'Juan', 'Pérez', '$2a$14$Q9QwQwQwQwQwQwQwQwQwQeQwQwQwQwQwQwQwQwQwQwQwQwQwQw', 'local'),
('ana', 'ana@email.com', 'Ana', 'García', '$2a$14$Q9QwQwQwQwQwQwQwQwQwQeQwQwQwQwQwQwQwQwQwQwQwQwQwQw', 'local');

-- Asignación de roles a usuarios
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1),
(2, 2),
(3, 2);

-- Asignación de apis a roles
INSERT INTO roles_apis (role_id, api_id) VALUES
(1, 1), (1, 2), (1, 3);

-- Tabla de categorías
CREATE TABLE categories (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT NULL,
  deleted_at DATETIME DEFAULT NULL
);

-- Tabla de productos
CREATE TABLE products (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10,2) NOT NULL,
  stock INT DEFAULT 0,
  category_id INT,
  image VARCHAR(255),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT NULL,
  deleted_at DATETIME DEFAULT NULL,
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- Tabla de carritos
CREATE TABLE carts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT NULL,
  deleted_at DATETIME DEFAULT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Tabla de ítems en el carrito
CREATE TABLE cart_items (
  id INT AUTO_INCREMENT PRIMARY KEY,
  cart_id INT NOT NULL,
  product_id INT NOT NULL,
  quantity INT NOT NULL DEFAULT 1,
  price DECIMAL(10,2) NOT NULL,
  image VARCHAR(255),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT NULL,
  deleted_at DATETIME DEFAULT NULL,
  FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Inserción de categorías
INSERT INTO categories (name, description) VALUES
('Cucha', 'Cuchas y camas para mascotas'),
('Funda', 'Fundas y accesorios para mascotas');

-- Inserción de productos
INSERT INTO products (name, description, price, stock, category_id, image) VALUES
('Cucha Deluxe Clásica', 'Diseñada para brindar confort y elegancia, ideal para mascotas medianas y grandes. Con base antideslizante y acolchado extra. Fabricada con tela impermeable de alta resistencia, interior de espuma viscoelástica y costuras reforzadas.', 97.00, 86, 1, '/img/cuchas/cucha1.jpg'),
('Cucha Eco Confort', 'Fabricada con materiales reciclados, liviana y resistente. Ideal para climas templados y uso en interiores. Estructura de plástico reciclado, funda de algodón orgánico y relleno de fibras sintéticas hipoalergénicas.', 5.00, 52, 1, '/img/cuchas/cucha2.jpg'),
('Cucha Térmica', 'Protege del frío y la humedad. Revestida con materiales térmicos, perfecta para invierno y espacios exteriores. Exterior de poliéster impermeable, interior de lana sintética y base aislante de goma EVA.', 74.00, 23, 1, '/img/cuchas/cucha3.jpg'),
('Cucha Modular Urban', 'Moderna, desmontable y fácil de transportar. Su diseño urbano combina con cualquier ambiente del hogar. Paneles de polipropileno, uniones de silicona flexible y funda de microfibra lavable.', 47.00, 87, 1, '/img/cuchas/cucha4.jpg'),
('Funda Estilo Campo', 'Textura rústica y resistente. Ideal para ambientes rurales o mascotas aventureras que disfrutan del aire libre. Confeccionada en lona gruesa de algodón y costuras dobles para mayor durabilidad.', 48.00, 33, 2, '/img/fundas/funda2.jpg'),
('Funda Ultra Soft', 'Máxima suavidad para el descanso de tu mascota. Lavable, antialérgica y disponible en varios colores. Tejido exterior de microfibra ultrasuave y relleno de vellón siliconado.', 97.00, 9, 2, '/img/fundas/funda1.jpg'),
('Funda Soft', 'Práctica, cómoda y acolchada. Ideal para usar sobre colchones o dentro de cuchas rígidas. Exterior de algodón peinado y relleno de espuma de poliuretano.', 97.00, 6, 2, '/img/fundas/funda3.jpg'),
('Funda Doble', 'Funda reversible con doble cara: una térmica para invierno y otra fresca para verano. ¡2 en 1! Cara térmica de polar y cara fresca de algodón, con relleno de fibra hueca siliconada.', 97.00, 6, 2, '/img/fundas/funda4.jpg');

-- Inserción de carritos
INSERT INTO carts (user_id) VALUES
(2),
(3);

-- Inserción de ítems en carritos
INSERT INTO cart_items (cart_id, product_id, quantity, price) VALUES
(1, 1, 1, 850.00),
(1, 2, 2, 50.00),
(2, 3, 1, 15.00);

