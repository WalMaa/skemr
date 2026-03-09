INSERT INTO projects (id, name, created_at, updated_at) VALUES
('aaaabbbb-aaaa-aaaa-aaaa-aaaabbbbcccc', 'Demo Project', NOW(), NOW());

INSERT INTO databases (id, display_name, db_name, username, password, host, port, database_type, project_id) VALUES
('11112222-3333-4444-5555-666677778888', 'Test Database', 'test_db', 'test_user', 'test_pass', 'localhost', 5433, 'postgres', 'aaaabbbb-aaaa-aaaa-aaaa-aaaabbbbcccc');


--- The token is "wg8h807NHOg4.6jfEHDrS_KrQDS65uB_L9VbRErPAKU0PHz48LVIEUJM"
INSERT INTO project_access_tokens (id, project_id, name, prefix, hash, last_used, expires_at, created_at, updated_at) VALUES
('22223333-4444-5555-6666-777788889999', 'aaaabbbb-aaaa-aaaa-aaaa-aaaabbbbcccc', 'Demo Token', 'wg8h807NHOg4', '$argon2id$v=19$m=65536,t=2,p=1$CF0S475vCJyRHflYiGzmTg$XQAdy7e3pU5Re1UqPlf9zexACIfcIGQEEqSDXGVRYUk', NULL, NULL, NOW(), NOW());