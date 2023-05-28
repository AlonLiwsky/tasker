CREATE TABLE IF NOT EXISTS task (
                                    id INT PRIMARY KEY AUTO_INCREMENT,
                                    name VARCHAR(255) NOT NULL
    );

CREATE TABLE IF NOT EXISTS step (
                                    id INT PRIMARY KEY AUTO_INCREMENT,
                                    task_id INT NOT NULL,
                                    step_type VARCHAR(255) NOT NULL,
    params VARCHAR(255),
    failure_step INT,
    position INT NOT NULL,
    FOREIGN KEY (task_id) REFERENCES task(id),
    FOREIGN KEY (failure_step) REFERENCES step(id),
    INDEX idx_position (position)
    );

CREATE TABLE IF NOT EXISTS scheduled_task (
                                              id INT PRIMARY KEY AUTO_INCREMENT,
                                              name VARCHAR(255) NOT NULL,
    chron VARCHAR(255) NOT NULL,
    retry_policy VARCHAR(255),
    task_id INT NOT NULL,
    enabled BOOLEAN NOT NULL,
    last_run DATETIME,
    first_run DATETIME,
    FOREIGN KEY (task_id) REFERENCES task(id)
    );

CREATE TABLE IF NOT EXISTS execution (
                                         id INT PRIMARY KEY AUTO_INCREMENT,
                                         scheduled_task_id INT NOT NULL,
                                         try_number INT NOT NULL,
                                         status VARCHAR(255) NOT NULL,
    requested_time DATETIME,
    executed_time DATETIME,
    last_status_change_time DATETIME,
    FOREIGN KEY (scheduled_task_id) REFERENCES scheduled_task(id)
    );
