ALTER TABLE system_metrics ADD PRIMARY KEY (id, created_at);
ALTER TABLE system_metrics DROP CONSTRAINT system_metrics_pkey;
select create_hypertable('public.system_metrics', 'created_at', if_not_exists => TRUE);
CREATE MATERIALIZED view if not exists  system_metrics_hourly
WITH (timescaledb.continuous) AS
SELECT
    -- (1) Kelompokkan waktu ke dalam interval per 1 jam
    time_bucket('1 hour', created_at) AS hour,

    -- (2) Hitung statistik untuk CPU
    AVG(cpu_usage) AS avg_cpu_usage,
    MAX(cpu_usage) AS max_cpu_usage,
    MIN(cpu_usage) AS min_cpu_usage,

    -- (3) Hitung statistik untuk Memori
    AVG(mem_usage_byte) AS avg_mem_usage,
    MAX(mem_usage_byte) AS max_mem_usage,
    -- Ambil nilai total memori (biasanya tidak berubah dalam 1 jam)
    AVG(mem_total_byte) AS mem_total,

    -- (4) Hitung statistik untuk Disk
    AVG(disk_usage_byte) AS avg_disk_usage,
    MAX(disk_usage_byte) AS max_disk_usage,
    AVG(disk_total_byte) AS disk_total
FROM
    system_metrics
GROUP BY
    hour;

SELECT add_continuous_aggregate_policy('system_metrics_hourly',
    start_offset => INTERVAL '3 hours',
    end_offset   => INTERVAL '1 hour',
    schedule_interval => INTERVAL '30 minutes'
);

CREATE MATERIALIZED view if not exists system_metrics_daily
WITH (timescaledb.continuous) AS
SELECT
    -- (1) Kelompokkan waktu ke dalam interval per 1 jam
    time_bucket('1 day', created_at) AS days,

    -- (2) Hitung statistik untuk CPU
    AVG(cpu_usage) AS avg_cpu_usage,
    MAX(cpu_usage) AS max_cpu_usage,
    MIN(cpu_usage) AS min_cpu_usage,

    -- (3) Hitung statistik untuk Memori
    AVG(mem_usage_byte) AS avg_mem_usage,
    MAX(mem_usage_byte) AS max_mem_usage,
    -- Ambil nilai total memori (biasanya tidak berubah dalam 1 jam)
    AVG(mem_total_byte) AS mem_total,

    -- (4) Hitung statistik untuk Disk
    AVG(disk_usage_byte) AS avg_disk_usage,
    MAX(disk_usage_byte) AS max_disk_usage,
    AVG(disk_total_byte) AS disk_total
FROM
    system_metrics
GROUP BY
    days;

SELECT add_continuous_aggregate_policy('system_metrics_daily',
    start_offset => INTERVAL '3 days',
    end_offset   => INTERVAL '1 day',
    schedule_interval => INTERVAL '1 hour'
);