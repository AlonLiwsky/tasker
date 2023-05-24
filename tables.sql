CREATE TABLE IF NOT EXISTS task (
                                    ID INT PRIMARY KEY AUTO_INCREMENT,
                                    Name VARCHAR(255) NOT NULL
    );

CREATE TABLE IF NOT EXISTS step (
                                    ID INT PRIMARY KEY AUTO_INCREMENT,
                                    TaskID INT NOT NULL,
                                    SuccessNextStep INT,
                                    FailureNextStep INT,
                                    StÏepID INT NOT NULL,
                                    Params VARCHAR(255),
    FOREIGN KEY (TaskID) REFERENCES task(ID),
    FOREIGN KEY (SuccessNextStep) REFERENCES step(ID),
    FOREIGN KEY (FailureNextStep) REFERENCES step(ID)
    );

CREATE TABLE IF NOT EXISTS scheduled_task (
                                              ID INT PRIMARY KEY AUTO_INCREMENT,
                                              Name VARCHAR(255) NOT NULL,
    Chron VARCHAR(255) NOT NULL,
    RetryPolicy VARCHAR(255),
    TaskID INT NOT NULL,
    Enabled BOOLEAN NOT NULL,
    LastRun DATETIME,
    FirstRun DATETIME,
    FOREIGN KEY (TaskID) REFERENCES task(ID)
    );

CREATE TABLE IF NOT EXISTS execution (
                                         ID INT PRIMARY KEY AUTO_INCREMENT,
                                         ScheduledTaskID INT NOT NULL,
                                         TryNumber INT NOT NULL,
                                         Status VARCHAR(255) NOT NULL,
    RequestedTime DATETIME,
    ExecutedTime DATETIME,
    LastStatusChangeTime DATETIME,
    FOREIGN KEY (ScheduledTaskID) REFERENCES scheduled_task(ID)
    );
