CREATE TABLE IF NOT EXISTS scans (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    image VARCHAR(500) NOT NULL,
    container_id VARCHAR(12) NOT NULL,
    total_time_ms FLOAT NOT NULL,
    total_checks INTEGER NOT NULL,
    passed INTEGER NOT NULL,
    warnings INTEGER NOT NULL,
    failed INTEGER NOT NULL,
    critical INTEGER NOT NULL,
    high INTEGER NOT NULL,
    medium INTEGER NOT NULL,
    low INTEGER NOT NULL,
    has_failures BOOLEAN NOT NULL,
    has_warnings BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS check_results (
    id SERIAL PRIMARY KEY,
    scan_id INTEGER NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    severity VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    duration_ms FLOAT NOT NULL,
    output TEXT,
    output_preview TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS check_issues (
    id SERIAL PRIMARY KEY,
    check_result_id INTEGER NOT NULL REFERENCES check_results(id) ON DELETE CASCADE,
    issue TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_scans_timestamp ON scans(timestamp DESC);
CREATE INDEX idx_scans_image ON scans(image);
CREATE INDEX idx_check_results_scan_id ON check_results(scan_id);
CREATE INDEX idx_check_issues_check_result_id ON check_issues(check_result_id);