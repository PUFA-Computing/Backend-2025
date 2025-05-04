ALTER TABLE users
ADD COLUMN institution_name VARCHAR(255),
ADD COLUMN gender VARCHAR(6),

ALTER COLUMN student_id DROP NOT NULL,
ALTER COLUMN student_id SET DEFAULT NULL,
ALTER COLUMN major DROP NOT NULL,
ALTER COLUMN major SET DEFAULT NULL;




ALTER TABLE events
ADD COLUMN open_for_all BOOLEAN NOT NULL DEFAULT FALSE;

-- Function trigger insert or update users table to update profile_picture column if gender male then https://sg.pufacomputing.live/Assets/male.jpeg else https://sg.pufacomputing.live/Assets/female.jpeg
CREATE OR REPLACE FUNCTION update_profile_picture_based_on_gender()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.gender = 'male' THEN
        NEW.profile_picture := 'https://sg.pufacomputing.live/Assets/male.jpeg';
    ELSIF NEW.gender = 'female' THEN
        NEW.profile_picture := 'https://sg.pufacomputing.live/Assets/female.jpeg';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_profile_picture_trigger
    BEFORE INSERT ON users
    FOR EACH ROW
EXECUTE FUNCTION update_profile_picture_based_on_gender();