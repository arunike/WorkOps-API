CREATE TABLE IF NOT EXISTS Offices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Departments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Associates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    title VARCHAR(255),
    department VARCHAR(255),
    office VARCHAR(255),
    status VARCHAR(50),
    start_date DATETIME,
    empl_status VARCHAR(50),
    salary INT,
    dob DATETIME,
    profile_picture VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    phone_number VARCHAR(50),
    gender VARCHAR(50),
    private_email VARCHAR(255),
    manager_id INT,
    FOREIGN KEY (manager_id) REFERENCES Associates(id)
);

CREATE TABLE IF NOT EXISTS Tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    requester_id INT NOT NULL,
    task_name VARCHAR(255) NOT NULL,
    task_value VARCHAR(255) NOT NULL,
    reason TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    target_value INT,
    approvers JSON,
    timestamp INT,
    comments TEXT
);

CREATE TABLE IF NOT EXISTS time_off_requests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    associate_id INT NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    reason TEXT,
    approver_id INT,
    status VARCHAR(50) DEFAULT 'Pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (associate_id) REFERENCES Associates(id),
    FOREIGN KEY (approver_id) REFERENCES Associates(id)
);

CREATE TABLE IF NOT EXISTS menu_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    menu_item VARCHAR(100) NOT NULL,
    permission_type VARCHAR(50) NOT NULL,
    permission_value VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_permission (menu_item, permission_type, permission_value)
);


CREATE TABLE IF NOT EXISTS Thanks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_id INT,
    to_id INT,
    message TEXT,
    category VARCHAR(255),
    timestamp BIGINT,
    FOREIGN KEY (from_id) REFERENCES Associates(id),
    FOREIGN KEY (to_id) REFERENCES Associates(id)
);

CREATE TABLE IF NOT EXISTS thanks_likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    thank_id INT NOT NULL,
    associate_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_like (thank_id, associate_id),
    FOREIGN KEY (thank_id) REFERENCES Thanks(id) ON DELETE CASCADE,
    FOREIGN KEY (associate_id) REFERENCES Associates(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS thanks_comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    thank_id INT NOT NULL,
    associate_id INT NOT NULL,
    comment TEXT NOT NULL,
    timestamp BIGINT NOT NULL,
    FOREIGN KEY (thank_id) REFERENCES Thanks(id) ON DELETE CASCADE,
    FOREIGN KEY (associate_id) REFERENCES Associates(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS TimeOffRequests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    associate_id INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    reason TEXT,
    status VARCHAR(50) DEFAULT 'Pending',
    FOREIGN KEY (associate_id) REFERENCES Associates(id)
);

CREATE TABLE IF NOT EXISTS DocumentCategories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS AppSettings (
    setting_key VARCHAR(255) PRIMARY KEY,
    setting_value TEXT
);

INSERT IGNORE INTO Offices (name) VALUES ('London'), ('New York'), ('Paris');
INSERT IGNORE INTO Departments (name) VALUES ('Design'), ('Sales'), ('Tech'), ('People'), ('Product'), ('Legal'), ('Finance'), ('Executive');
INSERT IGNORE INTO DocumentCategories (name) VALUES ('Contracts'), ('Identity'), ('Tax'), ('Payroll'), ('Other');

INSERT INTO Associates (id, first_name, last_name, title, department, office, status, start_date, empl_status, salary, dob, manager_id, email, password, phone_number, gender, private_email) VALUES
(1, 'Richie', 'Zhou', 'CEO', 'Executive', 'New York', 'Employed', '2015-01-01', 'Employed', 250000, '1975-05-20', NULL, 'richie.zhou@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0001', 'Male', 'richie.zhou@personal.com'),
(2, 'Melissa', 'Smith', 'Head of Design', 'Design', 'London', 'Employed', '2018-10-04', 'Employed', 100954, '1984-03-09', 1, 'melissa.smith@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0002', 'Female', 'melissa.smith@personal.com'),
(3, 'Anna', 'Jones', 'Head of Sales', 'Sales', 'New York', 'Employed', '2018-01-25', 'Employed', 134335, '1982-11-18', 1, 'anna.jones@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0003', 'Female', 'anna.jones@personal.com'),
(4, 'Jonathan', 'Lewis', 'Head of Tech', 'Tech', 'New York', 'Employed', '2018-01-12', 'Employed', 123284, '1988-01-30', 1, 'jonathan.lewis@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0004', 'Male', 'jonathan.lewis@personal.com'),
(5, 'Mia', 'Mao', 'Head of People', 'People', 'New York', 'Employed', '2017-08-28', 'Employed', 101081, '1984-09-08', 1, 'mia.mao@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0005', 'Female', 'mia.mao@personal.com'),
(6, 'Michelle', 'Thomas', 'Head of Product', 'Product', 'Paris', 'Employed', '2018-07-06', 'Employed', 117451, '1990-10-07', 1, 'michelle.thomas@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0006', 'Female', 'michelle.thomas@personal.com'),
(7, 'Rebecca', 'Young', 'Head of Legal', 'Legal', 'Paris', 'Employed', '2019-09-28', 'Employed', 118162, '1990-06-01', 1, 'rebecca.young@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0007', 'Female', 'rebecca.young@personal.com'),
(8, 'Donna', 'Smith', 'Head of Finance', 'Finance', 'London', 'Employed', '2018-09-08', 'Employed', 116867, '1981-08-27', 1, 'donna.smith@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0008', 'Female', 'donna.smith@personal.com'),
(9, 'Charles', 'Rodriguez', 'QA Engineer', 'Tech', 'London', 'Employed', '2022-06-05', 'Employed', 57333, '1995-12-08', 4, 'charles.rodriguez@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0009', 'Male', 'charles.rodriguez@personal.com'),
(10, 'Carol', 'Jackson', 'Compliance Officer', 'Legal', 'London', 'Employed', '2025-10-27', 'Employed', 72813, '1998-01-08', 7, 'carol.jackson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0010', 'Female', 'carol.jackson@personal.com'),
(11, 'Timothy', 'Davis', 'Full Stack Developer', 'Tech', 'London', 'Employed', '2021-11-20', 'Employed', 81371, '2000-03-09', 4, 'timothy.davis@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0011', 'Male', 'timothy.davis@personal.com'),
(12, 'Rebecca', 'King', 'Compliance Officer', 'Legal', 'New York', 'Employed', '2024-07-20', 'Employed', 83662, '1992-12-19', 7, 'rebecca.king@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0012', 'Female', 'rebecca.king@personal.com'),
(13, 'Deborah', 'Miller', 'Financial Analyst', 'Finance', 'London', 'Employed', '2021-04-25', 'Employed', 64920, '2000-04-26', 8, 'deborah.miller@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0013', 'Female', 'deborah.miller@personal.com'),
(14, 'Scott', 'Martinez', 'Software Engineer', 'Tech', 'New York', 'Employed', '2021-08-20', 'Employed', 50738, '1998-10-16', 4, 'scott.martinez@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0014', 'Male', 'scott.martinez@personal.com'),
(15, 'Patricia', 'Johnson', 'QA Engineer', 'Tech', 'Paris', 'Employed', '2022-12-13', 'Employed', 92288, '1996-09-05', 4, 'patricia.johnson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0015', 'Female', 'patricia.johnson@personal.com'),
(16, 'Shirley', 'White', 'Controller', 'Finance', 'Paris', 'Employed', '2025-09-10', 'Employed', 94624, '2000-04-02', 8, 'shirley.white@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0016', 'Female', 'shirley.white@personal.com'),
(17, 'Matthew', 'Jackson', 'UI/UX Designer', 'Design', 'London', 'Employed', '2021-11-01', 'Employed', 71709, '1992-05-16', 2, 'matthew.jackson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0017', 'Male', 'matthew.jackson@personal.com'),
(18, 'Justin', 'Green', 'Account Executive', 'Sales', 'London', 'Employed', '2025-03-22', 'Employed', 63187, '1999-05-14', 3, 'justin.green@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0018', 'Male', 'justin.green@personal.com'),
(19, 'Jonathan', 'Davis', 'Accountant', 'Finance', 'New York', 'Employed', '2023-05-16', 'Employed', 81688, '1990-01-28', 8, 'jonathan.davis@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0019', 'Male', 'jonathan.davis@personal.com'),
(20, 'Emma', 'Mitchell', 'Financial Analyst', 'Finance', 'New York', 'Employed', '2023-06-30', 'Employed', 65973, '1999-05-22', 8, 'emma.mitchell@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0020', 'Female', 'emma.mitchell@personal.com'),
(21, 'Ronald', 'Wilson', 'Designer', 'Design', 'London', 'Employed', '2025-07-06', 'Employed', 93629, '1997-06-19', 2, 'ronald.wilson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0021', 'Male', 'ronald.wilson@personal.com'),
(22, 'Carol', 'Jackson', 'Business Analyst', 'Product', 'New York', 'Employed', '2021-01-04', 'Employed', 78198, '1993-12-23', 6, 'carol.jackson2@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0022', 'Female', 'carol.jackson@personal.com'),
(23, 'Nicole', 'Hall', 'Product Owner', 'Product', 'Paris', 'Employed', '2024-08-23', 'Employed', 80145, '1996-06-10', 6, 'nicole.hall@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0023', 'Female', 'nicole.hall@personal.com'),
(24, 'Mary', 'Wilson', 'Business Analyst', 'Product', 'New York', 'Employed', '2024-09-29', 'Employed', 78276, '1993-05-31', 6, 'mary.wilson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0024', 'Female', 'mary.wilson@personal.com'),
(25, 'Ryan', 'Nguyen', 'Product Manager', 'Product', 'London', 'Employed', '2024-02-03', 'Employed', 71906, '1997-12-22', 6, 'ryan.nguyen@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0025', 'Male', 'ryan.nguyen@personal.com'),
(26, 'Jason', 'Brown', 'Full Stack Developer', 'Tech', 'London', 'Employed', '2022-05-02', 'Employed', 82202, '1996-01-29', 4, 'jason.brown@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0026', 'Male', 'jason.brown@personal.com'),
(27, 'Carol', 'Davis', 'Accountant', 'Finance', 'Paris', 'Employed', '2024-11-18', 'Employed', 82166, '1999-06-17', 8, 'carol.davis@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0027', 'Female', 'carol.davis@personal.com'),
(28, 'Rebecca', 'Flores', 'Talent Acquisition', 'People', 'New York', 'Employed', '2025-12-14', 'Employed', 49464, '1998-06-14', 5, 'rebecca.flores@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0028', 'Female', 'rebecca.flores@personal.com'),
(29, 'Kenneth', 'Scott', 'Legal Counsel', 'Legal', 'London', 'Employed', '2025-07-14', 'Employed', 48734, '1993-02-12', 7, 'kenneth.scott@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0029', 'Male', 'kenneth.scott@personal.com'),
(30, 'Laura', 'Clark', 'Financial Analyst', 'Finance', 'Paris', 'Employed', '2022-05-02', 'Employed', 86982, '1995-01-26', 8, 'laura.clark@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0030', 'Female', 'laura.clark@personal.com'),
(31, 'Stephen', 'Garcia', 'Software Engineer', 'Tech', 'London', 'Employed', '2024-07-14', 'Employed', 51707, '1996-07-02', 4, 'stephen.garcia@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0031', 'Male', 'stephen.garcia@personal.com'),
(32, 'Rebecca', 'Taylor', 'Senior Designer', 'Design', 'London', 'Employed', '2024-11-06', 'Employed', 59085, '1998-11-22', 2, 'rebecca.taylor@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0032', 'Female', 'rebecca.taylor@personal.com'),
(33, 'Amy', 'Allen', 'Account Executive', 'Sales', 'New York', 'Employed', '2022-03-05', 'Employed', 90121, '1993-08-14', 3, 'amy.allen@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0033', 'Female', 'amy.allen@personal.com'),
(34, 'Matthew', 'Brown', 'UI/UX Designer', 'Design', 'New York', 'Employed', '2022-02-11', 'Employed', 45866, '1996-11-20', 2, 'matthew.brown@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0034', 'Male', 'matthew.brown@personal.com'),
(35, 'Sandra', 'Flores', 'Software Engineer', 'Tech', 'Paris', 'Employed', '2023-12-10', 'Employed', 61959, '1999-01-26', 4, 'sandra.flores@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0035', 'Female', 'sandra.flores@personal.com'),
(36, 'Linda', 'Rodriguez', 'UI/UX Designer', 'Design', 'London', 'Employed', '2022-02-12', 'Employed', 76400, '1998-07-19', 2, 'linda.rodriguez@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0036', 'Female', 'linda.rodriguez@personal.com'),
(37, 'Edward', 'Flores', 'Accountant', 'Finance', 'Paris', 'Employed', '2022-01-10', 'Employed', 63620, '1992-02-12', 8, 'edward.flores@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0037', 'Male', 'edward.flores@personal.com'),
(38, 'Sandra', 'White', 'Senior Designer', 'Design', 'New York', 'Employed', '2024-09-05', 'Employed', 50429, '1990-06-26', 2, 'sandra.white@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0038', 'Female', 'sandra.white@personal.com'),
(39, 'Sharon', 'Wilson', 'DevOps Engineer', 'Tech', 'New York', 'Employed', '2024-05-09', 'Employed', 70092, '1994-04-07', 4, 'sharon.wilson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0039', 'Female', 'sharon.wilson@personal.com'),
(40, 'Sarah', 'Rivera', 'Recruiter', 'People', 'London', 'Employed', '2024-10-13', 'Employed', 94398, '1995-06-24', 5, 'sarah.rivera@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0040', 'Female', 'sarah.rivera@personal.com'),
(41, 'Ronald', 'Nelson', 'Designer', 'Design', 'Paris', 'Employed', '2024-11-05', 'Employed', 91747, '1990-08-08', 2, 'ronald.nelson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0041', 'Male', 'ronald.nelson@personal.com'),
(42, 'Steven', 'Johnson', 'Account Executive', 'Sales', 'New York', 'Employed', '2022-10-28', 'Employed', 91808, '1999-06-27', 3, 'steven.johnson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0042', 'Male', 'steven.johnson@personal.com'),
(43, 'Amanda', 'Baker', 'Recruiter', 'People', 'New York', 'Employed', '2025-02-11', 'Employed', 74733, '1997-07-20', 5, 'amanda.baker@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0043', 'Female', 'amanda.baker@personal.com'),
(44, 'Ronald', 'Walker', 'Sales Representative', 'Sales', 'New York', 'Employed', '2023-12-23', 'Employed', 93674, '1994-04-22', 3, 'ronald.walker@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0044', 'Male', 'ronald.walker@personal.com'),
(45, 'Emily', 'Rodriguez', 'Senior Designer', 'Design', 'New York', 'Employed', '2021-02-10', 'Employed', 87319, '1998-03-29', 2, 'emily.rodriguez@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0045', 'Female', 'emily.rodriguez@personal.com'),
(46, 'Ashley', 'Brown', 'UI/UX Designer', 'Design', 'London', 'Employed', '2024-09-30', 'Employed', 65698, '1993-02-06', 2, 'ashley.brown@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0046', 'Female', 'ashley.brown@personal.com'),
(47, 'Elizabeth', 'Hill', 'Full Stack Developer', 'Tech', 'London', 'Employed', '2025-04-01', 'Employed', 92648, '1992-12-15', 4, 'elizabeth.hill@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0047', 'Female', 'elizabeth.hill@personal.com'),
(48, 'Amy', 'Johnson', 'Compliance Officer', 'Legal', 'Paris', 'Employed', '2022-08-19', 'Employed', 74253, '2000-03-06', 7, 'amy.johnson@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0048', 'Female', 'amy.johnson@personal.com'),
(49, 'Scott', 'Martinez', 'Accountant', 'Finance', 'London', 'Employed', '2024-04-12', 'Employed', 84547, '1990-02-04', 8, 'scott.martinez2@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0049', 'Male', 'scott.martinez@personal.com'),
(50, 'Jeffrey', 'Thomas', 'UI/UX Designer', 'Design', 'London', 'Employed', '2021-06-10', 'Employed', 52332, '1991-11-06', 2, 'jeffrey.thomas@workops.com', '$2a$12$11rGxnGTikibUQFJHfhxOuzq7oqe7CPg6Lyddlhz6bHxVJfMwvRBO', '+1 (555) 000-0050', 'Male', 'jeffrey.thomas@personal.com');

-- Insert default menu permissions
INSERT INTO menu_permissions (menu_item, permission_type, permission_value) VALUES
-- Everyone can see these
('dashboard', 'everyone', NULL),
('time off', 'everyone', NULL),
('thanks', 'everyone', NULL),
('hierarchy', 'everyone', NULL),

-- Associates - People department and CEO
('associates', 'department', 'People'),
('associates', 'title', 'CEO'),

-- Admin - CEO and People department
('admin', 'title', 'CEO'),
('admin', 'department', 'People'),

-- Task - CEO and People department
('tasks', 'title', 'CEO'),
('tasks', 'department', 'People'),

-- My Team - Managers and CEO
('my team', 'title', 'CEO'),
('my team', 'title', 'Manager'),
('my team', 'title', 'Head of Design'),
('my team', 'title', 'Head of Sales'),
('my team', 'title', 'Head of Tech'),
('my team', 'title', 'Head of People'),
('my team', 'title', 'Head of Product'),
('my team', 'title', 'Head of Legal'),
('my team', 'title', 'Head of Finance'),

-- Time Entry Approvals - Managers and CEO
('time entry approvals', 'title', 'CEO'),
('time entry approvals', 'title', 'Manager'),
('time entry approvals', 'title', 'Head of Design'),
('time entry approvals', 'title', 'Head of Sales'),
('time entry approvals', 'title', 'Head of Tech'),
('time entry approvals', 'title', 'Head of People'),
('time entry approvals', 'title', 'Head of Product'),
('time entry approvals', 'title', 'Head of Legal'),
('time entry approvals', 'title', 'Head of Finance'),

-- Time Off Approvals - Managers and CEO
('time off approvals', 'title', 'CEO'),
('time off approvals', 'title', 'Manager'),
('time off approvals', 'title', 'Head of Design'),
('time off approvals', 'title', 'Head of Sales'),
('time off approvals', 'title', 'Head of Tech'),
('time off approvals', 'title', 'Head of People'),
('time off approvals', 'title', 'Head of Product'),
('time off approvals', 'title', 'Head of Legal'),
('time off approvals', 'title', 'Head of Finance');
