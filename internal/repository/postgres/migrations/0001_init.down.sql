DROP INDEX IF EXISTS  room_restrictions_reservation_id_idx, room_restrictions_room_id_idx,
                      room_restrictions_start_date_end_date_idx, room_restrictions_last_name_idx,
                      room_restrictions_email_idx CASCADE;
DROP TABLE IF EXISTS users, rooms, reservations, room_restrictions, restrictions CASCADE;
