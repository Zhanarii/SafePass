
CREATE TYPE user_role AS ENUM ('admin', 'supervisor', 'security', 'worker');
CREATE TYPE detection_status AS ENUM ('compliant', 'violation', 'warning');
CREATE TYPE ppe_item AS ENUM ('helmet', 'vest', 'boots', 'gloves', 'goggles', 'mask');
CREATE TYPE violation_severity AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE incident_status AS ENUM ('open', 'investigating', 'resolved', 'closed');
CREATE TYPE access_decision AS ENUM ('allowed', 'denied', 'manual_review');


CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    role user_role NOT NULL DEFAULT 'worker',
    department VARCHAR(100),
    badge_number VARCHAR(50) UNIQUE,
    photo_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    location GEOMETRY(Point, 4326), 
    type VARCHAR(50) NOT NULL, 
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE cameras (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    rtsp_url TEXT NOT NULL,
    position GEOMETRY(Point, 4326), 
    viewing_angle INTEGER, 
    fps INTEGER DEFAULT 10,
    resolution VARCHAR(20) DEFAULT '1920x1080',
    is_active BOOLEAN DEFAULT true,
    last_heartbeat TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE access_zones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    camera_id UUID REFERENCES cameras(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    required_ppe ppe_item[] NOT NULL, 
    danger_level violation_severity DEFAULT 'medium',
    access_rules JSONB, 
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE detections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    access_zone_id UUID REFERENCES access_zones(id) ON DELETE SET NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    frame_url TEXT, 
    detected_ppe ppe_item[], 
    missing_ppe ppe_item[], 
    confidence_scores JSONB, 
    bounding_boxes JSONB, 
    status detection_status NOT NULL,
    face_embedding BYTEA, 
    processing_time_ms INTEGER, 
    model_version VARCHAR(50)
);

CREATE INDEX idx_detections_camera_timestamp ON detections(camera_id, timestamp DESC);
CREATE INDEX idx_detections_status ON detections(status) WHERE status != 'compliant';
CREATE INDEX idx_detections_user ON detections(user_id) WHERE user_id IS NOT NULL;

CREATE TABLE violations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    detection_id UUID NOT NULL REFERENCES detections(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    access_zone_id UUID REFERENCES access_zones(id) ON DELETE SET NULL,
    violation_type VARCHAR(100) NOT NULL, 
    severity violation_severity NOT NULL,
    description TEXT,
    snapshot_url TEXT,
    video_url TEXT,
    access_decision access_decision DEFAULT 'denied',
    notified_at TIMESTAMP WITH TIME ZONE,
    acknowledged_by UUID REFERENCES users(id) ON DELETE SET NULL,
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_violations_user ON violations(user_id);
CREATE INDEX idx_violations_camera_date ON violations(camera_id, created_at DESC);
CREATE INDEX idx_violations_severity ON violations(severity);


CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    violation_id UUID NOT NULL REFERENCES violations(id) ON DELETE CASCADE,
    camunda_process_instance_id VARCHAR(100) UNIQUE,
    incident_number VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status incident_status NOT NULL DEFAULT 'open',
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    root_cause TEXT,
    corrective_actions TEXT,
    due_date TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    closed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_assigned ON incidents(assigned_to) WHERE assigned_to IS NOT NULL;


CREATE TABLE incident_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    comment TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE access_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    access_zone_id UUID REFERENCES access_zones(id) ON DELETE SET NULL,
    detection_id UUID REFERENCES detections(id) ON DELETE SET NULL,
    decision access_decision NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    badge_scanned BOOLEAN DEFAULT false,
    override_by UUID REFERENCES users(id) ON DELETE SET NULL, 
    override_reason TEXT
);

CREATE INDEX idx_access_logs_user ON access_logs(user_id, timestamp DESC);
CREATE INDEX idx_access_logs_camera ON access_logs(camera_id, timestamp DESC);


CREATE TABLE daily_statistics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_detections INTEGER DEFAULT 0,
    compliant_detections INTEGER DEFAULT 0,
    violations_count INTEGER DEFAULT 0,
    unique_workers INTEGER DEFAULT 0,
    avg_processing_time_ms FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(location_id, date)
);


CREATE TABLE notification_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    location_id UUID REFERENCES locations(id) ON DELETE CASCADE,
    rule_name VARCHAR(255) NOT NULL,
    trigger_condition JSONB NOT NULL, 
    notification_channels TEXT[] NOT NULL,
    recipients JSONB NOT NULL, 
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';


CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_locations_updated_at BEFORE UPDATE ON locations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cameras_updated_at BEFORE UPDATE ON cameras
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_access_zones_updated_at BEFORE UPDATE ON access_zones
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_incidents_updated_at BEFORE UPDATE ON incidents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE detections_2025_01 PARTITION OF detections
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE detections_2025_02 PARTITION OF detections
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');


COMMENT ON TABLE detections IS 'Каждая детекция от YOLO модели (high volume)';
COMMENT ON TABLE violations IS 'Агрегированные нарушения для анализа и отчётности';
COMMENT ON TABLE incidents IS 'Camunda процессы расследования нарушений';
COMMENT ON TABLE access_logs IS 'Журнал решений о доступе в зоны';
COMMENT ON TABLE daily_statistics IS 'Агрегированная статистика для дашбордов';