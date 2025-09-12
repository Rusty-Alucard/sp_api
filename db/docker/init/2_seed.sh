psql -U sp -d sp -c "\copy events(name,name_jp) from '/docker-entrypoint-initdb.d/events.csv' with csv QUOTE '\"'"
psql -U sp -d sp -c "\copy teams(fifa_trigramma,name,name_jp,confederation) from '/docker-entrypoint-initdb.d/teams.csv' with csv QUOTE '\"'"
