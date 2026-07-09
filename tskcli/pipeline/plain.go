package pipeline

// ASCIIArt is the --plain variant of the compact pipeline diagram.
//
// user_followup has two exits:
//   - refine → loops back to clarification (right rail)
//   - satisfied → joins no-followup into done → terminal
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
| clarification |<----------+
+-------+-------+           |
        | confirmed         |
        v                   |
+----------------+          |
| implementation |          |
+-------+--------+          |
        v                   |
 +--------------+           |
 | verification |           |
 +------+-------+           |
        v                   |
   +---------+              |
   | summary |              |
   +----+----+              |
        +-no followup---+   |
        |               |   |
        |questions      |   |
        v               |   |
+---------------+       |   |
| user_followup |       |   |
+---+-------+---+       |   |
    |       |           |   |
    |       +-satisfied>+   |
    |                   |   |
    |              +----+-+ |
    |              | done | |
    |              +--+---+ |
    |                 v     |
    |                 @     |
    +------refine-----------+`
