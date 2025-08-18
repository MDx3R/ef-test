CREATE INDEX idx_subscriptions_user_service_start_end 
ON subscriptions(user_id, service_name, start_date, end_date);

CREATE INDEX idx_subscriptions_user_service_start_active 
ON subscriptions(user_id, service_name, start_date)
WHERE end_date IS NULL;