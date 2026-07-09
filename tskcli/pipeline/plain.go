package pipeline

// ASCIIArt is the --plain variant of the compact pipeline diagram.
const ASCIIArt = `       *
       |
       v
  +--------+
  | create |
  +---+----+
      | claim
      v
 +------------+
 | in_process |
 +-----+------+
       | research
       v
+---------------+
| clarification |<-- refine
+-------+-------+
        | confirmed
        v
+----------------+
| implementation |
+-------+--------+
        v
 +--------------+
 | verification |
 +------+-------+
        v
   +---------+
   | summary |
   +----+----+
        +-no followup>+------+
        |             | done |
        |questions    +--+---+
        v                |
+---------------+        |
| user_followup |        |
+-------+-------+        |
        +----->+ satisfied
        +-------+--------+
                v
                @`