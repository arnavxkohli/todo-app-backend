CREATE TABLE users (
    UserID VARCHAR(36) PRIMARY KEY,
    Username VARCHAR(50) UNIQUE NOT NULL,
    Password VARCHAR(50) NOT NULL
);

CREATE TABLE todos (
    TodoID VARCHAR(36) PRIMARY KEY,
    CreatedDate DATE NOT NULL,
    DueDate DATE,
    Info VARCHAR(1000)
);

CREATE TABLE user_todos (
    TodoID VARCHAR(36) REFERENCES todos(TodoID),
    UserID VARCHAR(36) REFERENCES users(UserID)
);

