
-- Initialize database

DROP DATABASE go_geo WITH (FORCE); 
CREATE DATABASE go_geo;

\c go_geo

CREATE EXTENSION postgis;

CREATE SEQUENCE farms_id_seq;

CREATE TABLE farms (
    id integer DEFAULT nextval('farms_id_seq'::regclass) PRIMARY KEY,
    name text
);

CREATE SEQUENCE fields_id_seq;

CREATE TABLE fields (
    id integer DEFAULT nextval('fields_id_seq'::regclass) PRIMARY KEY,
	farm_id integer,
    name text
);

CREATE SEQUENCE crop_id_seq;

CREATE TABLE crops (
    id integer DEFAULT nextval('fields_id_seq'::regclass) PRIMARY KEY,
    name text
);

INSERT INTO crops (name) VALUES 
	('Corn'),
	('Soybeans');

CREATE SEQUENCE planting_pts_id_seq;

CREATE TABLE planting_pts (
    id integer DEFAULT nextval('planting_pts_id_seq'::regclass) PRIMARY KEY,
	field_id integer,
	crop_id integer,
	variety_id integer,
	time timestamp,
	section integer,
	swath_width_ft float,
	distance_ft float,
	heading_deg float,
	elevation_ft float,
	target_rate float,
	applied_rate float,
    geog geography(Point, 4326)
);

CREATE INDEX idx_planting_pts_geog
    ON planting_pts USING gist (geog);

CREATE SEQUENCE harvest_pts_id_seq;

CREATE TABLE harvest_pts (
    id integer DEFAULT nextval('harvest_pts_id_seq'::regclass) PRIMARY KEY,
    planting_pt_id integer,
	field_id integer,
	crop_id integer,
	time timestamp,
	section integer,
	swath_width_ft float,
	distance_ft float,
	heading_deg float,
	elevation_ft float,
	yield_bu_ac float,
	moisture_per float,
    geog geography(Point, 4326)
);

CREATE INDEX idx_harvest_pts_geog
    ON harvest_pts USING gist (geog);
