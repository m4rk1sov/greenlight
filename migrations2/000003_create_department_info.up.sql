CREATE TABLE IF NOT EXISTS departmentInfo (
                                               id BIGSERIAL PRIMARY KEY,
                                               departmentName VARCHAR(255) NOT NULL,
                                               staffQuantity INT NOT NULL,
                                               departmentDirector VARCHAR(255) NOT NULL,
                                               module_Info INT REFERENCES module_info(id)
);