INSERT INTO request_categories (name, description)
VALUES
    ('Electricity', 'Power, lighting and electrical faults'),
    ('Furniture', 'Damaged desks, chairs and fittings'),
    ('Plumbing', 'Leaking pipes, taps and drainage'),
    ('Internet', 'Wi-Fi and network connectivity'),
    ('Classroom Equipment', 'Projectors, boards and AV gear'),
    ('Hostel Maintenance', 'General hostel repairs')
ON CONFLICT (name) DO NOTHING;
