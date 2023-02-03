DROP TABLE IF EXISTS PersonEat;
CREATE TABLE PersonEat (
	PersonEat_id SERIAL PRIMARY KEY,
	name VARCHAR(255) DEFAULT NULL,
	descriptionText VARCHAR(255) DEFAULT NULL
);

DROP TABLE IF EXISTS PersonTeacher_Advisors;
CREATE TABLE PersonTeacher_Advisors (
	PersonTeacher_Advisors_id SERIAL PRIMARY KEY,
	PersonTeacher_id INT NOT NULL,
	advisors INT DEFAULT NULL
);

DROP TABLE IF EXISTS AddressBook_People;
CREATE TABLE AddressBook_People (
	AddressBook_People_id SERIAL PRIMARY KEY,
	AddressBook_id INT NOT NULL,
	people INT DEFAULT NULL
);

DROP TABLE IF EXISTS AddressBook;
CREATE TABLE AddressBook (
	AddressBook_id SERIAL PRIMARY KEY,
	title VARCHAR(255) DEFAULT NULL
);

DROP TABLE IF EXISTS Person_Phones;
CREATE TABLE Person_Phones (
	Person_Phones_id SERIAL PRIMARY KEY,
	Person_id INT NOT NULL,
	phones INT DEFAULT NULL
);

DROP TABLE IF EXISTS Person_Child;
CREATE TABLE Person_Child (
	Person_Child_id SERIAL PRIMARY KEY,
	Person_id INT NOT NULL,
	child VARCHAR(255)
);

DROP TABLE IF EXISTS Person;
CREATE TABLE Person (
	Person_id SERIAL PRIMARY KEY,
	name VARCHAR(255) DEFAULT NULL,
	id BIGINT DEFAULT NULL,
	email VARCHAR(255) DEFAULT NULL,
	partner INT DEFAULT NULL,
	breakfest INT DEFAULT NULL,
	lunch INT DEFAULT NULL,
	dinner INT DEFAULT NULL,
	primarySchool INT DEFAULT NULL,
	middleSchool INT DEFAULT NULL,
	highSchool INT DEFAULT NULL,
	vehicle INT DEFAULT NULL,
	bus VARCHAR(255) DEFAULT NULL,
	dadmom INT DEFAULT NULL,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS PersonParent;
CREATE TABLE PersonParent (
	PersonParent_id SERIAL PRIMARY KEY,
	father INT DEFAULT NULL,
	mother INT DEFAULT NULL,
	city VARCHAR(255) DEFAULT NULL
);

DROP TABLE IF EXISTS Person_Skills;
CREATE TABLE Person_Skills (
	Person_Skills_id SERIAL PRIMARY KEY,
	Person_id INT NOT NULL,
	skills INT DEFAULT NULL
);

DROP TABLE IF EXISTS PersonSpouse;
CREATE TABLE PersonSpouse (
	PersonSpouse_id SERIAL PRIMARY KEY,
	nameSpouse VARCHAR(255) DEFAULT NULL,
	emailSpouse VARCHAR(255) DEFAULT NULL
);

DROP TABLE IF EXISTS PersonDegreesEntry;
CREATE TABLE PersonDegreesEntry (
	PersonDegreesEntry_id SERIAL PRIMARY KEY,
	keyEscape VARCHAR(255) DEFAULT NULL,
	value INT DEFAULT NULL
);

DROP TABLE IF EXISTS PersonCar;
CREATE TABLE PersonCar (
	PersonCar_id SERIAL PRIMARY KEY,
	maker VARCHAR(255) DEFAULT NULL,
	model VARCHAR(255) DEFAULT NULL,
	year BIGINT DEFAULT NULL
);

DROP TABLE IF EXISTS Oneofs;
CREATE TABLE Oneofs (
	Oneofs_id SERIAL PRIMARY KEY,
	columnNames VARCHAR(255));
INSERT INTO Oneofs(tableName, oneof, columnNames) VALUES ('Person', 'Transport', 'vehicle,taxi,bus');

DROP TABLE IF EXISTS Person_Snacks;
CREATE TABLE Person_Snacks (
	Person_Snacks_id SERIAL PRIMARY KEY,
	Person_id INT NOT NULL,
	snacks INT DEFAULT NULL
);

DROP TABLE IF EXISTS PersonPhoneNumber;
CREATE TABLE PersonPhoneNumber (
	PersonPhoneNumber_id SERIAL PRIMARY KEY,
	numberEscape VARCHAR(255) DEFAULT NULL,
	typeEscape PersonPhoneNumber_typeEscape_Type DEFAULT NULL
);

DROP TABLE IF EXISTS PersonSkillsEntry;
CREATE TABLE PersonSkillsEntry (
	PersonSkillsEntry_id SERIAL PRIMARY KEY,
	keyEscape VARCHAR(255) DEFAULT NULL,
	value VARCHAR(255) DEFAULT NULL
);

DROP TABLE IF EXISTS Person_Degrees;
CREATE TABLE Person_Degrees (
	Person_Degrees_id SERIAL PRIMARY KEY,
	Person_id INT NOT NULL,
	degrees INT DEFAULT NULL
);

DROP TABLE IF EXISTS PersonEducation;
CREATE TABLE PersonEducation (
	PersonEducation_id SERIAL PRIMARY KEY,
	year BIGINT DEFAULT NULL,
	school VARCHAR(255) DEFAULT NULL,
	name PersonEducation_name_Type DEFAULT NULL,
	major VARCHAR(255) DEFAULT NULL
);

DROP TABLE IF EXISTS PersonTeacher;
CREATE TABLE PersonTeacher (
	PersonTeacher_id SERIAL PRIMARY KEY,
	fullname VARCHAR(255) DEFAULT NULL,
	school VARCHAR(255) DEFAULT NULL,
	startYear BIGINT DEFAULT NULL,
	endYear BIGINT DEFAULT NULL
);

ALTER TABLE PersonParent ADD CONSTRAINT constraint_PersonParent_father_Person_Person_id FOREIGN KEY (father) REFERENCES Person (Person_id);
ALTER TABLE AddressBook_People ADD CONSTRAINT constraint_AddressBook_People_people_Person_Person_id FOREIGN KEY (people) REFERENCES Person (Person_id);
ALTER TABLE Person_Child ADD CONSTRAINT constraint_Person_Child_Person_id_Person_Person_id FOREIGN KEY (Person_id) REFERENCES Person (Person_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_partner_PersonSpouse_PersonSpouse_id FOREIGN KEY (partner) REFERENCES PersonSpouse (PersonSpouse_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_middleSchool_PersonTeacher_PersonTeacher_id FOREIGN KEY (middleSchool) REFERENCES PersonTeacher (PersonTeacher_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_vehicle_PersonCar_PersonCar_id FOREIGN KEY (vehicle) REFERENCES PersonCar (PersonCar_id);
ALTER TABLE PersonTeacher_Advisors ADD CONSTRAINT constraint_PersonTeacher_Advisors_advisors_PersonTeacher_PersonTeacher_id FOREIGN KEY (advisors) REFERENCES PersonTeacher (PersonTeacher_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_highSchool_PersonTeacher_PersonTeacher_id FOREIGN KEY (highSchool) REFERENCES PersonTeacher (PersonTeacher_id);
ALTER TABLE PersonTeacher_Advisors ADD CONSTRAINT constraint_PersonTeacher_Advisors_PersonTeacher_id_PersonTeacher_PersonTeacher_id FOREIGN KEY (PersonTeacher_id) REFERENCES PersonTeacher (PersonTeacher_id);
ALTER TABLE AddressBook_People ADD CONSTRAINT constraint_AddressBook_People_AddressBook_id_AddressBook_AddressBook_id FOREIGN KEY (AddressBook_id) REFERENCES AddressBook (AddressBook_id);
ALTER TABLE Person_Phones ADD CONSTRAINT constraint_Person_Phones_Person_id_Person_Person_id FOREIGN KEY (Person_id) REFERENCES Person (Person_id);
ALTER TABLE Person_Phones ADD CONSTRAINT constraint_Person_Phones_phones_PersonPhoneNumber_PersonPhoneNumber_id FOREIGN KEY (phones) REFERENCES PersonPhoneNumber (PersonPhoneNumber_id);
ALTER TABLE Person_Skills ADD CONSTRAINT constraint_Person_Skills_Person_id_Person_Person_id FOREIGN KEY (Person_id) REFERENCES Person (Person_id);
ALTER TABLE Person_Degrees ADD CONSTRAINT constraint_Person_Degrees_degrees_PersonDegreesEntry_PersonDegreesEntry_id FOREIGN KEY (degrees) REFERENCES PersonDegreesEntry (PersonDegreesEntry_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_primarySchool_PersonTeacher_PersonTeacher_id FOREIGN KEY (primarySchool) REFERENCES PersonTeacher (PersonTeacher_id);
ALTER TABLE PersonDegreesEntry ADD CONSTRAINT constraint_PersonDegreesEntry_value_PersonEducation_PersonEducation_id FOREIGN KEY (value) REFERENCES PersonEducation (PersonEducation_id);
ALTER TABLE Person_Skills ADD CONSTRAINT constraint_Person_Skills_skills_PersonSkillsEntry_PersonSkillsEntry_id FOREIGN KEY (skills) REFERENCES PersonSkillsEntry (PersonSkillsEntry_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_breakfest_PersonEat_PersonEat_id FOREIGN KEY (breakfest) REFERENCES PersonEat (PersonEat_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_lunch_PersonEat_PersonEat_id FOREIGN KEY (lunch) REFERENCES PersonEat (PersonEat_id);
ALTER TABLE Person_Snacks ADD CONSTRAINT constraint_Person_Snacks_Person_id_Person_Person_id FOREIGN KEY (Person_id) REFERENCES Person (Person_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_dadmom_PersonParent_PersonParent_id FOREIGN KEY (dadmom) REFERENCES PersonParent (PersonParent_id);
ALTER TABLE Person_Degrees ADD CONSTRAINT constraint_Person_Degrees_Person_id_Person_Person_id FOREIGN KEY (Person_id) REFERENCES Person (Person_id);
ALTER TABLE Person ADD CONSTRAINT constraint_Person_dinner_PersonEat_PersonEat_id FOREIGN KEY (dinner) REFERENCES PersonEat (PersonEat_id);
ALTER TABLE Person_Snacks ADD CONSTRAINT constraint_Person_Snacks_snacks_PersonEat_PersonEat_id FOREIGN KEY (snacks) REFERENCES PersonEat (PersonEat_id);
ALTER TABLE PersonParent ADD CONSTRAINT constraint_PersonParent_mother_Person_Person_id FOREIGN KEY (mother) REFERENCES Person (Person_id);
