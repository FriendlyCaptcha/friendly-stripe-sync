-- name: SetSyncState :exec
INSERT INTO "stripe"."sync_state" (id, last_event) VALUES ('current_state', $1) ON CONFLICT (id) DO UPDATE SET last_event = EXCLUDED.last_event;

-- name: GetCurrentSyncState :one
SELECT * FROM "stripe"."sync_state" WHERE id = 'current_state';

-- name: DeleteCurrentSyncState :exec
DELETE FROM "stripe"."sync_state" WHERE id = 'current_state';