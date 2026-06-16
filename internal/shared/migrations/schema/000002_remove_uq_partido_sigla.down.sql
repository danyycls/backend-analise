-- Remove duplicatas de sigla que possam ter sido criadas antes de recriar a constraint
DELETE FROM partido p1 USING (
    SELECT sigla, MIN(id) AS min_id
    FROM partido
    GROUP BY sigla
    HAVING COUNT(*) > 1
) p2
WHERE p1.sigla = p2.sigla AND p1.id != p2.min_id;

ALTER TABLE partido ADD CONSTRAINT uq_partido_sigla UNIQUE (sigla);
