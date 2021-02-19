
CREATE database test;
USE test;

DROP TABLE IF EXISTS employee;
CREATE TABLE employee ( 
`employee_id` INT(50) NOT NULL ,
`firstname` VARCHAR(100) NOT NULL ,
`lastname` VARCHAR(100) NOT NULL , 
`dob` DATETIME NOT NULL , 
`password` TEXT NOT NULL 
);