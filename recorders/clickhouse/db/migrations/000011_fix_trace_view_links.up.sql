-- trace_view
CREATE OR REPLACE VIEW trace_view AS
SELECT
    TraceId AS traceID,
    SpanId AS spanID,
    SpanName AS operationName,
    ParentSpanId AS parentSpanID,
    ServiceName AS serviceName,
    Duration / 1000000 AS duration,
    Timestamp AS startTime,
    arrayMap(key -> map('key', key, 'value', SpanAttributes[key]), mapKeys(SpanAttributes)) AS tags,
    arrayMap(key -> map('key', key, 'value', ResourceAttributes[key]), mapKeys(ResourceAttributes)) AS serviceTags,
    arrayMap((traceId, spanId) -> map('traceID', traceId, 'spanID', spanId), Links.TraceId, Links.SpanId) as references
FROM traces
WHERE TraceId = {trace_id:String}
;
