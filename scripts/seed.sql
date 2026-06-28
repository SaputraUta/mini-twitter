-- Realistic seed for the fan-out benchmark.
-- ~10k users (implicit ids 1..10000), celebrities = ids 1..20.
-- Power-law follows: everyone follows all 20 celebs + ~30 random users.
-- ~500k tweets, 30% authored by celebrities.
--
-- Run:  docker exec -i mt-postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < scripts/seed.sql
-- Then: cd services/timeline && go run ./cmd/backfill   (fills Redis timelines)

-- 1. clean + reset id sequences
TRUNCATE tweets, follows RESTART IDENTITY;

-- 2. everyone follows the 20 celebrities (ids 1..20)
INSERT INTO follows (follower_id, followee_id)
SELECT u, c
FROM generate_series(1, 10000) u
CROSS JOIN generate_series(1, 20) c
WHERE u <> c
ON CONFLICT DO NOTHING;

-- 3. each user also follows ~30 random users
INSERT INTO follows (follower_id, followee_id)
SELECT u, (random() * 9999)::int + 1
FROM generate_series(1, 10000) u
CROSS JOIN generate_series(1, 30)
ON CONFLICT DO NOTHING;

-- 4. 500k tweets, 30% from celebrities (ids 1..20), rest spread across all users
INSERT INTO tweets (user_id, text)
SELECT
  CASE WHEN random() < 0.3 THEN (random() * 19)::int + 1
       ELSE (random() * 9999)::int + 1 END,
  'tweet#' || g
FROM generate_series(1, 500000) g;
