INSERT INTO drivers (id, lat, lng)
SELECT uuid_generate_v4(), random()*180-90, random()*360-180
FROM generate_series(1, 10000);