CREATE TABLE directory (
    user_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email TEXT,
    phone_number TEXT
);

INSERT INTO directory (email, phone_number) VALUES
    ('homer@thesimpsons.com', '+1 888 123 1111'),
    ('marge@thesimpsons.com', '+1 888 123 1112'),
    ('maggie@thesimpsons.com', '+1 888 123 1113'),
    ('lisa@thesimpsons.com', '+1 888 123 1114'),
    ('bart@thesimpsons.com', '+1 888 123 1115');