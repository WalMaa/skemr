INSERT INTO projects (id, name, created_at, updated_at) VALUES
('aaaabbbb-aaaa-aaaa-aaaa-aaaabbbbcccc', 'Demo Project', NOW(), NOW());

INSERT INTO databases (id, display_name, db_name, username, password, host, port, database_type, project_id) VALUES
('11112222-3333-4444-5555-666677778888', 'Test Database', 'test_db', 'test_user', 'test_pass', 'localhost', 5433, 'postgres', 'aaaabbbb-aaaa-aaaa-aaaa-aaaabbbbcccc');