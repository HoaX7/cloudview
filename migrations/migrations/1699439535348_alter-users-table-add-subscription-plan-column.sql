-- Up Migration
ALTER TABLE users ADD COLUMN subscription_plan_id int;
ALTER TABLE users ADD CONSTRAINT 
    fk_subscription_plan foreign key(subscription_plan_id) 
    references subscription_plans(id);
CREATE INDEX user_email_idx on users (email);

-- Down Migration
ALTER TABLE users DROP CONSTRAINT fk_subscription_plan;
ALTER TABLE users DROP COLUMN subscription_plan_id;
DROP INDEX user_email_idx;